package plugin_option

import (
	"encoding/json"
	"time"

	meta2 "github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type model struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Enabled     bool      `json:"enabled"`
	Installed   bool      `json:"installed"`
	Installable bool      `json:"installable"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newModel(op entity.PluginOption) model {
	return model{
		ID:          string(op.ID.GetID()),
		Version:     op.ID.GetVersion(),
		Enabled:     op.Enabled,
		Installed:   op.Installed,
		Installable: op.Installable,
		CreatedAt:   op.CreatedAt,
		UpdatedAt:   op.UpdatedAt,
	}
}

func (m model) Convert() entity.PluginOption {
	id := meta.ID{
		ID:      meta2.PluginID(m.ID),
		Version: m.Version,
	}
	return entity.PluginOption{
		ID:          id,
		Enabled:     m.Enabled,
		Installed:   m.Installed,
		Installable: m.Installable,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func encode(m model) ([]byte, error) {
	return json.Marshal(m)
}

func decode(data []byte) (model, error) {
	m := model{}
	if err := json.Unmarshal(data, &m); err != nil {
		return model{}, err
	}
	return m, nil
}
