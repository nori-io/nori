package graph_test

import (
	"github.com/secure2work/nori/core/plugins/dependency"
	"github.com/secure2work/nori-common/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDependencyGraph_Sort(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	httpPlugin := meta.Data{
		ID: meta.ID{
			ID:      "nori/http",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name: "Nori HTTP Interface",
		},
		Interface: meta.HTTP,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"http"},
	}

	mysqlPlugin := meta.Data{
		ID: meta.ID{
			ID:      "nori/sql",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name: "NoriCMS: MySQL Driver",
		},
		Interface: meta.SQL,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"sql", "mysql"},
	}

	cmsPlugin := meta.Data{
		ID: meta.ID{
			ID:      "nori/cms/posts/naive",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "Nori",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			meta.SQL.Dependency("1.0.0"),
			meta.HTTP.Dependency("1.0.0"),
			//    meta.HTTPTransport.Dependency(),
		},
		Description: meta.Description{
			Name:        "NoriCMS Naive Posts Plugin",
			Description: "Naive Posts Plugin",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/",
		},
		Tags: []string{"cms", "posts", "api"},
	}

	managerPlugin.Add(cmsPlugin)
	managerPlugin.Add(httpPlugin)
	managerPlugin.Add(mysqlPlugin)

	t.Log("Порядок плагинов до сортировки:", managerPlugin)

	/*  mConfig := mocks.Config{}
	  mConfig.On("String", AnythingOfType("string")).
		Return(func(s string) string {
		  return s
		})
	  manager := config.NewManager(&mConfig)
	  manager.Register(httpPlugin)
	  manager.Register(mysqlPlugin)
	  manager.Register(cmsPlugin)
	*/

	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Erroe in sorting")
	}

	for index, _ := range pluginsSorted {
		t.Log(index+1, "-й элемент для запуска", pluginsSorted[index].ID, " ", pluginsSorted[index].Version)
	}

	var (
		index1Http  int
		index2Mysql int
		index3Cms   int
	)
	for index, value := range pluginsSorted {
		if value.ID == "nori/http" {
			index1Http = index
		}

		if value.ID == "nori/sql" {
			index2Mysql = index
		}
		if value.ID == "nori/cms/posts/naive" {
			index3Cms = index
		}
	}

	//a.Equal(true, index1Http == 0 || index1Http == 1)
	//a.Equal(true, index2Mysql == 0 || index2Mysql == 1)
	a.Equal(true, index3Cms > index1Http)
	a.Equal(true, index3Cms > index2Mysql)

}