/*
Copyright 2018-2020 The App Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/nori-io/logger"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/app/container/env"
	"github.com/nori-io/nori/internal/app/container/handler"
	"github.com/nori-io/nori/internal/app/container/manager"
	ctnrNori "github.com/nori-io/nori/internal/app/container/nori"
	"github.com/nori-io/nori/internal/app/container/repository"
	"github.com/nori-io/nori/internal/app/container/service"
	"github.com/nori-io/nori/internal/nori"
	"go.uber.org/fx"
)

type App struct {
	fxOptions fx.Option
	//log       logger.Logger
	//grpc      *grpc.Server
	//http      *echo.Echo
	nori *nori.Nori
}

type Params struct {
	ConfigFile string
	Log        logger.Logger
}

func New(p Params) (*App, error) {
	var app = new(App)

	app.FxProvides(
		env.New(p.ConfigFile, p.Log),
		handler.New,
		manager.New,
		service.New,
		repository.New,
		ctnrNori.New,
	)

	return app, nil
}

func (app *App) FxProvides(ff ...func() fx.Option) {
	options := make([]fx.Option, len(ff))
	for i, f := range ff {
		options[i] = f()
	}
	app.fxOptions = fx.Options(options...)
}

func (app *App) Init() error {
	app.fxOptions = fx.Options(
		app.fxOptions,
		fx.NopLogger,

		fx.Invoke(
			func(nori *nori.Nori) (*App, error) {
				//app.http = http
				//app.grpc = grpc
				app.nori = nori
				return app, nil
			},
		),
	)

	err := fx.New(app.fxOptions).Start(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(context.Background())

	// todo: start grpc
	// todo: start http

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := app.nori.Run(ctx)
		if err != nil {
			// todo: log error
			println(err.Error())
			close(sig)
			return
		}
	}()

	<-sig

	cancel()

	wg.Wait()

	if err := app.nori.Stop(); err != nil {
		log.L().Error(err.Error())
	}

	return nil
}
