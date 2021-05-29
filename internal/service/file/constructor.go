package file

import (
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env            *env.Env
	FileRepository repository.FileRepository
}

func New(params Params) service.FileService {
	return &Service{
		Env:            params.Env,
		FileRepository: params.FileRepository,
	}
}
