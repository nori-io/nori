package hooks

import (
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env         *env.Env
	FileService service.FileService
}

func New(params Params) (service.HooksService, error) {
	files, err := params.FileService.GetAll(params.Env.Config.Nori.Hooks.Dir)
	if err != nil {
		return nil, err
	}
	for i := range files {
		p := files[i].Fn()
		hook, ok := p.(logger.Hook)
		if !ok {
			continue
		}
		params.Env.Logger.AddHook(hook)
	}

	return Service{}, nil
}
