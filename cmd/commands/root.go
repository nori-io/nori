/*
Copyright 2018-2020 The Nori Authors.
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

package commands

import (
	"fmt"
	"os"

	"github.com/nori-io/logger"
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
	//config := go_config.New()
	logger := logger.New(logger.SetOutWriter(os.Stderr), logger.SetJsonFormatter(""))

	//cobra.OnInitialize(initConfig(config, logger))
	//rootCmd.PersistentFlags().StringVar(&configFile, "cfg.file", "", fmt.Sprintf("defualt: `$HOME/%s`", defaultConfigFile))

	rootCmd.AddCommand(serverCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}

// initConfig reads config file and ENV variables if set.
//func initConfig(config go_config.Config, logger commonLogger.logger) func() {
//	return func() {
//		if configFile == "" {
//			// Find home directory.
//			home, err := homedir.Dirs()
//			if err != nil {
//				logger.Fatal(err.Error())
//			}
//			// build config file path
//			configFile = filepath.Join(home, defaultConfigFile)
//		}
//
//		fileSource, err := file.Source(
//			file.File{Path: configFile, Namespace: ""},
//		)
//		if err != nil {
//			logger.Fatal(err.Error())
//		}
//		config.UseSource(fileSource)
//
//		logger.Info("Using config file: %s", configFile)
//
//		config.UseSource(env.Source("NORI", ","))
//	}
//}
