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
	"path/filepath"

	"github.com/nori-io/logger"
	commonLogger "github.com/nori-io/nori-common/logger"

	"github.com/cheebo/go-config"
	"github.com/cheebo/go-config/sources/env"
	"github.com/cheebo/go-config/sources/file"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfgFile string

const (
	configDir  = ".config/nori"
	configName = "nori.json"
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
	logger := logger.New()

	cobra.OnInitialize(initConfig(config, logger))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s/%s)", configDir, configName))

	rootCmd.AddCommand(serverCmd(config, logger), versionCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(config go_config.Config, logger commonLogger.Logger) func() {
	return func() {
		config.SetDefault("nori.grpc.enable", true)
		config.SetDefault("nori.grpc.address", "0.0.0.0:29876")

		config.SetDefault("nori.rest.enable", true)
		config.SetDefault("nori.rest.base", "/")
		config.SetDefault("nori.rest.address", "0.0.0.0:28541")

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
			file.File{Path: cfgFile, Type: go_config.JSON, Namespace: ""},
		)
		if err != nil {
			logger.Fatal(err.Error())
		}
		config.UseSource(fileSource)

		logger.Info("Using config file: %s \n", cfgFile)

		config.UseSource(env.Source("NORI"))
	}
}
