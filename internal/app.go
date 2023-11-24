package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goverland-labs/datasource-snapshot/proto/votingpb"
	"github.com/goverland-labs/helpers-ens-resolver/proto"
	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/goverland-labs/core-api/protobuf/internalapi"

	"github.com/goverland-labs/core-storage/internal/communicate"
	"github.com/goverland-labs/core-storage/internal/config"
	"github.com/goverland-labs/core-storage/internal/dao"
	"github.com/goverland-labs/core-storage/internal/ensresolver"
	"github.com/goverland-labs/core-storage/internal/events"
	"github.com/goverland-labs/core-storage/internal/proposal"
	"github.com/goverland-labs/core-storage/internal/vote"
	"github.com/goverland-labs/core-storage/pkg/grpcsrv"
	"github.com/goverland-labs/core-storage/pkg/health"
	"github.com/goverland-labs/core-storage/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
	db      *gorm.DB

	proposalRepo    *proposal.Repo
	proposalService *proposal.Service

	daoIDRepo    *dao.DaoIDRepo
	daoIDService *dao.DaoIDService

	daoRepo       *dao.Repo
	daoUniqueRepo *dao.UniqueVoterRepo
	daoService    *dao.Service

	voteRepo    *vote.Repo
	voteService *vote.Service

	ensRepo    *ensresolver.Repo
	ensService *ensresolver.Service

	eventsRepo *events.Repo
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		a.initDB,

		// Init Dependencies
		a.initServices,

		// Init Workers: System
		a.initPrometheusWorker,
		a.initHealthWorker,
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initDB() error {
	db, err := gorm.Open(postgres.Open(a.cfg.DB.DSN), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		}),
	})
	if err != nil {
		return err
	}

	ps, err := db.DB()
	if err != nil {
		return err
	}
	ps.SetMaxOpenConns(a.cfg.DB.MaxOpenConnections)

	a.db = db
	if a.cfg.DB.Debug {
		a.db = db.Debug()
	}

	a.daoRepo = dao.NewRepo(a.db)
	a.daoUniqueRepo = dao.NewUniqueVoterRepo(a.db)
	a.daoIDRepo = dao.NewDaoIDRepo(a.db)
	a.proposalRepo = proposal.NewRepo(a.db)
	a.voteRepo = vote.NewRepo(a.db)
	a.eventsRepo = events.NewRepo(a.db)
	a.ensRepo = ensresolver.NewRepo(a.db)

	return err
}

func (a *Application) initServices() error {
	nc, err := nats.Connect(
		a.cfg.Nats.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(a.cfg.Nats.MaxReconnects),
		nats.ReconnectWait(a.cfg.Nats.ReconnectTimeout),
	)
	if err != nil {
		return err
	}

	pb, err := communicate.NewPublisher(nc)
	if err != nil {
		return err
	}

	err = a.initEnsResolver(pb)
	if err != nil {
		return fmt.Errorf("init dao: %w", err)
	}

	err = a.initDao(nc, pb)
	if err != nil {
		return fmt.Errorf("init dao: %w", err)
	}

	err = a.initProposal(nc, pb)
	if err != nil {
		return fmt.Errorf("init proposal: %w", err)
	}

	err = a.initVote(nc, pb)
	if err != nil {
		return fmt.Errorf("init vote: %w", err)
	}

	err = a.initAPI()
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	return nil
}

func (a *Application) initEnsResolver(pb *communicate.Publisher) error {
	conn, err := grpc.Dial(a.cfg.InternalAPI.EnsResolverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with ens resolver: %v", err)
	}

	a.ensService, err = ensresolver.NewService(a.ensRepo, proto.NewEnsClient(conn), pb)
	if err != nil {
		return fmt.Errorf("ensresolver.NewService: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("ens-resolver", a.ensService.Start))

	return nil
}

func (a *Application) initDao(nc *nats.Conn, pb *communicate.Publisher) error {
	a.daoIDService = dao.NewDaoIDService(a.daoIDRepo)

	service, err := dao.NewService(a.daoRepo, a.daoUniqueRepo, a.daoIDService, pb, a.proposalRepo)
	if err != nil {
		return fmt.Errorf("dao service: %w", err)
	}
	a.daoService = service
	if err = service.PrefillDaoIDs(); err != nil {
		return fmt.Errorf("PrefillDaoIDs: %w", err)
	}

	cs, err := dao.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("dao consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("dao-consumer", cs.Start))

	cw := dao.NewNewCategoryWorker(service)
	mc := dao.NewVotersCountWorker(service)
	pcw := dao.NewPopularCategoryWorker(service)
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-category-process-worker", cw.ProcessNew))
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-category-outdated-worker", cw.RemoveOutdated))
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-voters-worker", mc.ProcessNew))
	a.manager.AddWorker(process.NewCallbackWorker("dao-popular-category-process-worker", pcw.Process))

	return nil
}

func (a *Application) initProposal(nc *nats.Conn, pb *communicate.Publisher) error {
	erService, err := events.NewService(a.eventsRepo)
	if err != nil {
		return fmt.Errorf("new events service: %w", err)
	}

	service, err := proposal.NewService(a.proposalRepo, pb, erService, a.daoService, a.ensService)
	if err != nil {
		return fmt.Errorf("proposal service: %w", err)
	}
	a.proposalService = service

	cs, err := proposal.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("proposal consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("proposal-consumer", cs.Start))

	vw := proposal.NewVotingWorker(service)
	a.manager.AddWorker(process.NewCallbackWorker("voting-worker", vw.Start))

	tw := proposal.NewTopWorker(service)
	a.manager.AddWorker(process.NewCallbackWorker("proposal-top-worker", tw.Start))

	return nil
}

func (a *Application) initVote(nc *nats.Conn, pb *communicate.Publisher) error {
	dsConn, err := grpc.Dial(
		a.cfg.InternalAPI.DatasourceSnapshotAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create connection with datasource snapshot server: %v", err)
	}

	dsClient := votingpb.NewVotingClient(dsConn)

	service, err := vote.NewService(a.voteRepo, a.daoService, pb, dsClient)
	if err != nil {
		return fmt.Errorf("vote service: %w", err)
	}
	a.voteService = service

	cs, err := vote.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("vote consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("vote-consumer", cs.Start))

	return nil
}

func (a *Application) initAPI() error {
	authInterceptor := grpcsrv.NewAuthInterceptor()
	srv := grpcsrv.NewGrpcServer(
		[]string{
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		},
		authInterceptor.AuthAndIdentifyTickerFunc,
	)

	internalapi.RegisterDaoServer(srv, dao.NewServer(a.daoService))
	internalapi.RegisterProposalServer(srv, proposal.NewServer(a.proposalService))
	internalapi.RegisterVoteServer(srv, vote.NewServer(a.voteService))

	a.manager.AddWorker(grpcsrv.NewGrpcServerWorker("API", srv, a.cfg.InternalAPI.Bind))

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
