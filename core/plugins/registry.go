package plugins

import (
	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/sirupsen/logrus"
)

type RegistryManager interface {
	Registry

	Add(p Plugin)
	Remove(p Plugin)
}

type Registry interface {
	Resolve(dep meta.Dependency) interface{}

	Auth() interfaces.Auth
	Authorize() interfaces.Authorize
	Cache() interfaces.Cache
	Config() config.Manager
	Http() interfaces.Http
	Logger(meta meta.Meta) *logrus.Logger
	Mail() interfaces.Mail
	PubSub() interfaces.PubSub
	Session() interfaces.Session
	Sql() interfaces.SQL
	Templates() interfaces.Templates
}

type registry struct {
	log           *logrus.Logger
	plugins       map[string][]Plugin
	interfaces    map[meta.Interface]string
	configManager config.Manager
}

func NewRegistry(cm config.Manager, logger *logrus.Logger) RegistryManager {
	return &registry{
		log:           logger,
		plugins:       map[string][]Plugin{},
		interfaces:    map[meta.Interface]string{},
		configManager: cm,
	}
}

func (r registry) Add(p Plugin) {
	id := p.GetMeta().Id().ID
	_, ok := r.plugins[id]
	if !ok {
		r.plugins[id] = []Plugin{}
	}
	r.plugins[id] = append(r.plugins[id], p)

	// add alias for non Custom interface
	if p.GetMeta().GetInterface() != meta.Custom {
		r.interfaces[p.GetMeta().GetInterface()] = p.GetMeta().Id().ID
	}
}

func (r registry) Remove(p Plugin) {
	id := p.GetMeta().Id().ID
	_, ok := r.plugins[id]
	if !ok {
		return
	}

	for i, plug := range r.plugins[id] {
		if p.GetMeta().Id() == plug.GetMeta().Id() {
			r.plugins[id] = append(r.plugins[id][:i], r.plugins[id][i+1:]...)
			return
		}
	}

	// remove alias for non Custom interface
	if p.GetMeta().GetInterface() != meta.Custom {
		delete(r.interfaces, p.GetMeta().GetInterface())
	}
}

func (r registry) Resolve(dep meta.Dependency) interface{} {
	plugins, ok := r.plugins[dep.ID]
	if !ok {
		return nil
	}

	for _, plug := range plugins {
		cons, err := dep.GetConstraint()
		if err != nil {
			return nil
		}

		ver, err := plug.GetMeta().Id().GetVersion()
		if err != nil {
			return nil
		}

		if cons.Check(ver) {
			// @todo plugin must run, otherwise return nil (or not?)
			return plug.GetInstance()
		}
	}

	return nil
}

func (r registry) get(alias meta.Interface) interface{} {
	id, ok := r.interfaces[alias]
	if !ok {
		return nil
	}

	plugins, ok := r.plugins[id]
	if !ok {
		return nil
	}

	if len(plugins) == 0 {
		return nil
	}

	return plugins[0]
}

func (r registry) Auth() interfaces.Auth {
	item := r.get(meta.Auth)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Auth)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Auth", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Authorize() interfaces.Authorize {
	item := r.get(meta.Authorize)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Authorize)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Authorize", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Cache() interfaces.Cache {
	item := r.get(meta.Cache)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Cache)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Cache", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Config() config.Manager {
	return r.configManager
}

func (r registry) Http() interfaces.Http {
	item := r.get(meta.HTTP)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Http)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Http", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Logger(meta meta.Meta) *logrus.Logger {
	return r.log
}

func (r registry) Mail() interfaces.Mail {
	item := r.get(meta.Mail)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Mail)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Mail", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) PubSub() interfaces.PubSub {
	item := r.get(meta.PubSub)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.PubSub)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to PubSub", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Session() interfaces.Session {
	item := r.get(meta.Session)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Session)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Session", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Sql() interfaces.SQL {
	item := r.get(meta.SQL)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.SQL)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to SQL", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}

func (r registry) Templates() interfaces.Templates {
	item := r.get(meta.Templates)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Templates)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Templates", item.(Plugin).GetMeta().Id().String())
		return nil
	}
	return i
}
