package plugin

import (
	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/sirupsen/logrus"
)

type Registry interface {
	Resolve(dep meta.Dependency) (interface{}, error)

	Auth() (interfaces.Auth, error)
	Authorize() (interfaces.Authorize, error)
	Cache() (interfaces.Cache, error)
	Config() config.Manager
	Http() (interfaces.Http, error)
	HTTPTransport() (interfaces.HTTPTransport, error)
	Logger(meta meta.Meta) *logrus.Logger
	Mail() (interfaces.Mail, error)
	PubSub() (interfaces.PubSub, error)
	Session() (interfaces.Session, error)
	Sql() (interfaces.SQL, error)
	Templates() (interfaces.Templates, error)
}
