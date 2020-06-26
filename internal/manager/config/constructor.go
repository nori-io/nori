package config

import (
	go_config "github.com/cheebo/go-config"
	"github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/meta"
	"go.uber.org/fx"
)

type ManagerParams struct {
	fx.In
	//GoConfig go_config.Config
}

func New(params ManagerParams) (config.Manager, error) {
	m := new(manager)
	m.configs = make(map[meta.ID]*[]config.Variable)
	m.config = go_config.New()
	return m, nil
}
