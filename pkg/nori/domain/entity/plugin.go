package entity

import (
	"context"
	"fmt"

	"github.com/nori-io/common/v5/pkg/domain/event"
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/common/v5/pkg/domain/registry"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
	"github.com/nori-io/nori/pkg/nori/domain/errors"
)

type Plugin struct {
	file   File
	plugin plugin.Plugin
	state  enum.State
	notifiable plugin.Notifiable
	installable plugin.Installable
}

func New(file File) (*Plugin, error) {
	var (
		err error
		recovered interface{}
		p plugin.Plugin
		notifiable plugin.Notifiable
		installable plugin.Installable
	)

	func() {
		defer func() {
			recovered = recover()
		}()
		p = file.Fn()
	}()
	if err != nil {
		return nil, err
	}
	if recovered != nil {
		// todo: custom error
		return nil, fmt.Errorf("%v", recovered)
	}

	installable, _ = p.(plugin.Installable)
	notifiable, _ = p.(plugin.Notifiable)

	return &Plugin{
		file:   file,
		plugin: p,
		state:  enum.None,
		notifiable: notifiable,
		installable: installable,
	}, nil
}

func (p *Plugin) File() string {
	return p.file.Path
}

func (p *Plugin) Init(ctx context.Context, cr registry.ConfigRegistry, log logger.Logger) error {
	// already initialized
	if p.state == enum.Inited {
		return nil
	}
	if p.state != enum.None {
		return errors.PluginIncorrectState{
			Expected: enum.None,
			Current:  p.state,
		}
	}

	var (
		err error
		recovered interface{}
	)
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Init(ctx, cr.Register(p.plugin.Meta().GetID()), log.With(logger.Field{
			Key:   "plugin_id",
			Value: p.plugin.Meta().GetID().String(),
		}, logger.Field{
			Key:   "plugin_interface",
			Value: p.plugin.Meta().GetInterface().String(),
		}))
		if err == nil {
			p.state = enum.Inited
		}
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) Start(ctx context.Context, r registry.Registry) error {
	// already initialized
	if p.state == enum.Running {
		return nil
	}
	if p.state != enum.Inited {
		return errors.PluginIncorrectState{
			Expected: enum.Inited,
			Current:  p.state,
		}
	}

	var (
		err       error
		recovered interface{}
	)

	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Start(ctx, r)
		if err == nil {
			p.state = enum.Running
		}
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) Stop(ctx context.Context, r registry.Registry) error {
	if p.state == enum.Inited {
		return nil
	}
	if p.state != enum.Running {
		return errors.PluginIncorrectState{
			Expected: enum.Running,
			Current:  p.state,
		}
	}

	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.plugin.Stop(ctx, r)
		if err == nil {
			p.state = enum.Inited
		}
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) Install(ctx context.Context, r registry.Registry) error {
	if p.installable == nil {
		return errors.NonInstallablePlugin{
			ID:   p.Meta().GetID(),
			Path: p.file.Path,
		}
	}

	if p.state != enum.None {
		return errors.PluginIncorrectState{
			Expected: enum.None,
			Current:  p.state,
		}
	}

	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.installable.Install(ctx, r)
		if err == nil {
			p.state = enum.Inited
		}
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) UnInstall(ctx context.Context, r registry.Registry) error {
	if p.installable == nil {
		return errors.NonInstallablePlugin{
			ID:   p.Meta().GetID(),
			Path: p.file.Path,
		}
	}

	if p.state != enum.None && p.state != enum.Inited {
		return errors.PluginIncorrectState{
			Expected: enum.None,
			Current:  p.state,
		}
	}

	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.installable.UnInstall(ctx, r)
		if err == nil {
			p.state = enum.Inited
		}
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) Subscribe(e event.EventEmitter) error {
	if p.installable == nil {
		return errors.NonInstallablePlugin{
			ID:   p.Meta().GetID(),
			Path: p.file.Path,
		}
	}

	if p.state != enum.None && p.state != enum.Inited {
		return errors.PluginIncorrectState{
			Expected: enum.None,
			Current:  p.state,
		}
	}

	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		p.notifiable.Subscribe(e)
	}()
	if err != nil {
		return err
	}
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (p *Plugin) Meta() meta.Meta {
	return p.plugin.Meta()
}

func (p *Plugin) Instance() interface{} {
	return p.plugin.Instance()
}

func (p *Plugin) State() enum.State {
	return p.state
}

func (p *Plugin) IsInstallable() bool {
	return p.installable == nil
}

func (p *Plugin) IsNotifiable() bool {
	return p.notifiable == nil
}