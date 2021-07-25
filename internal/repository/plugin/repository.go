package plugin

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	domain_errors "github.com/nori-io/nori/pkg/nori/domain/errors"
)

func (r *PluginRepository) Add(plugin *entity.Plugin) error {
	if _, ok := r.files[plugin.File()]; ok {
		return nil
	}

	// todo
	r.plugins[plugin.Meta().GetID().String()] = plugin
	r.files[plugin.File()] = plugin

	return nil
}

func (r *PluginRepository) Remove(file string) error {
	delete(r.files, file)
	for id, plugin := range r.plugins {
		if plugin.File() == file {
			delete(r.plugins, id)
			break
		}
	}
	return nil
}

func (r *PluginRepository) Find(id meta.ID) (*entity.Plugin, error) {
	plugin, ok := r.plugins[id.String()]
	if !ok {
		return nil, domain_errors.NotFound{ID: id}
	}
	return plugin, nil
}

func (r *PluginRepository) FindAll() []*entity.Plugin {
	items := []*entity.Plugin{}
	for _, plugin := range r.plugins {
		items = append(items, plugin)
	}
	return items
}

func (r *PluginRepository) FindByIDs(ids []meta.ID) []*entity.Plugin {
	items := []*entity.Plugin{}
	for _, plugin := range r.plugins {
		for _, id := range ids {
			if id == plugin.Meta().GetID() {
				items = append(items, plugin)
			}
		}
	}
	return items
}
