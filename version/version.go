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
