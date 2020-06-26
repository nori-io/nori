package entity

import (
	"context"
	"fmt"

	"github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori/internal/domain/enum/status"
)

type Plugin struct {
	file          File
	plugin        plugin.Plugin
	status        status.Status
	isInstallable bool
	isHook        bool
}

func NewPlugin(f File, p plugin.Plugin) *Plugin {
	var (
		installable bool
		hook        bool
	)
	if _, ok := p.(plugin.Installable); ok {
		installable = true
	}
	if _, ok := p.(logger.Hook); ok {
		hook = true
	}
	return &Plugin{
		file:          f,
		plugin:        p,
		status:        0,
		isInstallable: installable,
		isHook:        hook,
	}
}

func (p *Plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Init(ctx, config, log)
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		return fmt.Errorf("%v", recovered)
	}

	p.status = status.Initialised
	return nil
}

func (p *Plugin) Start(ctx context.Context, registry plugin.Registry) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Start(ctx, registry)
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		return fmt.Errorf("%v", recovered)
	}

	p.status = status.Started
	return nil
}

func (p *Plugin) Stop(ctx context.Context, registry plugin.Registry) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Stop(ctx, registry)
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		return fmt.Errorf("%v", recovered)
	}

	p.status = status.Stopped
	return nil
}

func (p *Plugin) Install(ctx context.Context, registry plugin.Registry) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Stop(ctx, registry)
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		return fmt.Errorf("%v", recovered)
	}

	p.status = status.Nil
	return nil
}

func (p *Plugin) Meta() meta.Meta {
	return p.plugin.Meta()
}

func (p *Plugin) Instance() interface{} {
	return p.plugin.Instance()
}

func (p *Plugin) IsInstallable() bool {
	return p.isInstallable
}

func (p *Plugin) IsHook() bool {
	return p.isHook
}

func (p *Plugin) Status() status.Status {
	return p.status
}
