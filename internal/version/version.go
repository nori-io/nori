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

package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/nori-io/nori-common/v2/logger"

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
