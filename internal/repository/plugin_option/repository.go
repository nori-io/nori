package plugin_option

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/storage"
	"github.com/nori-io/common/v5/pkg/errors"
	"github.com/nori-io/nori/internal/domain/entity"
)

type Repository struct {
	Bucket storage.Bucket
}

func (r *Repository) Upsert(op entity.PluginOption) error {
	data, err := encode(newModel(op))
	if err != nil {
		return err
	}
	err = r.Bucket.Set(op.ID.String(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(id meta.ID) error {
	return r.Bucket.Delete(id.String())
}

func (r *Repository) Find(id meta.ID) (entity.PluginOption, error) {
	data, err := r.Bucket.Get(id.String())
	if err != nil {
		return entity.PluginOption{}, err
	}

	if len(data) == 0 {
		return entity.PluginOption{}, errors.EntityNotFound{Entity: id.String()}
	}

	m, err := decode(data)
	if err != nil {
		return entity.PluginOption{}, err
	}
	return m.Convert(), nil
}

func (r *Repository) FindAll() ([]entity.PluginOption, error) {
	items := []entity.PluginOption{}

	c := r.Bucket.Cursor()
	if c == nil {
		return nil, nil
	}

	for id, data := c.First(); id != ""; id, data = c.Next() {
		if len(data) == 0 {
			continue
		}
		m, err := decode(data)
		if err != nil {
			return nil, err
		}
		items = append(items, m.Convert())
	}

	return items, nil
}
