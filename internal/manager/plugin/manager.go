package plugin

import (
	"context"

	"github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/enum/status"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/pkg/errors"
)

type Manager struct {
	pluginRepository repository.PluginRepository
	registryService  plugin.Registry
	cm               config.Manager
	logger           logger.Logger
	bucket           storage.Bucket
}

func (m *Manager) Register(plugin *entity.Plugin) error {
	// todo: compare core version and plugin core.version
	err := m.pluginRepository.Register(plugin)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RegisterAll(plugins []*entity.Plugin) error {
	for _, plugin := range plugins {
		// todo: compare core version and plugin core.version
		err := m.pluginRepository.Register(plugin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) UnRegister(p *entity.Plugin) error {
	return m.pluginRepository.UnRegister(p)
}

func (m *Manager) GetAll() []*entity.Plugin {
	return m.pluginRepository.FindAll()
}

func (m *Manager) GetByID(id meta.ID) (*entity.Plugin, error) {
	plugin := m.pluginRepository.FindByID(id)
	if plugin == nil {
		return nil, errors.NotFound{ID: id}
	}
	return plugin, nil
}

func (m *Manager) GetHooks() []*entity.Plugin {
	return m.pluginRepository.FindHooks()
}

func (m *Manager) GetInstallable() []*entity.Plugin {
	return m.pluginRepository.FindInstallable()
}

func (m *Manager) GetRunning() []*entity.Plugin {
	return nil
}

func (m *Manager) Start(ctx context.Context, id meta.ID) error {
	plugin := m.pluginRepository.FindByID(id)
	if plugin == nil {
		return errors.NotFound{ID: id}
	}

	if plugin.Status() == status.Started {
		// already started
		return nil
	}

	err := plugin.Init(ctx, m.cm.Register(plugin.Meta().Id()), m.logger)
	if err != nil {
		return err
	}

	// start dependencies first
	for _, dep := range plugin.Meta().GetDependencies() {
		p, err := m.pluginRepository.Resolve(dep)
		if err != nil {
			return err
		}
		err = m.Start(ctx, p.Meta().Id())
		if err != nil {
			return err
		}
	}

	if plugin.Status() == status.Nil {
		// todo: fill params
		err := plugin.Init(ctx, nil, nil)
		if err != nil {
			return err
		}
	}

	return plugin.Start(ctx, m.registryService)
}

func (m *Manager) StartAll(ctx context.Context) error {
	ids, err := m.pluginRepository.Tree()
	if err != nil {
		return err
	}

	for _, p := range ids {
		err := m.Start(ctx, p.Meta().Id())
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Stop(ctx context.Context, id meta.ID) error {
	plugin := m.pluginRepository.FindByID(id)
	if plugin == nil {
		return errors.NotFound{ID: id}
	}

	if plugin.Status() == status.Nil {
		// cannot stop uninitialised and not started plugin
		// todo: add notification on potentially incorrect behavior on top level
		return nil
	}

	if plugin.Status() == status.Stopped {
		// already stopped
		return nil
	}

	// stop dependent plugins
	for _, did := range m.pluginRepository.FindDependent(id) {
		err := m.Stop(ctx, did.Meta().Id())
		if err != nil {
			return err
		}
	}

	return plugin.Stop(ctx, m.registryService)
}

func (m *Manager) StopAll(ctx context.Context) error {
	ids, err := m.pluginRepository.Tree()
	if err != nil {
		return err
	}

	// stop plugins in revers order
	for i := len(ids) - 1; i >= 0; i-- {
		err := m.Start(ctx, ids[i].Meta().Id())
		if err != nil {
			return err
		}
	}

	return nil
}
