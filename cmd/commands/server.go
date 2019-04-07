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
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"

	"github.com/gorilla/mux"

	"github.com/nori-io/nori/version"

	"strings"

	"github.com/spf13/cobra"

	"github.com/cheebo/go-config"
	"github.com/nori-io/nori/core/client/rest"
	configManager "github.com/nori-io/nori/core/config"
	"github.com/nori-io/nori/core/grpc"
	"github.com/nori-io/nori/core/plugins"
	"github.com/nori-io/nori/core/storage"
	"github.com/sirupsen/logrus"
)

// serverCmd represents the server command
func serverCmd(config go_config.Config, logger *logrus.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Run: func(cmd *cobra.Command, args []string) {
			noriVersion := version.NoriVersion(logger)

			logger.Infof("Nori Engine [version %s]", noriVersion.Version().String())

			// nori storage
			storage, err := storage.NewStorage(config, logger)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			// config manager: wrapper around go-config
			configManager := configManager.NewManager(config)

			// plugin manager
			pluginManager := plugins.NewManager(
				storage,
				configManager,
				noriVersion,
				plugins.NewPluginExtractor(),
				logger.WithField("component", "PluginManager").Logger)

			// Load Plugins
			dirs := getPluginsDir(config, logger)
			logger.Infof("Plugin dir(s): %s", strings.Join(dirs, ",\n"))
			err = pluginManager.AddDir(dirs)
			if err != nil {
				logger.Error(err)
			}

			// check
			err = pluginManager.StartAll(context.Background())
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			if config.Bool("nori.rest.enable") {
				go func(m plugins.Manager) {
					addr := config.String("nori.rest.address")
					logger.Infof("Starts Nori Core REST API on %s", addr)
					r := mux.NewRouter()
					rest.RegisterRoutes(config, r, m)

					box := packr.NewBox("../../html")
					fs := http.FileServer(box)
					r.Handle("/", fs)

					server := &http.Server{
						Addr:    addr,
						Handler: r,
					}
					err := server.ListenAndServe()
					if err != nil {
						logger.Error(err)
					}

				}(pluginManager)
			}

			// gRPC Config
			if config.Bool("nori.grpc.enable") {
				addr := config.String("nori.grpc.address")
				server := grpc.NewServer(
					dirs,
					addr,
					true,
					pluginManager,
					logger)
				server.SetCertificates(config.String("nori.grpc.tls.ca"), config.String("nori.grpc.tls.private"))
				err = server.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}
}

func getPluginsDir(config go_config.Config, logger *logrus.Logger) []string {
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
