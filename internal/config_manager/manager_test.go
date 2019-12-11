package config_manager_test

import (
	"testing"

	"github.com/cheebo/go-config"
	commonCfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/internal/config_manager"
	"github.com/stretchr/testify/assert"
)

func TestManager_Register(t *testing.T) {
	a := assert.New(t)

	id := "nori/testing"

	pluginMeta := meta.Data{
		ID: meta.ID{
			ID:      meta.PluginID(id),
			Version: "1.0.0",
		},
	}

	mConfig := go_config.New()

	manager := config_manager.NewManager(mConfig)
	cm := manager.Register(pluginMeta)
	cm.String("http.addr", "HTTP server addr")
	cm.String("http.enabled", "Enable HTTP server")

	vars := manager.PluginVariables(pluginMeta.ID)

	a.Len(vars, 2)
	if len(vars) == 2 {
		a.Equal(commonCfg.Variable{
			Name:        "http.addr",
			Description: "HTTP server addr",
		}, vars[0])
		a.Equal(commonCfg.Variable{
			Name:        "http.enabled",
			Description: "Enable HTTP server",
		}, vars[1])
	}
}
