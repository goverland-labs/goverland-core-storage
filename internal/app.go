package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/votingpb"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"
	"github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"

	"github.com/goverland-labs/goverland-core-storage/internal/config"
	"github.com/goverland-labs/goverland-core-storage/internal/dao"
	"github.com/goverland-labs/goverland-core-storage/internal/delegate"
	"github.com/goverland-labs/goverland-core-storage/internal/discord"
	"github.com/goverland-labs/goverland-core-storage/internal/ensresolver"
	"github.com/goverland-labs/goverland-core-storage/internal/events"
	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
	"github.com/goverland-labs/goverland-core-storage/internal/pubsub"
	"github.com/goverland-labs/goverland-core-storage/internal/stats"
	"github.com/goverland-labs/goverland-core-storage/internal/vote"
	"github.com/goverland-labs/goverland-core-storage/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-core-storage/pkg/health"
	"github.com/goverland-labs/goverland-core-storage/pkg/prometheus"
	zerionsdk "github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"
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

	delegateRepo    *delegate.Repo
	delegateService *delegate.Service

	ensRepo    *ensresolver.Repo
	ensService *ensresolver.Service

	eventsRepo    *events.Repo
	eventsService *events.Service

	statsService *stats.Service

	zerionClient  *zerionsdk.Client
	discordSender *discord.Sender
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
	a.delegateRepo = delegate.NewRepo(a.db)

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

	pb, err := natsclient.NewPublisher(nc)
	if err != nil {
		return err
	}

	a.initZerionAPI()
	a.initDiscordSender()

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

	err = a.initDelegates(nc, pb)
	if err != nil {
		return fmt.Errorf("init delegates: %w", err)
	}

	a.initStats()

	err = a.initAPI()
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	return nil
}

func (a *Application) initEnsResolver(pb *natsclient.Publisher) error {
	conn, err := grpc.NewClient(a.cfg.InternalAPI.EnsResolverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with ens resolver: %v", err)
	}

	a.ensService, err = ensresolver.NewService(a.ensRepo, enspb.NewEnsClient(conn), pb)
	if err != nil {
		return fmt.Errorf("ensresolver.NewService: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("ens-resolver", a.ensService.Start))

	return nil
}

func (a *Application) initDao(nc *nats.Conn, pb *natsclient.Publisher) error {
	a.daoIDService = dao.NewDaoIDService(a.daoIDRepo)

	topDAOCache := dao.NewTopDAOCache(a.daoRepo)
	fungibleChainRepo := dao.NewFungibleChainRepo(a.db)

	service, err := dao.NewService(a.daoRepo, a.daoUniqueRepo, a.daoIDService, pb, a.proposalRepo, topDAOCache, fungibleChainRepo, a.zerionClient, a.discordSender)
	if err != nil {
		return fmt.Errorf("dao service: %w", err)
	}
	a.daoService = service
	if err = service.PrefillDaoIDs(); err != nil {
		return fmt.Errorf("PrefillDaoIDs: %w", err)
	}

	fungibleChainWorker := dao.NewFungibleChainWorker(a.zerionClient, service, fungibleChainRepo)

	cs, err := dao.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("dao consumer: %w", err)
	}
	a.manager.AddWorker(process.NewCallbackWorker("dao-consumer", cs.Start))

	cw := dao.NewNewCategoryWorker(service)
	mc := dao.NewVotersCountWorker(service)
	pcw := dao.NewPopularCategoryWorker(service)
	avw := dao.NewCntCalculationWorker(service)
	rw := dao.NewRecommendationWorker(service)
	tpw := dao.NewTokenPriceWorker(service, a.zerionClient)
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-category-process-worker", cw.ProcessNew))
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-category-outdated-worker", cw.RemoveOutdated))
	a.manager.AddWorker(process.NewCallbackWorker("dao-new-voters-worker", mc.ProcessNew))
	a.manager.AddWorker(process.NewCallbackWorker("dao-popular-category-process-worker", pcw.Process))
	a.manager.AddWorker(process.NewCallbackWorker("dao-active-votes-worker", avw.ProcessVotes))
	a.manager.AddWorker(process.NewCallbackWorker("dao-proposals-counter-worker", avw.ProcessProposalCounters))
	a.manager.AddWorker(process.NewCallbackWorker("top-dao-cache-worker", topDAOCache.Start))
	a.manager.AddWorker(process.NewCallbackWorker("dao-recommendations", rw.Process))
	a.manager.AddWorker(process.NewCallbackWorker("token-price", tpw.Process))
	a.manager.AddWorker(process.NewCallbackWorker("fungible-chain-worker", fungibleChainWorker.Start))

	return nil
}

func (a *Application) initProposal(nc *nats.Conn, pb *natsclient.Publisher) error {
	erService, err := events.NewService(a.eventsRepo)
	if err != nil {
		return fmt.Errorf("new events service: %w", err)
	}

	a.eventsService = erService

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

func (a *Application) initDelegates(nc *nats.Conn, pb *natsclient.Publisher) error {
	dsConn, err := grpc.NewClient(
		a.cfg.InternalAPI.DatasourceSnapshotAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create connection with datasource snapshot server: %v", err)
	}

	delegateClient := delegatepb.NewDelegateClient(dsConn)
	service := delegate.NewService(a.delegateRepo, delegateClient, a.daoService, a.proposalService, a.ensService, pb, a.eventsService)
	a.delegateService = service

	cs, err := delegate.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("delegates consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("delegates-allowed-daos", service.UpdateAllowedDaos))
	a.manager.AddWorker(process.NewCallbackWorker("delegates-consumer", cs.Start))

	ltw := delegate.NewLifeTimeWorker(service)
	a.manager.AddWorker(process.NewCallbackWorker("delegates-life-time-worker", ltw.Start))

	return nil
}

func (a *Application) initVote(nc *nats.Conn, pb *natsclient.Publisher) error {
	dsConn, err := grpc.NewClient(
		a.cfg.InternalAPI.DatasourceSnapshotAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create connection with datasource snapshot server: %v", err)
	}

	dsClient := votingpb.NewVotingClient(dsConn)
	votesNotifier := pubsub.NewPubSub[string](1000) // TODO: const

	service, err := vote.NewService(votesNotifier, a.voteRepo, a.daoService, pb, a.ensService, dsClient)
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

func (a *Application) initStats() {
	a.statsService = stats.NewService(a.daoRepo, a.proposalRepo)
	//cw := stats.NewCalcTotalsWorker(a.statsService)

	//a.manager.AddWorker(process.NewCallbackWorker("calc-totals", cw.Start))
}

func (a *Application) initAPI() error {
	authInterceptor := grpcsrv.NewAuthInterceptor()
	srv := grpcsrv.NewGrpcServer(
		[]string{
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		},
		authInterceptor.AuthAndIdentifyTickerFunc,
	)

	storagepb.RegisterDaoServer(srv, dao.NewServer(a.daoService))
	storagepb.RegisterProposalServer(srv, proposal.NewServer(a.proposalService))
	storagepb.RegisterVoteServer(srv, vote.NewServer(a.voteService))
	storagepb.RegisterEnsServer(srv, ensresolver.NewServer(a.ensService))
	storagepb.RegisterStatsServer(srv, stats.NewServer(a.statsService))
	storagepb.RegisterDelegateServer(srv, delegate.NewServer(a.delegateService, a.daoService))

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

func (a *Application) initZerionAPI() {
	zc := zerionsdk.NewClient(a.cfg.Zerion.BaseURL, a.cfg.Zerion.Key, http.DefaultClient)
	a.zerionClient = zc
}

func (a *Application) initDiscordSender() {
	ds := discord.NewSender(a.cfg.Discord.NewDaosURL)
	a.discordSender = ds
}
