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

package dependency_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nori-io/nori/internal/dependency"
	"github.com/stretchr/testify/assert"

	"github.com/nori-io/nori-common/v2/meta"

	"github.com/nori-io/nori-common/v2/mocks/mock_meta"
)

func TestManager_Add(t *testing.T) {
	a := assert.New(t)

	manager := dependency.NewManager()

	id := meta.ID{
		ID:      "nori/test",
		Version: "1.0.0",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt := mock_meta.NewMockMeta(ctrl)
	mt.EXPECT().Id().Return(id).AnyTimes()
	mt.EXPECT().GetDependencies().Return([]meta.Dependency{}).AnyTimes()

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt := mock_meta.NewMockMeta(ctrl)
	mt.EXPECT().Id().Return(id).AnyTimes()
	mt.EXPECT().GetDependencies().Return([]meta.Dependency{}).AnyTimes()

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
		Constraint: "^1.0.0",
		Interface:  "nori/mocks",
	}

	unresolved := map[meta.ID][]meta.Dependency{
		id1: []meta.Dependency{dep},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt1 := mock_meta.NewMockMeta(ctrl)
	mt1.EXPECT().Id().Return(id1).AnyTimes()
	mt1.EXPECT().GetDependencies().Return([]meta.Dependency{dep}).AnyTimes()
	mt1.EXPECT().GetInterface().Return(meta.Interface("nori/test@1.0.0")).AnyTimes()

	mt2 := mock_meta.NewMockMeta(ctrl)
	mt2.EXPECT().Id().Return(id2).AnyTimes()
	mt2.EXPECT().GetDependencies().Return([]meta.Dependency{}).AnyTimes()
	mt2.EXPECT().GetInterface().Return(meta.Interface("nori/mocks@1.0.0")).AnyTimes()

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
		ID:      "nori/main",
		Version: "1.0.0",
	}

	id2 := meta.ID{
		ID:      "nori/dependency",
		Version: "1.0.0",
	}

	dep := meta.Dependency{
		Interface:  "nori/dependency",
		Constraint: "^1.0.0",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt1 := mock_meta.NewMockMeta(ctrl)
	mt1.EXPECT().Id().Return(id1).AnyTimes()
	mt1.EXPECT().GetDependencies().Return([]meta.Dependency{dep}).AnyTimes()
	mt1.EXPECT().GetInterface().Return(meta.Interface("test/mocks@1.0.0")).AnyTimes()

	mt2 := mock_meta.NewMockMeta(ctrl)
	mt2.EXPECT().Id().Return(id2).AnyTimes()
	mt2.EXPECT().GetDependencies().Return([]meta.Dependency{}).AnyTimes()
	mt2.EXPECT().GetInterface().Return(meta.Interface("nori/dependency@1.0.0")).AnyTimes()

	a.NoError(manager.Add(mt1))
	a.NoError(manager.Add(mt2))

	list, err := manager.Sort()
	a.NoError(err)

	ordered := []meta.ID{id2, id1}

	for i, id := range list {
		a.Equal(ordered[i], id)
	}
}
