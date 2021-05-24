package entity

import (
	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
)

type Plugin struct {
	File   string
	Fn     func() plugin.Plugin
	Plugin plugin.Plugin
	State  enum.State
}
