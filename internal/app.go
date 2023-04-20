package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/core-storage/internal/communicate"
	"github.com/goverland-labs/core-storage/internal/config"
	"github.com/goverland-labs/core-storage/internal/dao"
	"github.com/goverland-labs/core-storage/pkg/health"
	"github.com/goverland-labs/core-storage/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
	db      *gorm.DB
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

		// Init Workers: Application
		// TODO

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
	db, err := gorm.Open("postgres", a.cfg.DB.PostgresAddr)
	if err != nil {
		return err
	}

	a.db = db
	dao.AutoMigrate(db)

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

	return nil
}

func (a *Application) initDao(nc *nats.Conn, pb *communicate.Publisher) error {
	repo := dao.NewRepo(a.db)
	service, err := dao.NewService(repo, pb)
	if err != nil {
		return fmt.Errorf("dao service: %w", err)
	}

	_ = service

	cs, err := dao.NewConsumer(context.Background(), nc, service)
	if err != nil {
		return fmt.Errorf("dao consumer: %w", err)
	}

	a.manager.AddWorker(cs)

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

	err := a.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("close db connection")
	}

	i := 0
	for {
		if i > 10 {
			break
		}

		<-time.After(time.Millisecond * 200)
		i++
	}
}
