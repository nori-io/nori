package graph_test

import (
	"github.com/secure2work/nori-common/meta"
	"github.com/secure2work/nori/core/plugins/dependency"
	"github.com/stretchr/testify/assert"
	"testing"
)

//1) plugin-1 -> plugin-2 -> plugin-3 (all available)
/*func TestDependencyGraph_Sort1(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()

	plugin1 := meta.Data{
		ID: meta.ID{
			ID:      "plugin1",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin2", ">=1.0, <2.0", meta.Custom},
		},
		Description: meta.Description{
			Name: "Nori HTTP Interface",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"http"},
	}

	plugin2 := meta.Data{
		ID: meta.ID{
			ID:      "plugin2",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin3", ">=1.0, <2.0", meta.Custom},
		},
		Description: meta.Description{
			Name: "NoriCMS: MySQL Driver",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"sql", "mysql"},
	}

	plugin3 := meta.Data{
		ID: meta.ID{
			ID:      "plugin3",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "Nori",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{},
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

	managerPlugin.Add(plugin2)
	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin3)

	t.Log("Порядок плагинов до сортировки:", managerPlugin)
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Erroe in sorting")
	}
	for index, _ := range pluginsSorted {
		t.Log(index+1, "-й элемент для запуска", pluginsSorted[index].ID, " ", pluginsSorted[index].Version)
	}

	var (
		index1 int
		index2 int
		index3 int
	)
	for index, value := range pluginsSorted {
		if value.ID == "plugin1" {
			index1 = index
		}

		if value.ID == "plugin2" {
			index2 = index
		}
		if value.ID == "plugin3" {
			index3 = index
		}
	}

	a.Equal(true, index3 < index2)
	a.Equal(true, index2 < index1)
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin2.ID)
	managerPlugin.Remove(plugin3.ID)

}*/

//2) plugin-1 -> plugin-2 -> plugin-3 (3rd is absent)
/*func TestDependencyGraph_Sort2(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()

	//
	plugin1 := meta.Data{
		ID: meta.ID{
			ID:      "plugin1",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin2", ">=1.0, <2.0", meta.Custom},
		},
		Description: meta.Description{
			Name: "Nori HTTP Interface",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"http"},
	}

	plugin2 := meta.Data{
		ID: meta.ID{
			ID:      "plugin2",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin3", ">=1.0, <2.0", meta.Custom},
		},
		Description: meta.Description{
			Name: "NoriCMS: MySQL Driver",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"sql", "mysql"},
	}
	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin2)


	t.Log("Порядок плагинов до сортировки:", managerPlugin)
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Erroe in sorting")
	}
	for index, _ := range pluginsSorted {
		t.Log(index+1, "-й элемент для запуска", pluginsSorted[index].ID, " ", pluginsSorted[index].Version)
	}

	var (
		index1 int
		index2 int
		index3 int
	)
	for index, value := range pluginsSorted {
		if value.ID == "plugin1" {
			index1 = index
		}

		if value.ID == "plugin2" {
			index2 = index
		}
		if value.ID == "plugin3" {
			index3 = index
		}
	}
	//index3 initialized by default value=0
	a.Equal(true, index3==0)
	a.Equal(true, index2 < index1)
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin2.ID)

}
*/

//3) plugin-1 -> plugin-2 -> plugin-3 (2nd is absent)
func TestDependencyGraph_Sort3(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()

	plugin1 := meta.Data{
		ID: meta.ID{
			ID:      "plugin1",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "NoriCMS",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin2", ">=1.0, <2.0", meta.Custom},
		},
		Description: meta.Description{
			Name: "Nori HTTP Interface",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"http"},
	}


	plugin3 := meta.Data{
		ID: meta.ID{
			ID:      "plugin3",
			Version: "1.0",
		},
		Author: meta.Author{
			Name: "Nori",
			URI:  "https://noricms.com",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{},
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


	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin3)

	t.Log("Порядок плагинов до сортировки:", managerPlugin)
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Erroe in sorting")
	}
	for index, _ := range pluginsSorted {
		t.Log(index+1, "-й элемент для запуска", pluginsSorted[index].ID, " ", pluginsSorted[index].Version)
	}

	var (
		index1 int
		index2 int
		index3 int
	)
	for index, value := range pluginsSorted {
		if value.ID == "plugin1" {
			index1 = index
		}

		if value.ID == "plugin2" {
			index2 = index
		}
		if value.ID == "plugin3" {
			index3 = index
		}
	}

	a.Equal(true, index3 < index1||index3>index1)
	//index2 initialized by default value=0
	a.Equal(true, index2==0)
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin3.ID)

}

//4) plugin -> interface (all available)