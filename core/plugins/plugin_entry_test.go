// Copyright © 2018 Secure2Work info@secure2work.com
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package plugins

import (
	"context"
	"strings"
	"testing"

	"github.com/secure2work/nori/core/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPluginEntry_calcWeightKinds(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[entities.HTTP.String()] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[entities.Cache.String()] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	pList := SortPlugins(pMap)

	for i, p := range pList {
		assert.Equal(p.Plugin().GetMeta().GetId(), order[i])
	}

	assert.Nil(nil)
}

func TestPluginEntry_calcWeightCustom(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http", "nori/stats"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	item2 := new(MockPlugin)
	item2.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Stats",
		Id:           order[2],
		Dependencies: []string{},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[entities.HTTP.String()] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[order[2]] = &pluginEntry{
		plugin: item2,
		weight: -1,
	}

	pMap[entities.Cache.String()] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	pList := SortPlugins(pMap)

	for i, p := range pList {
		assert.Equal(p.Plugin().GetMeta().GetId(), order[i])
	}

	assert.Nil(nil)
}

func TestPluginEntry_calcWeightCustomWithDeps(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http", "nori/stats", "nori/admin", "nori/test"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	item2 := new(MockPlugin)
	item2.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Stats",
		Id:           order[2],
		Dependencies: []string{},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	item3 := new(MockPlugin)
	item3.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Admin",
		Id:           order[3],
		Dependencies: []string{"nori/stats"},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	item4 := new(MockPlugin)
	item4.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Test",
		Id:           order[4],
		Dependencies: []string{"http"},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[order[4]] = &pluginEntry{
		plugin: item4,
		weight: -1,
	}

	pMap[entities.HTTP.String()] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[order[3]] = &pluginEntry{
		plugin: item3,
		weight: -1,
	}

	pMap[order[2]] = &pluginEntry{
		plugin: item2,
		weight: -1,
	}

	pMap[entities.Cache.String()] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	pList := SortPlugins(pMap)

	for i, p := range pList {
		assert.Equal(p.Plugin().GetMeta().GetId(), order[i])
	}

	assert.Nil(nil)
}

func TestPluginEntry_checkDepsNone(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[strings.ToLower(entities.HTTP.String())] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[strings.ToLower(entities.Cache.String())] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	err := CheckDependencies(pMap)

	assert.Len(err, 0)
}

func TestPluginEntry_checkDepsOne(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http", "nori/stats"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	item2 := new(MockPlugin)
	item2.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Stats",
		Id:           order[2],
		Dependencies: []string{"cache"},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[strings.ToLower(entities.HTTP.String())] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[order[2]] = &pluginEntry{
		plugin: item2,
		weight: -1,
	}

	pMap[strings.ToLower(entities.Cache.String())] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	err := CheckDependencies(pMap)

	assert.Len(err, 0)
}

func TestPluginEntry_checkDepsWithVer(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http", "nori/stats", "nori/admin"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	item2 := new(MockPlugin)
	item2.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Stats",
		Id:           order[2],
		Dependencies: []string{"cache"},
		Interface:    entities.Custom,
		Version:      "1.1",
	})

	item3 := new(MockPlugin)
	item3.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Admin",
		Id:           order[3],
		Dependencies: []string{"nori/stats:>=1.0"},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[strings.ToLower(entities.HTTP.String())] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[order[2]] = &pluginEntry{
		plugin: item2,
		weight: -1,
	}

	pMap[order[3]] = &pluginEntry{
		plugin: item3,
		weight: -1,
	}

	pMap[strings.ToLower(entities.Cache.String())] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	err := CheckDependencies(pMap)

	assert.Len(err, 0)
}

func TestPluginEntry_checkDepsUnresolved(t *testing.T) {
	assert := assert.New(t)

	order := []string{"nori/cache/memory", "nori/http", "nori/stats"}

	item0 := new(MockPlugin)
	item0.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori HTTP",
		Id:           order[1],
		Dependencies: []string{},
		Interface:    entities.HTTP,
		Version:      "1.0",
	})

	item1 := new(MockPlugin)
	item1.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Cache",
		Id:           order[0],
		Dependencies: []string{},
		Interface:    entities.Cache,
		Version:      "1.0",
	})

	item2 := new(MockPlugin)
	item2.On("GetMeta").Return(&entities.PluginMetaStruct{
		PluginName:   "Nori Stats",
		Id:           order[2],
		Dependencies: []string{"cache", "nori/pages:>=1.0"},
		Interface:    entities.Custom,
		Version:      "1.0",
	})

	pMap := map[string]PluginEntry{}

	pMap[strings.ToLower(entities.HTTP.String())] = &pluginEntry{
		plugin: item0,
		weight: -1,
	}

	pMap[order[2]] = &pluginEntry{
		plugin: item2,
		weight: -1,
	}

	pMap[strings.ToLower(entities.Cache.String())] = &pluginEntry{
		plugin: item1,
		weight: -1,
	}

	//
	err := CheckDependencies(pMap)

	assert.Len(err, 1)
}

type MockPlugin struct {
	mock.Mock
}

func (_ MockPlugin) GetInstance() interface{} {
	return nil
}

func (m MockPlugin) GetMeta() entities.PluginMeta {
	args := m.Mock.Called()
	return args.Get(0).(entities.PluginMeta)
}

func (_ MockPlugin) Install(_ context.Context, _ PluginRegistry) error {
	return nil
}

func (_ MockPlugin) UnInstall(_ context.Context, _ PluginRegistry) error {
	return nil
}

func (_ MockPlugin) Start(_ context.Context, _ PluginRegistry) error {
	return nil
}

func (_ MockPlugin) Stop(_ context.Context) error {
	return nil
}
