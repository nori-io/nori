package plugins_test

import (
	"testing"

	"github.com/nori-io/nori/core/errors"
	"github.com/nori-io/nori/core/plugins"
	"github.com/stretchr/testify/assert"
)

func TestPluginExtractor_Get(t *testing.T) {
	a := assert.New(t)

	pe := plugins.NewPluginExtractor()
	plugin, err := pe.Get("testdata/plugin.so")

	a.Nil(err)
	a.NotNil(plugin)
}

func TestPluginExtractor_Get_FileDoesNotExist(t *testing.T) {
	a := assert.New(t)

	pe := plugins.NewPluginExtractor()
	plugin, err := pe.Get("testdata/nofile.so")

	a.Error(err)
	a.IsType(errors.FileDoesNotExist{}, err)
	a.Nil(plugin)
}

func TestPluginExtractor_Get_FileOpenError(t *testing.T) {
	a := assert.New(t)

	pe := plugins.NewPluginExtractor()
	plugin, err := pe.Get("testdata/empty.so")

	a.Error(err)
	a.IsType(errors.FileOpenError{}, err)
	a.Nil(plugin)
}

func TestPluginExtractor_GetLookupError(t *testing.T) {
	a := assert.New(t)

	pe := plugins.NewPluginExtractor()
	plugin, err := pe.Get("testdata/novariable.so")
	a.Error(err)
	a.Nil(plugin)
}

func TestPluginExtractor_Get_NoPluginInterfaceError(t *testing.T) {
	a := assert.New(t)

	pe := plugins.NewPluginExtractor()
	plugin, err := pe.Get("testdata/interface.so")

	a.Error(err)
	a.IsType(errors.NoPluginInterfaceError{}, err)
	a.Nil(plugin)
}
