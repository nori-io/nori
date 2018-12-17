package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"

	"github.com/secure2work/nori/core/plugins/meta"

	"github.com/sirupsen/logrus"
)

type Manager interface {
	AddFile(path string) (Plugin, error)
	AddDir(paths []string) error

	Install(id meta.ID, ctx context.Context) error

	Meta(id meta.ID) (meta.Meta, error)
	Metas() []meta.Meta

	Start(id meta.ID, ctx context.Context) error
	StartAll(ctx context.Context) error

	Stop(id meta.ID, ctx context.Context) error
	StopAll(ctx context.Context) error

	UnInstall(id meta.ID, ctx context.Context) error
}

func NewManager(registry RegistryManager, logger *logrus.Logger) Manager {
	return &manager{
		files:    map[string]meta.ID{},
		log:      logger,
		plugins:  map[meta.ID]Plugin{},
		registry: registry,
	}
}

type manager struct {
	files    map[string]meta.ID
	log      *logrus.Logger
	plugins  map[meta.ID]Plugin
	registry RegistryManager
}

func (m *manager) AddFile(path string) (Plugin, error) {
	file, err := plugin.Open(path)
	if err != nil {
		e := FileOpenError{
			Path: path,
			Err:  err,
		}
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	instance, err := file.Lookup("Plugin")
	if err != nil {
		e := LookupError{
			Path: path,
			Err:  err,
		}
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	p, ok := instance.(Plugin)
	if !ok {
		e := TypeAssertError{
			Path: path,
		}
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	if p.GetMeta().GetInterface() != meta.Custom {
		m.registry.Add(p)
	}

	m.plugins[p.GetMeta().Id()] = p
	m.files[path] = p.GetMeta().Id()

	return p, nil
}

func (m *manager) AddDir(paths []string) error {
	var err error
	for _, dir := range paths {
		var dirs []os.FileInfo
		if dirs, err = ioutil.ReadDir(dir); err != nil {
			return err
		}
		for _, d := range dirs {
			if d.IsDir() {
				continue
			}
			if path.Ext(d.Name()) != ".so" {
				continue
			}

			// load plugin
			filePath := filepath.Join(dir, d.Name())
			p, err := m.AddFile(filePath)
			if err != nil {
				m.log.Error(err)
				continue
			}

			m.log.Infof(
				"Found '%s' by '%s'",
				p.GetMeta().Id().String(),
				p.GetMeta().GetAuthor().Name,
			)
		}
	}
	return nil
}

func (m *manager) Install(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return NotFound{ID: id}
	}
	return p.Install(ctx, nil)
}

func (m *manager) Meta(id meta.ID) (meta.Meta, error) {
	p, ok := m.plugins[id]
	if !ok {
		return nil, NotFound{
			ID: id,
		}
	}
	return p.GetMeta(), nil
}

func (m *manager) Metas() []meta.Meta {
	var meta []meta.Meta
	for _, p := range m.plugins {
		meta = append(meta, p.GetMeta())
	}
	return meta
}

func (m *manager) Start(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return NotFound{ID: id}
	}
	// todo replace nil
	err := p.Init(ctx, nil)
	if err != nil {
		return err
	}
	return p.Start(ctx, nil)
}

func (m *manager) StartAll(ctx context.Context) error {
	// todo start plugins in topological order
	return nil
}

func (m *manager) Stop(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return NotFound{ID: id}
	}
	return p.Stop(ctx, nil)
}

func (m *manager) StopAll(ctx context.Context) error {
	// todo stop plugins in topological order
	return nil
}

func (m *manager) UnInstall(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return NotFound{ID: id}
	}
	return p.UnInstall(ctx, nil)
}
