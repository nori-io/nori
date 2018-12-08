// Copyright © 2018 Secure2Work info@secure2work.com
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

package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"strings"

	"github.com/cheebo/go-config"
	"github.com/secure2work/nori/core/grpc"
	"github.com/secure2work/nori/core/plugins"
	"github.com/secure2work/nori/core/storage"
	"github.com/sirupsen/logrus"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()

		noriStorage := storage.GetNoriStorage(config, logger)
		if noriStorage == nil {
			logger.Error("can't create NoriStorage")
			os.Exit(1)
		}

		pluginManager := plugins.GetPluginManager(noriStorage, logger, config)

		// Load Plugins
		dirs := getPluginsDir(config, logger)
		logger.Infof("Plugin dir(s): \n- %s", strings.Join(dirs, ",\n- "))
		err := pluginManager.Load(dirs)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		// Get list of installed plugins
		installedPluginsList, err := noriStorage.GetPluginMetas()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		err = pluginManager.Run(context.Background(), installedPluginsList)
		if err != nil {
			os.Exit(1)
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

func init() {
	rootCmd.AddCommand(serverCmd)

	config.SetDefault("nori.grpc.enable", true)
	config.SetDefault("nori.grpc.address", "0.0.0.0:29876")

	serverCmd.Flags().String("pem", "server.pem", "path to pem file")
	serverCmd.Flags().String("key", "server.key", "path to key file")

	viper.BindPFlag("pem", serverCmd.Flags().Lookup("pem"))
	viper.BindPFlag("key", serverCmd.Flags().Lookup("key"))
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
