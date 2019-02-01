package graph_test

import (
	"github.com/secure2work/nori-common/meta"
	"github.com/secure2work/nori/core/plugins/dependency"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	managerPlugin = dependency.NewManager()
	// depend of plugin2
	plugin1 = meta.Data{
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
	//depend of plugin3
	plugin2 = meta.Data{
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
	// without dependencies
	plugin3 = meta.Data{
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
	// without dependencies
	plugin4 = meta.Data{
		ID: meta.ID{
			ID:      "plugin4",
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
	// without dependencies
	pluginHttp = meta.Data{
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
	//without dependencies
	pluginMysql = meta.Data{
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
	// depend of  pluginHttp, pluginMysql
	pluginCms = meta.Data{
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
)

//1) plugin1 -> plugin2 -> plugin3 (all available) order for adding - 1 3 2
func TestDependencyGraph_Sort1(t *testing.T) {
	a := assert.New(t)
	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin3)
	managerPlugin.Add(plugin2)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
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
	a.Equal(3, len(pluginsSorted))
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin2.ID)
	managerPlugin.Remove(plugin3.ID)

}

//2) plugin1 -> plugin2 -> plugin3 (3rd is unavailable)
func TestDependencyGraph_Sort2(t *testing.T) {
	a := assert.New(t)
	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin2)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}
	a.Equal(0, len(pluginsSorted))
	a.NotEqual(err, nil)
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin2.ID)
}

//3) plugin1 -> plugin2 -> plugin3 (2nd is unavailable)
func TestDependencyGraph_Sort3(t *testing.T) {
	a := assert.New(t)
	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin3)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}
	a.Equal(0, len(pluginsSorted))
	a.NotEqual(err, nil)
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin3.ID)
}

//4) plugin1 -> interfaceHttp (all available)
func TestDependencyGraph_Sort4(t *testing.T) {
	a := assert.New(t)
	plugin1.Dependencies = []meta.Dependency{
		meta.HTTP.Dependency("1.0"),
	}
	managerPlugin.Add(plugin1)
	managerPlugin.Add(pluginHttp)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}

	var (
		index1    int
		indexHttp int
	)
	for index, value := range pluginsSorted {
		if value.ID == "plugin1" {
			index1 = index
		}
		if value.ID == "http" {
			indexHttp = index
		}
	}

	a.Equal(true, indexHttp < index1)
	a.Equal(2, len(pluginsSorted))

	plugin1.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(pluginHttp.ID)

}

//5) plugin1-> interfaceHttp (interface is unavailable)
func TestDependencyGraph_Sort5(t *testing.T) {
	a := assert.New(t)
	plugin1.Dependencies = []meta.Dependency{
		meta.HTTP.Dependency("1.0"),
	}

	managerPlugin.Add(plugin1)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}
	for index, _ := range pluginsSorted {
		t.Log(index+1, " element in list for start:", pluginsSorted[index].ID, " ", pluginsSorted[index].Version)
	}

	var (
		index1 int
	)
	for index, value := range pluginsSorted {
		if value.ID == "plugin1" {
			index1 = index
		}
	}

	a.Equal(true, index1 == 0)
	a.NotEqual(err, nil)
	a.Equal(0, len(pluginsSorted))
	plugin1.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}
	managerPlugin.Remove(plugin1.ID)

}

//6) plugin1 -> plugin2, plugin 3 -> plugin2, plugin 2 -> plugin4 (all available)
func TestDependencyGraph_Sort6(t *testing.T) {
	a := assert.New(t)
	plugin2.Dependencies = []meta.Dependency{
		{"plugin4", ">=1.0, <2.0", meta.Custom},
	}
	plugin3.Dependencies = []meta.Dependency{
		{"plugin4", ">=1.0, <2.0", meta.Custom},
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
		Dependencies: []meta.Dependency{
			{"plugin2", ">=1.0, <2.0", meta.Custom},
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

	managerPlugin.Add(plugin1)
	managerPlugin.Add(plugin2)
	managerPlugin.Add(plugin3)
	managerPlugin.Add(plugin4)

	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}

	var (
		index1 int
		index2 int
		index3 int
		index4 int
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
		if value.ID == "plugin4" {
			index4 = index
		}
	}

	a.Equal(true, index4 < index2)
	a.Equal(true, index2 < index3)
	a.Equal(true, index2 < index1)
	a.Equal(4, len(pluginsSorted))
	plugin2.Dependencies = []meta.Dependency{
		{"plugin3", ">=1.0, <2.0", meta.Custom},
	}
	plugin3.Dependencies = []meta.Dependency{}
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin2.ID)
	managerPlugin.Remove(plugin3.ID)
	managerPlugin.Remove(plugin4.ID)

}

//7) plugin1 -> plugin2, plugin 3 -> plugin2, (plugin 2 is unavailable)
func TestDependencyGraph_Sort7(t *testing.T) {
	a := assert.New(t)
	plugin3.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}
	plugin4 := meta.Data{
		ID: meta.ID{
			ID:      "plugin4",
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
	managerPlugin.Add(plugin4)

	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}

	a.NotEqual(nil, err)
	a.Equal(0, len(pluginsSorted))
	plugin3.Dependencies = []meta.Dependency{}
	managerPlugin.Remove(plugin1.ID)
	managerPlugin.Remove(plugin3.ID)
	managerPlugin.Remove(plugin4.ID)

}

//8) pluginCms->pluginMysql, pluginCms->pluginHttp
func TestDependencyGraph_Sort8(t *testing.T) {
	a := assert.New(t)
	managerPlugin.Add(pluginCms)
	managerPlugin.Add(pluginHttp)
	managerPlugin.Add(pluginMysql)
	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
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
	a.Equal(true, index3Cms > index1Http)
	a.Equal(true, index3Cms > index2Mysql)
	a.Equal(3, len(pluginsSorted))
	managerPlugin.Remove(pluginCms.ID)
	managerPlugin.Remove(pluginHttp.ID)
	managerPlugin.Remove(pluginMysql.ID)

}

//9) ring -plugin1->plugin1
func TestDependencyGraph_Sort9(t *testing.T) {
	a := assert.New(t)
	plugin1.Dependencies = []meta.Dependency{
		{"plugin1", ">=1.0, <2.0", meta.Custom},
	}

	managerPlugin.Add(plugin1)

	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}

	a.Equal(err, nil)
	plugin1.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}
	managerPlugin.Remove(plugin1.ID)

}

//10)ring plugin2->plugin3, plugin3->plugin2
func TestDependencyGraph_Sort10(t *testing.T) {
	a := assert.New(t)
	plugin3.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}

	managerPlugin.Add(plugin2)
	managerPlugin.Add(plugin3)

	t.Log("Plugins' order until sorting:")
	for index, value := range managerPlugin.GetDependencyGraph() {
		t.Log("Plugin n.", index+1, " in list until sotring:", value.ID, " version: ", value.Version, " Dependencies:", value.String())
	}
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting")
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}

	a.NotEqual(err, nil)
	plugin2.Dependencies = []meta.Dependency{
		{"plugin2", ">=1.0, <2.0", meta.Custom},
	}
	plugin3.Dependencies = []meta.Dependency{}
	managerPlugin.Remove(plugin2.ID)
	managerPlugin.Remove(plugin3.ID)

}