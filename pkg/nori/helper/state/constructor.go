package state

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
)

func New() *StateHelper {
	return &StateHelper{
		plugins: map[meta.ID]enum.State{},
	}
}
