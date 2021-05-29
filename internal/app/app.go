package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/handler/grpc"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	PluginManager service.PluginManager
	GRPC          *grpc.Server
}

type App struct {
	pluginManager service.PluginManager
	grpc          *grpc.Server
}

func New(params Params) (*App, error) {
	return &App{
		pluginManager: params.PluginManager,
		grpc:          params.GRPC,
	}, nil
}

func (a *App) Run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGHUP)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer func() {
			log.L().Info("gRPC server stopped")
			wg.Done()
		}()
		if err := a.grpc.Start(); err != nil {
			log.L().Fatal("gRPC error: %s", err.Error())
		}
	}()

	go func() {
		select {
		case <-sig:
			a.grpc.Shutdown()

			err := a.pluginManager.StopAll(context.Background())
			if err != nil {
				log.L().Error(err.Error())
			}
			log.L().Info("All plugins stopped")
			wg.Done()
		}
	}()

	wg.Wait()
}
