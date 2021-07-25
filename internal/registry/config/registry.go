package config

import (
	"strings"

	go_config "github.com/cheebo/go-config"
	"github.com/nori-io/common/v5/pkg/domain/config"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/registry"
)

type ConfigRegistry struct {
	config          go_config.Config
	pluginVariables map[meta.ID]*[]registry.Variable
}

func (m *ConfigRegistry) Register(id meta.ID) config.Config {
	vars := make([]registry.Variable, 0)
	m.pluginVariables[id] = &vars
	return &Config{
		variables: &vars,
		config:    m.config.Sub("plugins." + strings.ReplaceAll(id.String(), ".", "_")),
	}
}

func (m *ConfigRegistry) PluginVariables(id meta.ID) []registry.Variable {
	vars, ok := m.pluginVariables[id]
	if !ok {
		return []registry.Variable{}
	}
	return *vars
}
