package nori

import (
	"context"
	"fmt"
	"sync"

	"github.com/nori-io/common/v5/pkg/domain/event"
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/common/v5/pkg/domain/registry"
	"github.com/nori-io/nori/pkg/errors"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
	dregistry "github.com/nori-io/nori/pkg/nori/domain/registry"
	"github.com/nori-io/nori/pkg/nori/event_emitter"
	"github.com/nori-io/nori/pkg/nori/helper/dependency_graph"
	"github.com/nori-io/nori/pkg/nori/helper/state"
	"github.com/nori-io/nori/pkg/nori/registry/config"
	rregistry "github.com/nori-io/nori/pkg/nori/registry/registry"
)

type Engine struct {
	eventEmitter    event.EventEmitter
	dependencyGraph *dependency_graph.DependencyGraphHelper
	registry        dregistry.Registry
	state           *state.StateHelper
	config          registry.ConfigRegistry
	logger          logger.Logger
	mu              sync.Mutex
}

func New(logger logger.Logger) (Nori, error) {
	r := rregistry.New()
	cm, err := config.New(config.Params{File: "/etc/nori/config.yml"})
	if err != nil {
		return nil, err
	}
	return &Engine{
		eventEmitter:    event_emitter.NewEventEmitter(),
		dependencyGraph: dependency_graph.New(r),
		registry:        r,
		state:           state.New(),
		config:          cm,
		logger:          logger,
		mu:              sync.Mutex{},
	}, nil
}

func (n *Engine) Add(p plugin.Plugin) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// check exists
	if exist := n.registry.GetByID(p.Meta().GetID()); exist != nil {
		return errors.AlreadyExists{ID: p.Meta().GetID()}
	}

	// add
	if err := n.registry.Add(p); err != nil {
		return err
	}
	if err := n.dependencyGraph.Add(p); err != nil {
		if err := n.registry.Remove(p); err != nil {
			panic("cannot rollback nori.add")
		}
		return err
	}

	n.state.SetState(p.Meta().GetID(), enum.None)
	return nil
}

func (n *Engine) Remove(p plugin.Plugin) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// check state
	state, err := n.state.GetState(p.Meta().GetID())
	if err != nil {
		return err
	}
	if state == enum.Running {
		return fmt.Errorf("cannot remove running plugin")
	}

	// remove
	if err := n.dependencyGraph.Remove(p); err != nil {
		return err
	}
	return n.registry.Remove(p)
}

func (n *Engine) Init(ctx context.Context, id meta.ID) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.init(ctx, id)
}

func (n *Engine) InitAll(ctx context.Context) error {
	// check unresolved
	if n.hasUnresolvedDependencies() {
		return errors.DependenciesNotFound{Dependencies: n.dependencyGraph.GetUnresolved()}
	}

	// init all
	list, err := n.dependencyGraph.GetSorted()
	if err != nil {
		return err
	}
	for _, id := range list {
		if err := n.init(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (n *Engine) Start(ctx context.Context, id meta.ID) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.start(ctx, id)
}

func (n *Engine) StartAll(ctx context.Context) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// check unresolved
	if n.hasUnresolvedDependencies() {
		return errors.DependenciesNotFound{Dependencies: n.dependencyGraph.GetUnresolved()}
	}

	// start all
	list, err := n.dependencyGraph.GetSorted()
	if err != nil {
		return err
	}

	for _, id := range list {
		if err := n.start(ctx, id); err != nil {
			return err
		}
	}

	n.eventEmitter.Emit(enum.Event_NoriPluginsStarted, nil)

	return nil
}

func (n *Engine) Stop(ctx context.Context, id meta.ID) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.stop(ctx, id)
}

func (n *Engine) StopAll(ctx context.Context) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// stop all
	ids, err := n.dependencyGraph.GetSortedReversed()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := n.stop(ctx, id); err != nil {
			return err
		}
	}

	n.eventEmitter.Emit(enum.Event_NoriPluginsStopped, nil)

	return nil
}

func (n *Engine) Install(ctx context.Context, p plugin.Plugin) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			ID:   p.Meta().GetID(),
			Path: "",
		}
	}

	if n.registry.GetByID(p.Meta().GetID()) != nil {
		return errors.AlreadyExists{ID: p.Meta().GetID()}
	}

	err := n.install(ctx, installable)

	if err != nil {
		n.getLogger(p).Error(fmt.Sprintf("plugin %s has been installed with error %s", p.Meta().GetID().String(), err.Error()))
	} else {
		n.getLogger(p).Info(fmt.Sprintf("plugin %s has been installed", p.Meta().GetID().String()))
	}

	return err
}

func (n *Engine) UnInstall(ctx context.Context, id meta.ID) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	state, err := n.state.GetState(id)
	if err != nil {
		return err
	}
	if state == enum.Running {
		return fmt.Errorf("plugin %s is running, cannot uninstall it", id.String())
	}

	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	installable, ok := (p).(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			ID:   id,
			Path: "",
		}
	}

	err = n.uninstall(ctx, installable)

	if err != nil {
		n.getLogger(p).Error(fmt.Sprintf("plugin %s [%s] has been uninstalled with error %s", id.String(), p.Meta().GetInterface().String(), err.Error()))
	} else {
		n.getLogger(p).Info(fmt.Sprintf("plugin %s [%s] has been uninstalled", id.String(), p.Meta().GetInterface().String()))
	}

	return err
}

func (n *Engine) GetByFilter(filter Filter) []meta.ID {
	return n.state.GetAllByState(filter.State)
}

func (n *Engine) GetPluginVariables(id meta.ID) []registry.Variable {
	return n.config.PluginVariables(id)
}

func (n *Engine) GetState(id meta.ID) (enum.State, error) {
	return n.state.GetState(id)
}

func (n *Engine) init(ctx context.Context, id meta.ID) error {
	// get plugin
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	// check state
	if err := n.preInitPlugin(id); err != nil {
		switch err.(type) {
		case ErrorPluginAlreadyInited:
			return nil
		default:
			return err
		}
	}

	// init dependencies
	for _, did := range n.dependencyGraph.GetDependencies(id) {
		if err := n.init(ctx, did); err != nil {
			return err
		}
	}

	// init
	err := n.initPlugin(ctx, p)
	l := n.getLogger(p)
	if err != nil {
		// todo: fix msg, plugin state has not changed to 'inited'
		l.Error(fmt.Sprintf("plugin %s has inited with error %s", id.String(), err.Error()))
		return err
	} else {
		l.Info(fmt.Sprintf("plugin %s has inited", id.String()))
	}

	notifiable, ok := p.(plugin.Notifiable)
	if ok {
		notifiable.Subscribe(n.eventEmitter)
		l.Info(fmt.Sprintf("plugin %s subscribed to events", id.String()))
	}

	return nil
}

func (n *Engine) start(ctx context.Context, id meta.ID) error {
	// get plugin
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	// start dependencies
	for _, did := range n.dependencyGraph.GetDependencies(id) {
		if err := n.start(ctx, did); err != nil {
			return err
		}
	}

	// start
	err := n.startPlugin(ctx, p)

	if err != nil {
		// todo: fix msg, plugin state has not changed to 'started'
		n.getLogger(p).Error(fmt.Sprintf("plugin %s has started with error %s", id.String(), err.Error()))
	} else {
		n.getLogger(p).Info(fmt.Sprintf("plugin %s has started", id.String()))
	}

	return err
}

func (n *Engine) stop(ctx context.Context, id meta.ID) error {
	// get plugin
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	// check state
	if err := n.preStopPlugin(id); err != nil {
		return err
	}

	// stop dependent
	for _, did := range n.dependencyGraph.GetDependent(id) {
		if err := n.stop(ctx, did); err != nil {
			return err
		}
		n.state.SetState(id, enum.Inited)
	}

	// stop
	err := n.stopPlugin(ctx, p)
	if err != nil {
		// todo: fix msg, plugin state has not changed to 'stopped'
		n.getLogger(p).Error(fmt.Sprintf("plugin %s has stopped with error %s", id.String(), err.Error()))
	} else {
		n.getLogger(p).Info(fmt.Sprintf("plugin %s has stopped", id.String()))
	}
	return err
}

func (n *Engine) getLogger(p plugin.Plugin) logger.Logger {
	return n.logger.With(logger.Field{
		Key:   "plugin_id",
		Value: p.Meta().GetID().String(),
	}).With(logger.Field{
		Key:   "plugin_interface",
		Value: p.Meta().GetInterface().String(),
	})
}
