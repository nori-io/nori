package service

import (
	"context"

	"github.com/nori-io/nori-common/v2/meta"
)

type StorageService interface {
	Install(ctx context.Context, m meta.Meta) error
	UnInstall(ctx context.Context, m meta.Meta) error
}
