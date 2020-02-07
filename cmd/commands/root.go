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
	"fmt"
	"os"
	"path/filepath"

	"github.com/nori-io/logger"
	commonLogger "github.com/nori-io/nori-common/v2/logger"

	"github.com/cheebo/go-config"
	"github.com/cheebo/go-config/sources/env"
	"github.com/cheebo/go-config/sources/file"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

const (
	defaultConfigFile = ".config/nori/nori.yml"
)

// root command
var rootCmd = &cobra.Command{
	Use:   "nori [command]",
	Short: fmt.Sprintf(`todo Version()`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	config := go_config.New()
	logger := logger.New(logger.SetOutWriter(os.Stderr), logger.SetJsonFormatter(""))

	cobra.OnInitialize(initConfig(config, logger))
	rootCmd.PersistentFlags().StringVar(&configFile, "cfg.file", "", fmt.Sprintf("defualt: `$HOME/%s`", defaultConfigFile))

	rootCmd.AddCommand(serverCmd(config, logger), versionCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

func initConf(config go_config.Config, logger commonLogger.Logger) func() {
	return func() {
		// 1. check -config flag, if empty then
		// 2. check NORI_CONFIG env variable, if empty then
		// 3. load default config

		cfg := "file,/home/sergey/nori"

		// use env variables as config source
		config.UseSource(env.Source("NORI", ","))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(config go_config.Config, logger commonLogger.Logger) func() {
	return func() {
		//config.SetDefault("nori.grpc.enable", true)
		//config.SetDefault("nori.grpc.address", "0.0.0.0:29876")

		//config.SetDefault("nori.rest.enable", true)
		//config.SetDefault("nori.rest.base", "/")
		//config.SetDefault("nori.rest.address", "0.0.0.0:28541")

		if cfgFile == "" {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				logger.Fatal(err.Error())
			}
			// build config file path
			cfgFile = filepath.Join(home, configDir, configName)
		}

		fileSource, err := file.Source(
			file.File{Path: cfgFile, Namespace: ""},
		)
		if err != nil {
			logger.Fatal(err.Error())
		}
		config.UseSource(fileSource)

		logger.Info("Using config file: %s", cfgFile)

		config.UseSource(env.Source("NORI", ","))
	}
}
