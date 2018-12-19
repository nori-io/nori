package plugins

import (
	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"
	"github.com/sirupsen/logrus"
)

type registry struct {
	rm            RegistryManager
	log           *logrus.Logger
	configManager config.Manager
}

func NewRegistry(rm RegistryManager, cm config.Manager, logger *logrus.Logger) plugin.Registry {
	return registry{
		log:           logger,
		rm:            rm,
		configManager: cm,
	}
}

func (r registry) Resolve(dep meta.Dependency) interface{} {
	plugin := r.rm.Resolve(dep)
	if plugin != nil {
		return plugin.Instance()
	}
	return nil
}

func (r registry) Auth() interfaces.Auth {
	item := r.rm.GetInterface(meta.Auth)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Auth)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Auth", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Authorize() interfaces.Authorize {
	item := r.rm.GetInterface(meta.Authorize)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Authorize)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Authorize", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Cache() interfaces.Cache {
	item := r.rm.GetInterface(meta.Cache)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Cache)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Cache", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Config() config.Manager {
	return r.configManager
}

func (r registry) Http() interfaces.Http {
	item := r.rm.GetInterface(meta.HTTP)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Http)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Http", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Logger(meta meta.Meta) *logrus.Logger {
	return r.log.WithField("plugin", meta.Id().String()).Logger
}

func (r registry) Mail() interfaces.Mail {
	item := r.rm.GetInterface(meta.Mail)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Mail)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Mail", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) PubSub() interfaces.PubSub {
	item := r.rm.GetInterface(meta.PubSub)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.PubSub)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to PubSub", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Session() interfaces.Session {
	item := r.rm.GetInterface(meta.Session)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Session)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Session", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Sql() interfaces.SQL {
	item := r.rm.GetInterface(meta.SQL)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.SQL)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to SQL", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}

func (r registry) Templates() interfaces.Templates {
	item := r.rm.GetInterface(meta.Templates)
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Templates)
	if !ok {
		r.log.Errorf("Can't convert plugin %s interface to Templates", item.(plugin.Plugin).Meta().Id().String())
		return nil
	}
	return i
}
