package registry

import (
	"sync"

	"github.com/nori-io/common/v5/pkg/domain/plugin"
)

func New() *Registry {
	return &Registry{
		mx:      sync.Mutex{},
		plugins: []plugin.Plugin{},
	}
}
