package dependency_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nori-io/nori-common/meta"

	"github.com/nori-io/nori/core/plugins/dependency"
	metaMock "github.com/nori-io/nori/core/plugins/mocks"
)

func TestManager_Add(t *testing.T) {
	a := assert.New(t)

	manager := dependency.NewManager()

	id := meta.ID{
		ID:      "nori/test",
		Version: "1.0.0",
	}

	mt := &metaMock.Meta{}
	mt.On("Id").Return(id)
	mt.On("GetDependencies").Return([]meta.Dependency{})

	err := manager.Add(mt)
	a.NoError(err)

	a.True(manager.Has(id))
}

func TestManager_Remove(t *testing.T) {
	a := assert.New(t)

	manager := dependency.NewManager()

	id := meta.ID{
		ID:      "nori/test",
		Version: "1.0.0",
	}

	mt := &metaMock.Meta{}
	mt.On("Id").Return(id)
	mt.On("GetDependencies").Return([]meta.Dependency{})

	err := manager.Add(mt)
	a.NoError(err)

	manager.Remove(id)

	a.False(manager.Has(id))
}

func TestManager_Resolve(t *testing.T) {
	a := assert.New(t)

	manager := dependency.NewManager()

	id1 := meta.ID{
		ID:      "nori/test",
		Version: "1.0.0",
	}

	id2 := meta.ID{
		ID:      "nori/mocks",
		Version: "1.0.0",
	}

	dep := meta.Dependency{
		ID:         "nori/mocks",
		Constraint: ">=1.0.0",
	}

	unresolved := map[meta.ID][]meta.Dependency{
		id1: []meta.Dependency{dep},
	}

	mt1 := &metaMock.Meta{}
	mt1.On("Id").Return(id1)
	mt1.On("GetDependencies").Return([]meta.Dependency{dep})

	mt2 := &metaMock.Meta{}
	mt2.On("Id").Return(id2)
	mt2.On("GetDependencies").Return([]meta.Dependency{})

	a.NoError(manager.Add(mt1))

	// check unresolved
	a.True(reflect.DeepEqual(unresolved, manager.UnResolvedDependencies()))

	a.NoError(manager.Add(mt2))

	// check unresolved
	a.True(reflect.DeepEqual(map[meta.ID][]meta.Dependency{}, manager.UnResolvedDependencies()))

	// check resolved
	id, err := manager.Resolve(dep)
	a.NoError(err)
	a.Equal(id2, id)
}

func TestManager_Sort(t *testing.T) {
	a := assert.New(t)

	manager := dependency.NewManager()

	id1 := meta.ID{
		ID:      "nori/test",
		Version: "1.0.0",
	}

	id2 := meta.ID{
		ID:      "nori/mocks",
		Version: "1.0.0",
	}

	dep := meta.Dependency{
		ID:         "nori/mocks",
		Constraint: ">=1.0.0",
	}

	mt1 := &metaMock.Meta{}
	mt1.On("Id").Return(id1)
	mt1.On("GetDependencies").Return([]meta.Dependency{dep})

	mt2 := &metaMock.Meta{}
	mt2.On("Id").Return(id2)
	mt2.On("GetDependencies").Return([]meta.Dependency{})

	a.NoError(manager.Add(mt1))
	a.NoError(manager.Add(mt2))

	list, err := manager.Sort()
	a.NoError(err)

	ordered := []meta.ID{id2, id1}

	for i, id := range list {
		a.Equal(ordered[i], id)
	}
}
