package plugins

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
)

func LogFieldsMeta(m meta.Meta) []logger.Field {
	return []logger.Field{
		{"plugin_id", string(m.Id().ID)},
		{"plugin_version", m.Id().Version},
		{"interface", m.GetInterface().String()},
	}
}
