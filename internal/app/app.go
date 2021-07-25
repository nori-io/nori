package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nori-io/common/v5/pkg/domain/storage"
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	nori_http "github.com/nori-io/nori/internal/handler/http"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	PluginManager service.PluginManager
	Http          *nori_http.Handler
	Env *env.Env
}

type App struct {
	pluginManager service.PluginManager
	http          *nori_http.Handler
	storage storage.Storage
}

func New(params Params) (*App, error) {
	return &App{
		pluginManager: params.PluginManager,
		http:          params.Http,
		storage: params.Env.Storage,
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
			wg.Done()
		}()
		if err := a.http.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.L().Info("http server stopped")
			} else {
				log.L().Fatal("http error: %s", err.Error())
			}
		}
	}()

	go func() {
		select {
		case <-sig:
			a.http.Stop(context.Background())
			if err := a.storage.Close(); err != nil {
				log.L().Error(err.Error())
			}

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
