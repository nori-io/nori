package plugins_test

import (
	"testing"

	"github.com/nori-io/nori/core/log"

	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/core/mocks"

	ncmocks "github.com/nori-io/nori-common/mocks"
	"github.com/nori-io/nori/core/plugins"
	"github.com/nori-io/nori/version"
	"github.com/stretchr/testify/assert"
)

func TestManager_AddFile(t *testing.T) {
	a := assert.New(t)

	// storage
	storage := &mocks.Storage{}
	storagePlugins := &mocks.Plugins{}
	storage.On("Plugins").Return(storagePlugins)
	// config manager
	cfgManager := &ncmocks.Manager{}
	// logger
	logger := log.New()
	// version
	ver := version.NoriVersion(logger)
	// plugin extractor
	mockedPlugin := &ncmocks.Plugin{}
	mockedPlugin.On("Meta").Return(meta.Data{
		ID: meta.ID{
			ID:      "nori/test",
			Version: "1.0.2",
		},
		Core: meta.Core{
			VersionConstraint: version.CurrentVersion,
		},
		Interface: meta.Interface(""),
	})
	pluginExtractor := &mocks.PluginExtractor{}
	pluginExtractor.On("Get", "nori/plugin.so").Return(mockedPlugin, nil)

	m := plugins.NewManager(storage, cfgManager, ver, pluginExtractor, logger)

	pl, err := m.AddFile("nori/plugin.so")
	a.Nil(err)
	a.NotNil(pl)
}
