package plugin

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/pkg/errors"
)

func (r *PluginRepository) Create(file *entity.File) (*entity.Plugin, error) {
	if file.Fn == nil {
		return nil, errors.NoPluginInterfaceError{Path: file.Path}
	}

	if p, ok := r.files[file.Path]; ok {
		return p, nil
	}

	plugin := &entity.Plugin{
		File:   file.Path,
		Fn:     file.Fn,
		Plugin: file.Fn(),
	}

	// todo
	r.plugins[plugin.Plugin.Meta().GetID().String()] = plugin
	r.files[file.Path] = plugin

	return plugin, nil
}

func (r *PluginRepository) Delete(file *entity.File) error {
	delete(r.files, file.Path)
	for id, plugin := range r.plugins {
		if plugin.File == file.Path {
			delete(r.plugins, id)
			break
		}
	}
	return nil
}

func (r *PluginRepository) Find(id meta.ID) (*entity.Plugin, error) {
	plugin, ok := r.plugins[id.String()]
	if !ok {
		return nil, errors.NotFound{ID: id}
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
			if id == plugin.Plugin.Meta().GetID() {
				items = append(items, plugin)
			}
		}
	}
	return items
}
