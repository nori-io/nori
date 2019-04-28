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
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nori-io/nori-common/logger"

	"github.com/nori-io/nori/core/grpc"

	"github.com/nori-io/nori/core/client/rest"

	"github.com/nori-io/nori/version"

	"strings"

	"github.com/spf13/cobra"

	"github.com/cheebo/go-config"
	noriConfig "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori/core/config"
	"github.com/nori-io/nori/core/plugins"
	noriStorage "github.com/nori-io/nori/core/storage"
)

type app struct {
	pluginManager plugins.Manager
	configManager noriConfig.Manager
	storage       noriStorage.Storage
}

type channels struct {
	shutdown  chan struct{}
	interrupt chan os.Signal
}

// serverCmd represents the server command
func serverCmd(goConfig go_config.Config, logger logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Run: func(cmd *cobra.Command, args []string) {
			// shutdown and interrupt
			ch := channels{
				shutdown:  make(chan struct{}, 1),
				interrupt: make(chan os.Signal, 1),
			}
			signal.Notify(ch.interrupt, syscall.SIGTERM)
			signal.Notify(ch.interrupt, syscall.SIGINT)

			noriVersion := version.NoriVersion(logger)

			logger.Infof("Nori Engine [version %s]", noriVersion.Version().String())

			// nori storage
			storage, err := noriStorage.NewStorage(goConfig, logger)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// goConfig manager: wrapper around go-goConfig
			configManager := config.NewManager(goConfig)

			// plugin manager
			pluginManager := plugins.NewManager(
				storage,
				configManager,
				noriVersion,
				plugins.NewPluginExtractor(),
				logger.WithField("component", "PluginManager"))

			// Load Plugins
			dirs := getPluginsDir(goConfig, logger)
			logger.Infof("Plugin dir(s): %s", strings.Join(dirs, ",\n"))
			err = pluginManager.AddDir(dirs)
			if err != nil {
				logger.Error(err)
			}

			wg := &sync.WaitGroup{}

			// check
			wg.Add(1)
			err = pluginManager.StartAll(context.Background())
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// REST API
			if goConfig.Bool("nori.rest.enable") {
				rest.New(goConfig.String("nori.rest.address"),
					goConfig.String("nori.rest.base"),
					pluginManager,
					wg,
					ch.shutdown,
					logger,
				)
			}

			// gRPC
			if goConfig.Bool("nori.grpc.enable") {
				server := grpc.NewServer(
					dirs,
					goConfig.String("nori.grpc.address"),
					true,
					pluginManager,
					wg,
					ch.shutdown,
					logger)
				server.SetCertificates(goConfig.String("nori.grpc.tls.ca"), goConfig.String("nori.grpc.tls.private"))
				err = server.Run()
				if err != nil {
					logger.Fatal(err)
				}
			}

			go func() {
				<-ch.interrupt
				close(ch.shutdown)

				if err := pluginManager.StopAll(context.Background()); err != nil {
					logger.Error(err)
				}
				wg.Done()
			}()

			wg.Wait()
			logger.Info("Nori Plugin Engine stopped")
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
