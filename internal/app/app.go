package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/domain/service"
	nori_http "github.com/nori-io/nori/internal/handler/http"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	PluginManager service.PluginManager
	Http          *nori_http.Handler
}

type App struct {
	pluginManager service.PluginManager
	http          *nori_http.Handler
}

func New(params Params) (*App, error) {
	return &App{
		pluginManager: params.PluginManager,
		http:          params.Http,
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
