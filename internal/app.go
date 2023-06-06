package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/goverland-labs/core-storage/internal/communicate"
	"github.com/goverland-labs/core-storage/internal/config"
	"github.com/goverland-labs/core-storage/internal/dao"
	"github.com/goverland-labs/core-storage/internal/events"
	"github.com/goverland-labs/core-storage/internal/proposal"
	"github.com/goverland-labs/core-storage/internal/vote"
	"github.com/goverland-labs/core-storage/pkg/grpcsrv"
	"github.com/goverland-labs/core-storage/pkg/health"
	"github.com/goverland-labs/core-storage/pkg/prometheus"
	"github.com/goverland-labs/core-storage/protobuf/internalapi"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
	db      *gorm.DB

	daoService      *dao.Service
	proposalService *proposal.Service
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
	db, err := gorm.Open(postgres.Open(a.cfg.DB.DSN), &gorm.Config{})
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

	err = a.initDao(nc, pb)
	if err != nil {
		return fmt.Errorf("init dao: %w", err)
	}

	err = a.initProposal(nc, pb)
	if err != nil {
		return fmt.Errorf("init proposal: %w", err)
	}

	err = a.initVote(nc)
	if err != nil {
		return fmt.Errorf("init vote: %w", err)
	}

	err = a.initAPI()
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	return nil
}

func (a *Application) initDao(nc *nats.Conn, pb *communicate.Publisher) error {
	repo := dao.NewRepo(a.db)
	service, err := dao.NewService(repo, pb)
	if err != nil {
		return fmt.Errorf("dao service: %w", err)
	}
	a.daoService = service

	cs, err := dao.NewConsumer(nc, service)
	if err != nil {
		return fmt.Errorf("dao consumer: %w", err)
	}

	a.manager.AddWorker(process.NewCallbackWorker("dao-consumer", cs.Start))

	return nil
}

func (a *Application) initProposal(nc *nats.Conn, pb *communicate.Publisher) error {
	erRepo := events.NewRepo(a.db)
	erService, err := events.NewService(erRepo)
	if err != nil {
		return fmt.Errorf("new events service: %w", err)
	}

	repo := proposal.NewRepo(a.db)
	service, err := proposal.NewService(repo, pb, erService)
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

	return nil
}

func (a *Application) initVote(nc *nats.Conn) error {
	repo := vote.NewRepo(a.db)
	service, err := vote.NewService(repo)
	if err != nil {
		return fmt.Errorf("vote service: %w", err)
	}

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
