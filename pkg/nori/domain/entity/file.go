package entity

import (
	"github.com/nori-io/common/v5/pkg/domain/plugin"
)

type File struct {
	Path string
	Fn   func() plugin.Plugin
}
