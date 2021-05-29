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
	"github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/app"
	"github.com/nori-io/nori/internal/app/container"
	"github.com/spf13/cobra"
)

var (
	// Cmd server command
	serverCmd = func() *cobra.Command {
		var configFile string

		cmd := &cobra.Command{
			Use:           "server",
			Short:         "server",
			SilenceUsage:  true,
			SilenceErrors: true,
			Run: func(c *cobra.Command, v []string) {
				container := container.New(configFile)
				err := container.Invoke(func(app *app.App) {
					app.Run()
				})
				if err != nil {
					logger.L().Fatal(err.Error())
				}
			},
		}

		cmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")

		return cmd
	}()
)
