package nori

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/nori-io/nori-common/logger"

	go_config "github.com/cheebo/go-config"
	"github.com/nori-io/nori/internal/plugins"
)

type Nori interface {
	Exec() error
}

type nori struct {
	config        Config
	logger        logger.Logger
	pluginManager plugins.Manager
	sig           chan os.Signal
}

func NewNori(cfg go_config.Config, log logger.Logger, sig chan os.Signal) Nori {
	c := Config{
		Plugins: struct{ Dir []interface{} }{Dir: nil},
	}
	err := cfg.Unmarshal(&c.Plugins, "plugins")
	if err != nil {
		log.Fatal("%v", err)
	}
	return &nori{
		config:        c,
		logger:        log,
		pluginManager: plugins.NewManager(cfg, log),
		sig:           sig,
	}
}

func (n *nori) Exec() error {
	// todo: load config
	// todo: logger config

	// todo: load files

	dirs := make([]string, len(n.config.Plugins.Dir))
	for i := 0; i < len(n.config.Plugins.Dir); i++ {
		dirs[i] = n.config.Plugins.Dir[i].(string)
	}

	m, err := n.pluginManager.AddDir(dirs)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range m {
		fmt.Println(d.Id())
	}

	// todo: load plugins
	// todo: start plugins

	// todo: start REST API server
	// todo: start gRPC server

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-n.sig

		if err := n.pluginManager.StopAll(context.Background()); err != nil {
			n.logger.Error("%v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
