package file

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/domain/manager"
	"github.com/nori-io/nori/internal/domain/repository"
	"go.uber.org/fx"
)

type ManagerParams struct {
	fx.In

	FileRepository repository.FileRepository
	Logger         logger.Logger
}

func New(params ManagerParams) (manager.File, error) {
	return &Manager{
		fileRepository: params.FileRepository,
		logger:         params.Logger,
	}, nil
}
