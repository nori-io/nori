package plugins

import (
	"github.com/secure2work/nori/core/plugins/meta"
)

type FileTable map[string]meta.ID

func (ft FileTable) Find(id meta.ID) string {
	for p, i := range ft {
		if i == id {
			return p
		}
	}
	return ""
}
