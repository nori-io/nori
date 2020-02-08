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

package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cheebo/go-config"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/nori"
	"github.com/spf13/cobra"
)

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
				os.Exit(1)
			}
		},
	}
}
