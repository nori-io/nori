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
	"os"

	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
)

const NoriMajorVersion = "1.0"

func NoriVersion(log *logrus.Logger) Version {
	return Version{
		logger:  log,
		version: "1.0.0",
	}
}

type Version struct {
	logger  *logrus.Logger
	version string
}

func (v Version) Version() *version.Version {
	ver, err := version.NewVersion(v.version)
	if err != nil {
		v.logger.Infof("Can't process Nori version [%s]", v.version)
		os.Exit(1)
	}
	return ver
}

func (v Version) Original() string {
	return v.version
}
