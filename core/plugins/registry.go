package plugins

import (
	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/errors"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"
	"github.com/sirupsen/logrus"
)

type registry struct {
	registryManager RegistryManager
	log             *logrus.Logger
	configManager   config.Manager
}

func NewRegistry(rm RegistryManager, cm config.Manager, logger *logrus.Logger) plugin.Registry {
	return registry{
		log:             logger,
		registryManager: rm,
		configManager:   cm,
	}
}

func (r registry) Resolve(dep meta.Dependency) (interface{}, error) {
	plugin := r.registryManager.Resolve(dep)
	if plugin != nil {
		return plugin.Instance(), nil
	}
	return nil, errors.DependencyNotFound{
		Dependency: dep,
	}
}

func (r registry) Auth() (interfaces.Auth, error) {
	item, err := r.registryManager.GetInterface(meta.Auth)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Auth)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Auth,
		}
	}
	return i, nil
}

func (r registry) Authorize() (interfaces.Authorize, error) {
	item, err := r.registryManager.GetInterface(meta.Authorize)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Authorize)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Authorize,
		}
	}
	return i, nil
}

func (r registry) Cache() (interfaces.Cache, error) {
	item, err := r.registryManager.GetInterface(meta.Cache)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Cache)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Cache,
		}
	}
	return i, nil
}

func (r registry) Config() config.Manager {
	return r.configManager
}

func (r registry) Http() (interfaces.Http, error) {
	item, err := r.registryManager.GetInterface(meta.HTTP)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Http)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTP,
		}
	}
	return i, nil
}

func (r registry) HTTPTransport() (interfaces.HTTPTransport, error) {
	item, err := r.registryManager.GetInterface(meta.HTTPTransport)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.HTTPTransport)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Logger(meta meta.Meta) *logrus.Logger {
	return r.log.WithField("plugin", meta.Id().String()).Logger
}

func (r registry) Mail() (interfaces.Mail, error) {
	item, err := r.registryManager.GetInterface(meta.Mail)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Mail)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Mail,
		}
	}
	return i, nil
}

func (r registry) PubSub() (interfaces.PubSub, error) {
	item, err := r.registryManager.GetInterface(meta.PubSub)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.PubSub)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.PubSub,
		}
	}
	return i, nil
}

func (r registry) Session() (interfaces.Session, error) {
	item, err := r.registryManager.GetInterface(meta.Session)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Session)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Sql() (interfaces.SQL, error) {
	item, err := r.registryManager.GetInterface(meta.SQL)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.SQL)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Templates() (interfaces.Templates, error) {
	item, err := r.registryManager.GetInterface(meta.Templates)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Templates)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}
