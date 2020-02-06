package graph_test

import (
	"fmt"
	"testing"

	"github.com/nori-io/nori-common/v2/meta"
	"github.com/stretchr/testify/assert"

	"github.com/nori-io/nori/internal/dependency"
	"github.com/nori-io/nori/pkg/errors"
)

const (
	pluginOne     = "plugin1"
	pluginTwo     = "plugin2"
	pluginThree   = "plugin3"
	pluginFour    = "plugin4"
	pluginRingOne = "pluginRingOne"
	pluginRingTwo = "pluginRingTwo"
	pluginHttp="pluginHttp"

	InterfaceHttp  = meta.Interface("nori/Http@1.0.0")
	SQLInterface   = meta.Interface("nori/Sql@1.0.0")
	Authentication = meta.Interface("nori/Authentication@1.0.0")
	InterfaceOne=meta.Interface("nori/InterfaceOne@1.0.0")
	InterfaceTwo=meta.Interface("nori/InterfaceTwo@1.0.0")
	InterfaceThree=meta.Interface("nori/InterfaceThree@1.0.0")
	InterfaceFour=meta.Interface("nori/InterfaceFour@1.0.0")


	RingOne = meta.Interface("nori/RingOne@1.0.0")
	RingTwo = meta.Interface("nori/RingTwo@1.0.0")
)

// plugin with SelfRing
func plugin_RingOne(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      pluginRingOne,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{{Constraint: ">=1.0.0, <2.0.0", Interface: RingOne}},
		Interface:   RingOne,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// plugin with SelfRing
func plugin_RingTwo(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      pluginRingTwo,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{{Constraint:">=1.0.0, <2.0.0" , Interface:RingTwo }},
		Interface:    RingTwo,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

// plugin1 with AuthenticationInterface depends on HttpInterface
func plugin1(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      pluginOne,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{{Constraint: ">=1.0.0, <2.0.0", Interface: InterfaceHttp}},
		Interface:    InterfaceOne,
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
			ID:      pluginTwo,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{
			{Constraint:">=1.0.0, <2.0.0", Interface:InterfaceThree},
		},
		Interface: InterfaceTwo,
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
			ID:      pluginThree,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Interface: InterfaceThree,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}

func pluginHTTP(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginHTTP",
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Interface: InterfaceHttp,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}


//without dependencies
/*func plugin4(deps ...meta.Dependency) meta.Meta {
	custom := meta.Interface("nori/Custom@0.0.1")
	data := meta.Data{
		ID: meta.ID{
			ID:      pluginFour,
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Interface: custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}
*/
// without dependencies

// without dependencies
/*func pluginMysql(deps ...meta.Dependency) meta.Meta {
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginMysql",
			Version: "1.0.0",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Interface: SQLInterface,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}
*/
// depend of  pluginHTTP, pluginMysql
/*func pluginCms(deps ...meta.Dependency) meta.Meta {
	custom := meta.Interface("nori/Custom@1.0.0")
	data := meta.Data{
		ID: meta.ID{
			ID:      "pluginCms",
			Version: "1.0.0",
		},
		Dependencies: []meta.Dependency{
			HttpInterface.Dependency(),
			SQLInterface.Dependency(),
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Interface: custom,
	}
	if len(deps) > 0 {
		data.Dependencies = deps
	}
	return data
}
*/
//1) plugin1 -> plugin2 -> plugin3 (all available) order for adding - 1 3 2
func TestDependencyGraph_AllPluginsAvailable(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin3()))
	a.Nil(managerPlugin.Add(plugin2()))
	a.Nil(managerPlugin.Add(pluginHTTP()))
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
	)
	for index, value := range pluginsSorted {
		if value.ID == pluginOne {
			index1 = index
		}
		if value.ID == pluginTwo {
			index2 = index
		}
		if value.ID == pluginThree {
			index3 = index
		}
		if value.ID==pluginHttp{
			indexHttp=index
		}
	}
	fmt.Println(indexHttp)
	a.Equal(true, (indexHttp==0)||(indexHttp==1)||(indexHttp==2))
	a.Equal(true, (index3==0)||(index3==1)||(index3==2))
	a.Equal(true, (index2==1)||(index2==2)||(index2==3))
	a.Equal(true, (index1==1)||(index1==2)||(index1==3))
	a.Equal(4, len(pluginsSorted))
}

//2) plugin1 -> plugin2 -> plugin3 (3rd is unavailable)
func TestDependencyGraph_UnavailablePlugin3(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin2()))
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
func TestDependencyGraph_UnavailablePlugin2(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin3()))
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
/*func TestDependencyGraph_AllPluginsAvailableWithInterfaceDependency(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(HttpInterface.Dependency())))
	a.Nil(managerPlugin.Add(pluginHTTP()))
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
			t.Log("Plugin n.", i, " in list until sorting:", value.Id(), "Plugin doesn't have dependencies")
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
		indexHTTP int
	)
	for index, value := range pluginsSorted {
		if value.ID == pluginOne {
			index1 = index
		}
		if value.ID == "http" {
			indexHTTP = index
		}
	}

	a.Equal(true, indexHTTP < index1)
	a.Equal(2, len(pluginsSorted))
}
*/
//5) plugin1-> interfaceHttp (interface is unavailable)
/*func TestDependencyGraph_UnavailableInterface(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(HttpInterface.Dependency())))
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
*/
//6) plugin1 -> plugin2, plugin 3 -> plugin2, plugin 2 -> plugin4 (all available)
/*func TestDependencyGraph_AllPluginsAvailable2(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin2(meta.Dependency{pluginFour, ">=1.0.0, <2.0.0", meta.Interface("")})))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", meta.Interface("")})))
	a.Nil(managerPlugin.Add(plugin4()))
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
		if value.ID == pluginOne {
			index1 = index
		}
		if value.ID == pluginTwo {
			index2 = index
		}
		if value.ID == pluginThree {
			index3 = index
		}
		if value.ID == pluginFour {
			index4 = index
		}
	}
	a.Equal(err, nil)
	a.Equal(true, index4 < index2)
	a.Equal(true, index2 < index3)
	a.Equal(true, index2 < index1)
	a.Equal(4, len(pluginsSorted))
}
*/
//7) plugin1 -> plugin2, plugin 3 -> plugin2, (plugin 2 is unavailable)
/*func TestDependencyGraph_UnavailablePlugin(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", meta.Interface("")})))

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
*/
//8) pluginCms->pluginMysql, pluginCms->pluginHTTP
/*func TestDependencyGraph_PluginsCmsMySqlHttp(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(pluginCms()))
	a.Nil(managerPlugin.Add(pluginHTTP()))
	a.Nil(managerPlugin.Add(pluginMysql()))
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
		if value.ID == "pluginHTTP" {
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
*/
//9) ring -plugin1->plugin1, plugin2->plugin2
func TestDependencyGraph_Loops(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Equal(errors.LoopVertexFound{Dependency: struct {
		Constraint string
		Interface  meta.Interface
	}{Constraint: ">=1.0.0, <2.0.0", Interface: RingOne}}, managerPlugin.Add(plugin_RingOne(meta.Dependency{Constraint:">=1.0.0, <2.0.0", Interface:RingOne})))

		a.Equal(errors.LoopVertexFound{Dependency: struct {
		Constraint string
		Interface  meta.Interface
	}{Constraint: ">=1.0.0, <2.0.0", Interface: RingTwo}}, managerPlugin.Add(plugin_RingTwo(meta.Dependency{Constraint:">=1.0.0, <2.0.0", Interface:RingTwo})))

}

//10)ring plugin2->plugin3, plugin3->plugin2
/*func TestDependencyGraph_Ring(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin2()))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", meta.Interface("")})))
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
*/
//10)ring plugin1->plugin2, plugin2->plugin3, plugin3->plugin2, plugin3->plugin1
/*func TestDependencyGraph_Ring2(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin2()))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", meta.Interface("")},
		meta.Dependency{
			ID:         pluginOne,
			Constraint: ">=1.0.0, <2.0.0",
			Interface:  meta.Interface(""),
		})))

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
*/
//11) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql
// (plugins with such interfaces added),pluginHTTP and pluginMysql -> plugin3, pluginHTTP->interface SQL
/*func TestDependencyGraph_Sort1(t *testing.T) {
	a := assert.New(t)
	var custom meta.Interface = ""
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", custom},
		HttpInterface.Dependency())))
	a.Nil(managerPlugin.Add(pluginHTTP(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom}, SQLInterface.Dependency())))
	a.Nil(managerPlugin.Add(plugin3()))
	a.Nil(managerPlugin.Add(plugin2(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginMysql(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginCms()))

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
		index1     int
		index2     int
		index3     int
		indexHTTP  int
		indexMysql int
		indexCms   int
	)
	for index, value := range pluginsSorted {
		if value.ID == pluginOne {
			index1 = index
		}
		if value.ID == pluginTwo {
			index2 = index
		}
		if value.ID == pluginThree {
			index3 = index
		}
		if value.ID == "pluginCms" {
			indexCms = index
		}
		if value.ID == "pluginHTTP" {
			indexHTTP = index
		}
		if value.ID == "pluginMysql" {
			indexMysql = index
		}
	}
	a.Equal(true, index3 < indexMysql)
	a.Equal(true, index3 < indexHTTP)
	a.Equal(true, indexHTTP < indexCms)
	a.Equal(true, indexHTTP > indexMysql)
	a.Equal(true, indexMysql < indexCms)
	a.Equal(true, index3 < index2)
	a.Equal(true, index2 < index1)
	a.Equal(true, indexHTTP < index1)

	a.Equal(6, len(pluginsSorted))
}
*/
//as 11 test but have ring pluginCms->pluginCms
//12) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHTTP and pluginMysql -> plugin3,
/*func TestDependencyGraph_SortWithRing(t *testing.T) {
	a := assert.New(t)
	var custom meta.Interface = ""
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", custom},
		HttpInterface.Dependency())))
	a.Nil(managerPlugin.Add(pluginHTTP(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(plugin3()))
	a.Nil(managerPlugin.Add(plugin2()))
	a.Nil(managerPlugin.Add(pluginMysql(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Equal(errors.LoopVertexFound{Dependency: struct {
		ID         meta.PluginID
		Constraint string
		Interface  meta.Interface
	}{ID: "pluginCms", Constraint: ">=1.0.0, <2.0.0", Interface: ""}}, managerPlugin.Add(pluginCms(meta.Dependency{"pluginCms", ">=1.0.0, <2.0.0", custom})))

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
	_, err := managerPlugin.Sort()
	a.NoError(err)
	t.Log(err)

}
*/
//as 11 test but have ring pluginHTTP->plugin3, plugin3->pluginHTTP
//13) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHTTP and pluginMysql -> plugin3
/*func TestDependencyGraph_SortWithRing2(t *testing.T) {
	a := assert.New(t)
	var custom meta.Interface = ""
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", custom},
		HttpInterface.Dependency())))
	a.Nil(managerPlugin.Add(pluginHTTP(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{"pluginHTTP", ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(plugin2(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginMysql(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginCms()))

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
*/
//// ring through interface as 11 test but have ring pluginCms>interfaceHttp, plugin3->pluginCms
//14) plugin1 -> plugin2 -> plugin3 order for adding - 1 3 2, plugin1->Interface Http, pluginCms->interfaceHttp and interfaceMysql (plugins with such interfaces added),pluginHTTP and pluginMysql -> plugin3
/*func TestDependencyGraph_SortWithRing3(t *testing.T) {
	a := assert.New(t)
	var custom meta.Interface = ""
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1(meta.Dependency{pluginTwo, ">=1.0.0, <2.0.0", custom},
		HttpInterface.Dependency())))
	a.Nil(managerPlugin.Add(pluginHTTP(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{"pluginCms", ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(plugin2(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginMysql(meta.Dependency{pluginThree, ">=1.0.0, <2.0.0", custom})))
	a.Nil(managerPlugin.Add(pluginCms()))

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
*/
// ring through 1 plugin, between plugin1 and plugin3
//15) plugin1 -> plugin2 -> plugin3, plugin3->1 (all available) order for adding - 1 2 3
/*func TestDependencyGraph_SortWithRing4(t *testing.T) {
	a := assert.New(t)
	managerPlugin := dependency.NewManager()
	a.Nil(managerPlugin.Add(plugin1()))
	a.Nil(managerPlugin.Add(plugin2()))
	a.Nil(managerPlugin.Add(plugin3(meta.Dependency{pluginOne, ">=1.0.0, <2.0.0", meta.Interface("")})))

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
*/
