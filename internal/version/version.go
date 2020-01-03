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

package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/nori-io/nori-common/logger"

	"github.com/hashicorp/go-version"
)

const CoreVersion = "0.2.0"

var (
	// The git commit that was compiled. These will be filled in by the compiler.
	GitCommit string

	// The main version number that is being run at the moment.
	NoriVersion = "0.2.0"

	// A pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = ""
)

func Version(log logger.FieldLogger) core {
	return core{
		logger:  log,
		version: CoreVersion,
	}
}

type core struct {
	logger  logger.FieldLogger
	version string
}

func (v core) Version() *version.Version {
	ver, err := version.NewVersion(v.version)
	ver.Segments()
	if err != nil {
		v.logger.Info("Can't process Nori version [%s]", v.version)
		os.Exit(1)
	}
	return ver
}

func (v core) Original() string {
	return v.version
}

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func GetHumanVersion() string {
	version := NoriVersion
	release := VersionPrerelease
	if release == "" {
		release = "dev"
	}

	if release != "" {
		if !strings.HasSuffix(version, "-"+release) {
			version += fmt.Sprintf("-%s", release)
		}
		if GitCommit != "" {
			version += fmt.Sprintf(" (%s)", GitCommit)
		}
	}
	return strings.Replace(version, "'", "", -1)
}
