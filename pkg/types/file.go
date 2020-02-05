package types

import (
	"github.com/nori-io/nori-common/v2/plugin"
)

type File struct {
	Plugin plugin.Plugin
	Path   string
}
