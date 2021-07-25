package nori

import (
	"context"
	"fmt"
	"sync"

	"github.com/nori-io/common/v5/pkg/domain/event"
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	common_registry "github.com/nori-io/common/v5/pkg/domain/registry"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
	"github.com/nori-io/nori/pkg/nori/domain/errors"
	domain_registry "github.com/nori-io/nori/pkg/nori/domain/registry"
	"github.com/nori-io/nori/pkg/nori/event_emitter"
	"github.com/nori-io/nori/pkg/nori/helper/dependency_graph"
	"github.com/nori-io/nori/pkg/nori/registry/registry"
)

type Engine struct {
	eventEmitter    event.EventEmitter
	dependencyGraph *dependency_graph.DependencyGraphHelper
	registry        domain_registry.Registry
	configRegistry  common_registry.ConfigRegistry
	logger          logger.Logger
	mu              sync.Mutex
}

func New(configRegistry common_registry.ConfigRegistry, logger logger.Logger) (Nori, error) {
	r := registry.New()
	return &Engine{
		eventEmitter:    event_emitter.NewEventEmitter(),
		dependencyGraph: dependency_graph.New(r),
		registry:        r,
		configRegistry:  configRegistry,
		logger:          logger,
		mu:              sync.Mutex{},
	}, nil
}

func (n *Engine) Add(p *entity.Plugin) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	id := p.Meta().GetID()

	// check exists
	if exist := n.registry.GetByID(id); exist != nil {
		return errors.AlreadyExists{ID: id}
	}

	// add
	if err := n.registry.Add(p); err != nil {
		return err
	}
	if err := n.dependencyGraph.Add(id, p.Meta().GetDependencies()); err != nil {
		if err := n.registry.Remove(id); err != nil {
			panic("cannot rollback nori.add")
		}
		return err
	}

	return nil
}

func (n *Engine) Remove(id meta.ID) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// check state
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}
	if p.State() == enum.Running {
		return fmt.Errorf("cannot remove running plugin")
	}

	// remove
	if err := n.dependencyGraph.Remove(p.Meta().GetID()); err != nil {
		return err
	}
	return n.registry.Remove(p.Meta().GetID())
}

func (n *Engine) Init(ctx context.Context, id meta.ID) error {
	// lock
	n.mu.Lock()
	defer n.mu.Unlock()

	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	return p.Init(ctx, n.configRegistry, n.logger)
}

func (n *Engine) InitAll(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.dependencyGraph.HasUnresolved() {
		return errors.DependenciesNotFound{Dependencies: n.dependencyGraph.GetUnresolved()}
	}

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

	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	return p.Start(ctx, n.registry)
}

func (n *Engine) StartAll(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.dependencyGraph.HasUnresolved() {
		return errors.DependenciesNotFound{Dependencies: n.dependencyGraph.GetUnresolved()}
	}

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

	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	err := p.Stop(ctx, n.registry)
	if err != nil {
		n.getLogger(p.Meta()).Error(fmt.Sprintf("plugin %s has been installed with error %s", p.Meta().GetID().String(), err.Error()))
	} else {
		n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s has been installed", p.Meta().GetID().String()))
	}
	return err
}

func (n *Engine) StopAll(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

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

func (n *Engine) Install(ctx context.Context, plugin *entity.Plugin) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	id := plugin.Meta().GetID()

	if n.registry.GetByID(id) != nil {
		return errors.AlreadyExists{ID: id}
	}

	if err := plugin.Install(ctx, n.registry); err != nil {
		n.getLogger(plugin.Meta()).Error(fmt.Sprintf("plugin %s has been installed with error %s", plugin.Meta().GetID().String(), err.Error()))
		return err
	}

	if err := n.Add(plugin); err != nil {
		n.getLogger(plugin.Meta()).Info(fmt.Sprintf("rallback installation of plugin %s ", plugin.Meta().GetID().String()))
		return plugin.UnInstall(ctx, n.registry)
	}

	n.getLogger(plugin.Meta()).Info(fmt.Sprintf("plugin %s has been installed", plugin.Meta().GetID().String()))
	return nil
}

func (n *Engine) UnInstall(ctx context.Context, id meta.ID) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	if err := p.UnInstall(ctx, n.registry); err != nil {
		n.getLogger(p.Meta()).Error(fmt.Sprintf("plugin %s [%s] has been uninstalled with error %s", id.String(), p.Meta().GetInterface().String(), err.Error()))
	}

	n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s [%s] has been uninstalled", id.String(), p.Meta().GetInterface().String()))
	return nil
}

func (n *Engine) GetByFilter(filter Filter) []meta.ID {
	//return n.state.GetAllByState(filter.State)
	// todo
	return nil
}

func (n *Engine) GetPluginVariables(id meta.ID) []common_registry.Variable {
	return n.configRegistry.PluginVariables(id)
}

func (n *Engine) GetState(id meta.ID) (enum.State, error) {
	p := n.registry.GetByID(id)
	if p == nil {
		return 0, errors.NotFound{ID: id}
	}
	return p.State(), nil
}

func (n *Engine) init(ctx context.Context, id meta.ID) error {
	// get plugin
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	// init dependencies
	for _, did := range n.dependencyGraph.GetDependencies(id) {
		if err := n.init(ctx, did); err != nil {
			return err
		}
	}

	// init
	err := p.Init(ctx, n.configRegistry, n.logger)
	if err != nil {
		// todo: fix msg, plugin state has not changed to 'inited'
		n.getLogger(p.Meta()).Error(fmt.Sprintf("plugin %s has inited with error %s", id.String(), err.Error()))
		return err
	}

	n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s has inited", id.String()))

	// subscribe
	if p.IsNotifiable() {
		if err := p.Subscribe(n.eventEmitter); err == nil {
			n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s subscribed to events", id.String()))
		}
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
	if err := p.Start(ctx, n.registry); err != nil {
		// todo: fix msg, plugin state has not changed to 'started'
		n.getLogger(p.Meta()).Error(fmt.Sprintf("plugin %s has started with error %s", id.String(), err.Error()))
	}

	n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s has started", id.String()))
	return nil
}

func (n *Engine) stop(ctx context.Context, id meta.ID) error {
	// get plugin
	p := n.registry.GetByID(id)
	if p == nil {
		return errors.NotFound{ID: id}
	}

	// stop dependent
	for _, did := range n.dependencyGraph.GetDependent(id) {
		if err := n.stop(ctx, did); err != nil {
			return err
		}
	}

	// stop
	if err := p.Stop(ctx, n.registry); err != nil {
		// todo: fix msg, plugin state has not changed to 'stopped'
		n.getLogger(p.Meta()).Error(fmt.Sprintf("plugin %s has stopped with error %s", id.String(), err.Error()))
	}

	n.getLogger(p.Meta()).Info(fmt.Sprintf("plugin %s has stopped", id.String()))
	return nil
}

func (n *Engine) getLogger(m meta.Meta) logger.Logger {
	return n.logger.With(logger.Field{
		Key:   "plugin_id",
		Value: m.GetID().String(),
	}).With(logger.Field{
		Key:   "plugin_interface",
		Value: m.GetInterface().String(),
	})
}
