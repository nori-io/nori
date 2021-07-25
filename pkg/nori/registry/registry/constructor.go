package registry

import (
	"sync"

	"github.com/nori-io/nori/pkg/nori/domain/entity"
	"github.com/nori-io/nori/pkg/nori/domain/registry"
)

func New() registry.Registry {
	return &Registry{
		mx:      sync.Mutex{},
		plugins: []*entity.Plugin{},
	}
}
