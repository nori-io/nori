// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cheebo/go-config"
	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori/internal/nori"
	"github.com/spf13/cobra"
)

type channels struct {
	shutdown  chan struct{}
	interrupt chan os.Signal
}

// serverCmd represents the server command
func serverCmd(cfg go_config.Config, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Run: func(cmd *cobra.Command, args []string) {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGTERM)
			signal.Notify(sig, syscall.SIGINT)
			signal.Notify(sig, syscall.SIGHUP)

			app := nori.NewNori(cfg, log, sig)
			if err := app.Exec(); err != nil {
				log.Error(err.Error())
			}

			//// shutdown and interrupt
			//ch := channels{
			//	shutdown:  make(chan struct{}, 1),
			//	interrupt: make(chan os.Signal, 1),
			//}
			//signal.Notify(ch.interrupt, syscall.SIGTERM)
			//signal.Notify(ch.interrupt, syscall.SIGINT)
			//
			//noriVersion := version.NoriVersion(log)
			//
			//log.Infof("Nori Engine [version %s]", noriVersion.Version().String())
			//
			//// cfg manager: wrapper around go-cfg
			//configManager := config.NewManager(cfg)
			//
			//// nori storage
			//storage, err := noriStorage.NewStorage(cfg, log)
			//if err != nil {
			//	log.Error(err)
			//	os.Exit(1)
			//}
			//
			//// plugin manager
			//pluginManager := plugins.NewManager(
			//	storage,
			//	configManager,
			//	noriVersion,
			//	plugins.NewPluginExtractor(),
			//	log.WithField("component", "PluginManager"))
			//
			//// Load Plugins
			//dirs := getPluginsDir(cfg, log)
			//log.Infof("Plugin dir(s): %s", strings.Join(dirs, ",\n"))
			//err = pluginManager.AddDir(dirs)
			//if err != nil {
			//	log.Error(err)
			//}
			//
			//wg := &sync.WaitGroup{}
			//
			//// check
			//wg.Add(1)
			//err = pluginManager.StartAll(context.Background())
			//if err != nil {
			//	log.Error(err)
			//	os.Exit(1)
			//}
			//
			//// REST API
			//if cfg.Bool("nori.rest.enable") {
			//	rest.New(cfg.String("nori.rest.address"),
			//		cfg.String("nori.rest.base"),
			//		pluginManager,
			//		wg,
			//		ch.shutdown,
			//		log,
			//	)
			//}
			//
			//// gRPC
			//if cfg.Bool("nori.grpc.enable") {
			//	server := grpc.NewServer(
			//		dirs,
			//		cfg.String("nori.grpc.address"),
			//		true,
			//		pluginManager,
			//		wg,
			//		ch.shutdown,
			//		log)
			//	server.SetCertificates(cfg.String("nori.grpc.tls.ca"), cfg.String("nori.grpc.tls.private"))
			//	err = server.Run()
			//	if err != nil {
			//		log.Error(err)
			//		os.Exit(1)
			//	}
			//}
			//
			//go func() {
			//	<-ch.interrupt
			//	close(ch.shutdown)
			//
			//	if err := pluginManager.StopAll(context.Background()); err != nil {
			//		log.Error(err)
			//	}
			//	wg.Done()
			//}()
			//
			//wg.Wait()
			//log.Info("Nori Plugin Engine stopped")
		},
	}
}

func getPluginsDir(config go_config.Config, logger logger.Logger) []string {
	dirs := config.Slice("plugins.dir", ",")
	if len(dirs) == 0 {
		logger.Error("plugins.dir not defined or has incorrect format")
		os.Exit(1)
	}
	var list []string
	for _, d := range dirs {
		item, ok := d.(string)
		if !ok {
			continue
		}
		list = append(list, item)
	}
	return list
}
