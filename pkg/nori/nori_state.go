package nori

import (
	"context"
	"fmt"

	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
)

func (n *Engine) preInitPlugin(id meta.ID) error {
	// check state
	state, err := n.state.GetState(id)
	if err != nil {
		return err
	}
	if state == enum.Inited || state == enum.Running {
		return ErrorPluginAlreadyInited{ID: id}
	}
	if state != enum.None {
		return fmt.Errorf("cannot initialize %s plugin that has %s state", id.String(), state.String())
	}
	return nil
}

func (n *Engine) initPlugin(ctx context.Context, p plugin.Plugin) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.Init(ctx, n.config.Register(p.Meta().GetID()), n.logger.With(logger.Field{
			Key:   "plugin_id",
			Value: p.Meta().GetID().String(),
		}, logger.Field{
			Key:   "plugin_interface",
			Value: p.Meta().GetInterface().String(),
		}))
		if err == nil {
			n.state.SetState(p.Meta().GetID(), enum.Inited)
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

func (n *Engine) preStartPlugin(id meta.ID) error {
	// check state
	state, err := n.state.GetState(id)
	if err != nil {
		return err
	}
	if state == enum.Running {
		return nil
	}
	if state != enum.Inited {
		return fmt.Errorf("plugin must have %s state but have %s ", enum.Inited.String(), state)
	}
	return nil
}

func (n *Engine) startPlugin(ctx context.Context, p plugin.Plugin) error {
	var (
		err       error
		recovered interface{}
	)

	id := p.Meta().GetID()

	// check plugin state
	// todo: return error?
	state, err := n.state.GetState(id)
	if err != nil {
		return err
	}
	if state == enum.Running {
		return nil
	}
	if state == enum.None {
		if err := n.init(ctx, id); err != nil {
			return err
		}
	}

	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.Start(ctx, n.registry)
		if err == nil {
			n.state.SetState(id, enum.Running)
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

func (n *Engine) preStopPlugin(id meta.ID) error {
	// check state
	state, err := n.state.GetState(id)
	if err != nil {
		return err
	}
	if state == enum.Inited {
		return nil
	}
	if state != enum.Running {
		return nil
	}
	return nil
}

func (n *Engine) stopPlugin(ctx context.Context, p plugin.Plugin) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.Stop(ctx, n.registry)
		if err == nil {
			n.state.SetState(p.Meta().GetID(), enum.Inited)
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

func (n *Engine) install(ctx context.Context, p plugin.Installable) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.Install(ctx, n.registry)
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

func (n *Engine) uninstall(ctx context.Context, p plugin.Installable) error {
	var err error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		err = p.UnInstall(ctx, n.registry)
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

func (n *Engine) subscribe(p plugin.Notifiable) error {
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		p.Subscribe(n.eventEmitter)
	}()
	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}
	return nil
}

func (n *Engine) hasUnresolvedDependencies() bool {
	return n.dependencyGraph.HasUnresolved()
}
