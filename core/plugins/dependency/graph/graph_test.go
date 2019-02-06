package graph_test

import (
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/core/plugins/dependency"
	"github.com/stretchr/testify/assert"
	"testing"
)

// depend of plugin2
func plugin1(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "plugin1",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin2", ">=1.0, <2.0", meta.Custom},
		},
		Interface: meta.Custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

//depend of plugin3
func plugin2(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "plugin2",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Dependencies: []meta.Dependency{
			{"plugin3", ">=1.0, <2.0", meta.Custom},
		},
		Interface: meta.Custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// without dependencies
func plugin3(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "plugin3",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Interface: meta.Custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

//without dependencies
func plugin4(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "plugin4",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Interface: meta.Custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// without dependencies
func pluginHttp(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginHttp",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Interface: meta.HTTP,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// without dependencies
func pluginMysql(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginMysql",
			Version: "1.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Interface: meta.SQL,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// depend of  pluginHttp, pluginMysql
func pluginCms(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginCms",
			Version: "1.0",
		},
		Dependencies: []meta.Dependency{meta.HTTP.Dependency("1.0"), meta.SQL.Dependency("1.0")},
		Core: meta.Core{
			VersionConstraint: ">=1.0, <2.0",
		},
		Interface: meta.Custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

//1) plugin1 -> plugin2 -> plugin3 (all available) order for adding - 1 3 2
func TestDependencyGraph_Sort1(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1())
	managerPlugin.Add(plugin3())
	managerPlugin.Add(plugin2())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting:", err.Error())
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
}

//2) plugin1 -> plugin2 -> plugin3 (3rd is unavailable)
func TestDependencyGraph_Sort2(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1())
	managerPlugin.Add(plugin2())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)
}

//3) plugin1 -> plugin2 -> plugin3 (2nd is unavailable)
func TestDependencyGraph_Sort3(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1())
	managerPlugin.Add(plugin3())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)
}

//4) plugin1 -> interfaceHttp (all available)
func TestDependencyGraph_Sort4(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(
		meta.HTTP.Dependency("1.0")))
	managerPlugin.Add(pluginHttp())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting:", err.Error())
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
}

//5) plugin1-> interfaceHttp (interface is unavailable)
func TestDependencyGraph_Sort5(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.HTTP.Dependency("1.0")))
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)

}

//6) plugin1 -> plugin2, plugin 3 -> plugin2, plugin 2 -> plugin4 (all available)
func TestDependencyGraph_Sort6(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1())
	managerPlugin.Add(plugin2(meta.Dependency{"plugin4", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin3(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin4())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting:", err.Error())
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
	a.Equal(err, nil)
	a.Equal(true, index4 < index2)
	a.Equal(true, index2 < index3)
	a.Equal(true, index2 < index1)
	a.Equal(4, len(pluginsSorted))
}

//7) plugin1 -> plugin2, plugin 3 -> plugin2, (plugin 2 is unavailable)
func TestDependencyGraph_Sort7(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1())
	managerPlugin.Add(plugin3(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom}))

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)

}

//8) pluginCms->pluginMysql, pluginCms->pluginHttp
func TestDependencyGraph_Sort8(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(pluginCms())
	managerPlugin.Add(pluginHttp())
	managerPlugin.Add(pluginMysql())
	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting:", err.Error())
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
		if value.ID == "pluginHttp" {
			index1Http = index
		}
		if value.ID == "pluginMysql" {
			index2Mysql = index
		}
		if value.ID == "pluginCms" {
			index3Cms = index
		}
	}
	a.Equal(true, index3Cms > index1Http)
	a.Equal(true, index3Cms > index2Mysql)
	a.Equal(3, len(pluginsSorted))
}

//9) ring -plugin1->plugin1
func TestDependencyGraph_Sort9(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.Dependency{"plugin1", ">=1.0, <2.0", meta.Custom}))

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)
}

//10)ring plugin2->plugin3, plugin3->plugin2
func TestDependencyGraph_Sort10(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin2())
	managerPlugin.Add(plugin3(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom}))

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)

}

//11) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHttp and pluginMysql -> plugin3
func TestDependencyGraph_Sort11(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom},
		meta.HTTP.Dependency("1.0")))
	managerPlugin.Add(pluginHttp(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin3())
	managerPlugin.Add(plugin2(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginMysql(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginCms())

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	pluginsSorted, err := managerPlugin.Sort()
	if err != nil {
		t.Log("Error in sorting:", err.Error())
	}
	t.Log("Plugins' order after sorting:")
	for index, _ := range pluginsSorted {
		t.Log("Plugin n.", index+1, " in list for start:", pluginsSorted[index].String())
	}
	var (
		index1 int
		index2 int
		index3 int
		indexHttp int
		indexMysql int
		indexCms int
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
		if value.ID=="pluginCms"{
			indexCms=index
		}
		if value.ID=="pluginHttp"{
			indexHttp=index
		}
		if value.ID=="pluginMysql"{
			indexMysql=index
		}
	}
	a.Equal(true,index3<indexMysql)
	a.Equal(true,index3<indexHttp)
	a.Equal(true,indexHttp<indexCms)
	a.Equal(true,indexMysql<indexCms)
	a.Equal(true, index3 < index2)
	a.Equal(true, index2 < index1)
	a.Equal(true,indexHttp<index1)

	a.Equal(6, len(pluginsSorted))
}

//as 11 test but have ring pluginCms->pluginCms
//12) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHttp and pluginMysql -> plugin3,
func TestDependencyGraph_Sort12(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom},
		meta.HTTP.Dependency("1.0")))
	managerPlugin.Add(pluginHttp(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin3())
	managerPlugin.Add(plugin2())
	managerPlugin.Add(pluginMysql(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginCms(meta.Dependency{"pluginCms", ">=1.0, <2.0", meta.Custom}))

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)

}

//as 11 test but have ring pluginHttp->plugin3, plugin3->pluginHttp
//13) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHttp and pluginMysql -> plugin3
func TestDependencyGraph_Sort13(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom},
		meta.HTTP.Dependency("1.0")))
	managerPlugin.Add(pluginHttp(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin3(meta.Dependency{"pluginHttp", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin2(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginMysql(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginCms())

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)
}

//// ring through interface as 11 test but have ring pluginCms>interfaceHttp, plugin3->pluginCms
//14) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHttp and pluginMysql -> plugin3
func TestDependencyGraph_Sort14(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	managerPlugin.Add(plugin1(meta.Dependency{"plugin2", ">=1.0, <2.0", meta.Custom},
		meta.HTTP.Dependency("1.0")))
	managerPlugin.Add(pluginHttp(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin3(meta.Dependency{"pluginCms",">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(plugin2(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginMysql(meta.Dependency{"plugin3", ">=1.0, <2.0", meta.Custom}))
	managerPlugin.Add(pluginCms())

	t.Log("Plugins' order until sorting:")
	pluginsList := managerPlugin.GetPluginsList()
	i := 0
	for _, value := range pluginsList {
		i++
		if len(value.GetDependencies()) > 0 {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), " Dependencies:")
			j := 0
			for _, depvalue := range value.GetDependencies() {
				j++
				t.Log("Dependence n.", j, "for", value.Id().ID, "is", depvalue.String())
			}
		} else {
			t.Log("Plugin n.", i, " in list until sotring:", value.Id(), "Plugin doesn't have dependencies")
		}
	}
	t.Log()
	_, err := managerPlugin.Sort()
	a.Error(err, "Error in sorting")
	t.Log(err)
}


