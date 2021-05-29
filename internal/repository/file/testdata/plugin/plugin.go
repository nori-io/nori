package main

import (
	"context"

	"github.com/nori-io/common/v5/pkg/domain/config"
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	p "github.com/nori-io/common/v5/pkg/domain/plugin"
	pl "github.com/nori-io/common/v5/pkg/domain/registry"
	m "github.com/nori-io/common/v5/pkg/meta"
)

func New(log logger.FieldLogger) p.Plugin {
	return plugin{
		log: log,
	}
}

type plugin struct {
	log logger.FieldLogger
}

func (p plugin) Meta() meta.Meta {
	return m.Meta{}
}

func (p plugin) Instance() interface{} {
	return nil
}

func (p plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	return nil
}

func (p plugin) Start(ctx context.Context, registry pl.Registry) error {
	return nil
}

func (p plugin) Stop(ctx context.Context, registry pl.Registry) error {
	return nil
}
