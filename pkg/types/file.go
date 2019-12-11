package types

import (
	"github.com/nori-io/nori-common/plugin"
)

type File struct {
	Plugin plugin.Plugin
	Path   string
}
