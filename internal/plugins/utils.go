package plugins

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
)

func LogFieldsMeta(m meta.Meta) []logger.Field {
	return []logger.Field{
		{"plugin", string(m.Id().String())},
		{"interface", m.GetInterface().String()},
	}
}
