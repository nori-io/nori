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
	"github.com/nori-io/nori-common/v2/logger"
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
				os.Exit(1)
			}
		},
	}
}
