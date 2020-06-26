package service

import (
	"context"

	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type HubClientService interface {
	GetPluginByID(ctx context.Context, id meta.ID) (entity.Plugin, error)
}
