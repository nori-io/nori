package plugins

import (
	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori-common/meta"
)

func LogFieldsMeta(m meta.Meta) []logger.Field {
	return []logger.Field{
		{"plugin_id", string(m.Id().ID)},
		{"plugin_version", m.Id().Version},
		{"interface", m.GetInterface().String()},
	}
}
