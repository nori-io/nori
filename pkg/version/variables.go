package version

import (
	"runtime/debug"
	"strings"
)

var (
	GOOS      string
	GOARCH    string
	GOVERSION string
)

func GetCommonPkgVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	for _, dep := range buildInfo.Deps {
		if strings.HasPrefix(dep.Path, "github.com/nori-io/common") {
			if dep.Version == "" {
				return ""
			}
			return dep.Version[1:]
		}
	}
	return ""
}