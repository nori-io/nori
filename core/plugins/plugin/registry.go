package plugin

import (
	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/sirupsen/logrus"
)

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
