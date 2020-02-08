/*
Copyright 2019-2020 The Nori Authors.
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

package nori

import (
	"context"
	"os"
	"sync"

	"github.com/nori-io/nori-common/v2/logger"

	go_config "github.com/cheebo/go-config"
	"github.com/nori-io/nori/internal/plugins"
)

type Nori interface {
	Exec() error
}

type nori struct {
	config        Config
	log           logger.Logger
	pluginManager plugins.Manager
	sig           chan os.Signal
}

func NewNori(cfg go_config.Config, log logger.Logger, sig chan os.Signal) Nori {
	c := Config{
		Plugins: struct{ Dir []interface{} }{Dir: nil},
	}
	// plugins
	err := cfg.Unmarshal(&c.Plugins, "plugins")
	if err != nil {
		log.Fatal("%v", err)
	}
	// nori
	err = cfg.Unmarshal(&c.Nori, "nori")
	if err != nil {
		log.Fatal("%v", err)
	}

	return &nori{
		config:        c,
		log:           log,
		pluginManager: plugins.NewManager(cfg, log),
		sig:           sig,
	}
}

func (n *nori) Exec() error {
	ctx := context.Background()
	// todo: load config
	// todo: log config

	// hooks
	// storage

	// load plugin files
	m, err := n.pluginManager.AddDir(pluginDir(n.config.Plugins.Dir))
	if err != nil {
		n.log.Fatal("Cannot load plugins: %s", err.Error())
	}
	for _, d := range m {
		n.log.Info("Found plugin [%s] interface [%s]", d.Id(), d.GetInterface())
	}

	//filepath.Match()

	ctx, cancelFunc := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	// three services to shutdown: PluginManager, gRPC, HTTP
	wg.Add(3)

	// start plugins
	err = n.pluginManager.StartAll(ctx)
	if err != nil {
		n.log.Error("PluginManager cannot start all plugins: [%s]", err.Error())
	}

	go n.pm(ctx, wg)
	go n.rest(ctx, wg)
	go n.gRPC(ctx, wg)

	<-n.sig
	cancelFunc()
	wg.Wait()
	return nil
}

func (n *nori) rest(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// todo start REST server
	select {
	case <-ctx.Done():
		// todo: shutdown
		n.log.Info("Nori REST Server went down")
	}
}

func (n *nori) gRPC(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// todo start gRPC server
	select {
	case <-ctx.Done():
		// todo: shutdown
		n.log.Info("Nori gRPC Server went down")
	}
}

func (n *nori) pm(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// todo
	select {
	case <-ctx.Done():
		if err := n.pluginManager.StopAll(ctx); err != nil {
			n.log.Error("Plugin Manager stopped all with error [%s]", err.Error())
		} else {
			n.log.Info("Plugin Manager stopped all")
		}
	}
}

func pluginDir(list []interface{}) []string {
	dirs := make([]string, len(list))
	for i := 0; i < len(list); i++ {
		dirs[i] = list[i].(string)
	}
	return dirs
}
