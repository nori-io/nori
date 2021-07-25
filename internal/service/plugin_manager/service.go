package plugin_manager

import (
	"context"
	"fmt"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	errors2 "github.com/nori-io/common/v5/pkg/errors"
	nori_entity "github.com/nori-io/nori/pkg/nori/domain/entity"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/pkg/nori"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
	errors3 "github.com/nori-io/nori/pkg/nori/domain/errors"
)

func (s *PluginManager) Enable(ctx context.Context, id meta.ID) error {
	p, err := s.PluginService.Get(id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("plugin %s not found", id.String())
	}

	po, err := s.PluginOptionService.Get(id)
	if err != nil {
		if _, ok := err.(errors2.EntityNotFound); !ok {
			return err
		}
	}

	data := service.PluginOptionUpsertData{
		ID:          id,
		Enabled:     true,
		Installed:   false,
		Installable: false,
	}
	if !po.IsEmpty() {
		data.Installed = po.Installed
		data.Installable = po.Installable
	}

	_, err = s.PluginOptionService.Upsert(data)

	return s.Nori.Add(p)
}

func (s *PluginManager) Disable(ctx context.Context, id meta.ID) error {
	p, err := s.PluginService.Get(id)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("plugin %s not found", id.String())
	}

	if err := s.Nori.Stop(ctx, id); err != nil {
		return err
	}

	po, err := s.PluginOptionService.Get(id)
	if err != nil {
		if _, ok := err.(errors2.EntityNotFound); !ok {
			return err
		}
	}
	data := service.PluginOptionUpsertData{
		ID:          id,
		Enabled:     false,
		Installed:   false,
		Installable: false,
	}
	if !po.IsEmpty() {
		data.Installed = po.Installed
		data.Installable = po.Installable
	}

	_, err = s.PluginOptionService.Upsert(data)
	if err != nil {
		return err
	}

	return s.Nori.Remove(p.Meta().GetID())
}

func (s *PluginManager) Install(ctx context.Context, id meta.ID) error {
	p, err := s.PluginService.Get(id)
	if err != nil {
		return err
	}

	if !p.IsInstallable() {
		return errors3.NonInstallablePlugin{
			ID:   id,
			Path: p.File(),
		}
	}

	po, err := s.PluginOptionService.Get(id)
	if err != nil {
		if _, ok := err.(errors2.EntityNotFound); !ok {
			return err
		}
	}
	if po.Installed {
		return nil
	}

	err = s.Nori.Install(ctx, p)
	if err != nil {
		return err
	}

	_, err = s.PluginOptionService.Upsert(service.PluginOptionUpsertData{
		ID:          id,
		Enabled:     false,
		Installed:   true,
		Installable: true,
	})

	return nil
}

func (s *PluginManager) UnInstall(ctx context.Context, id meta.ID) error {
	p, err := s.PluginService.Get(id)
	if err != nil {
		return err
	}

	if !p.IsInstallable() {
		return errors3.NonInstallablePlugin{
			ID:   id,
			Path: p.File(),
		}
	}

	po, err := s.PluginOptionService.Get(id)
	if err != nil {
		if _, ok := err.(errors2.EntityNotFound); !ok {
			return err
		}
	}
	if !po.Installed {
		return nil
	}

	err = s.Nori.UnInstall(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.PluginOptionService.Upsert(service.PluginOptionUpsertData{
		ID:          id,
		Enabled:     false,
		Installed:   false,
		Installable: true,
	})

	return nil
}

func (s *PluginManager) Start(ctx context.Context, id meta.ID) error {
	// todo: check inited
	err := s.Nori.Start(ctx, id)
	if err != nil {
		s.Env.Logger.Error(fmt.Sprintf("Plugin %s started with error %s", id.String(), err.Error()))
	} else {
		s.Env.Logger.Info(fmt.Sprintf("Plugin %s started", id.String()))
	}
	return err
}

func (s *PluginManager) Stop(ctx context.Context, id meta.ID) error {
	err := s.Nori.Stop(ctx, id)
	if err != nil {
		s.Env.Logger.Error(fmt.Sprintf("Plugin %s stopped with error %s", id.String(), err.Error()))
	} else {
		s.Env.Logger.Info(fmt.Sprintf("Plugin %s stopped", id.String()))
	}
	return err
}

func (s *PluginManager) StartAll(ctx context.Context) error {
	return s.Nori.StartAll(ctx)
}

func (s *PluginManager) StopAll(ctx context.Context) error {
	return s.Nori.StopAll(ctx)
}

func (s *PluginManager) GetByFilter(filter service.GetByFilterData) ([]*nori_entity.Plugin, error) {
	ids := s.Nori.GetByFilter(nori.Filter{
		State: enum.New(filter.State.Value()),
	})
	plugins := s.PluginService.GetByIDs(ids)
	return plugins, nil
}
