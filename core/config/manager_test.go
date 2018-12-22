package config_test

import (
	"testing"

	"github.com/secure2work/nori/core/config"

	"github.com/secure2work/nori/core/plugins/meta"

	"github.com/secure2work/nori/core/config/mock"
	"github.com/stretchr/testify/assert"

	. "github.com/stretchr/testify/mock"
)

func TestManager_Register(t *testing.T) {
	a := assert.New(t)

	id := "nori/testing"

	pluginMeta := meta.Data{
		ID: meta.ID{
			ID:      meta.PluginID(id),
			Version: "1.0",
		},
	}

	mConfig := mocks.Config{}
	mConfig.On("String", AnythingOfType("string")).
		Return(func(s string) string {
			return s
		})

	manager := config.NewManager(&mConfig)
	cm := manager.Register(pluginMeta)
	cm.String("http.addr", "HTTP server addr")
	cm.String("http.enabled", "Enable HTTP server")

	vars := manager.PluginVariables(pluginMeta.ID)

	a.Len(vars, 2)
	if len(vars) == 2 {
		a.Equal(config.Variable{
			Name:        "http.addr",
			Description: "HTTP server addr",
		}, vars[0])
		a.Equal(config.Variable{
			Name:        "http.enabled",
			Description: "Enable HTTP server",
		}, vars[1])
	}
}
