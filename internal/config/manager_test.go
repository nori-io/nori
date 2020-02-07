package config_test

import (
	"testing"

	"github.com/cheebo/go-config"
	commonCfg "github.com/nori-io/nori-common/v2/config"
	typesCfg "github.com/nori-io/nori-common/v2/config/types"
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/internal/config"
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

	manager := config.NewManager(mConfig)
	cm := manager.Register(pluginMeta)
	cm.Int("http.addr", "HTTP server addr")
	cm.String("http.enabled", "Enable HTTP server")

	vars := manager.PluginVariables(pluginMeta.ID)

	a.Len(vars, 2)
	if len(vars) == 2 {
		a.Equal(commonCfg.Variable{
			Name:        "http.addr",
			Description: "HTTP server addr",
			Type:        typesCfg.Int,
		}, vars[0])
		a.Equal(commonCfg.Variable{
			Name:        "http.enabled",
			Description: "Enable HTTP server",
			Type:        typesCfg.String,
		}, vars[1])
	}
}

func TestManager_PluginVariables(t *testing.T) {
	a := assert.New(t)
	mConfig := go_config.New()

	manager := config.NewManager(mConfig)
	vars := manager.PluginVariables(meta.ID{
		ID:      "nori",
		Version: "1.0.0",
	})

	a.Len(vars, 0)
}
