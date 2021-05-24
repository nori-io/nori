package entity

import (
	"time"

	"github.com/nori-io/common/v5/pkg/domain/meta"
)

type PluginOption struct {
	ID          meta.ID
	Enabled     bool
	Installed   bool
	Installable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p PluginOption) IsEmpty() bool {
	return p.ID == nil
}
