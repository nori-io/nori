package service

import (
	"bytes"

	"github.com/nori-io/nori/internal/domain/entity"
)

type FileService interface {
	Open(name string, src bytes.Buffer) (*entity.File, error)
	Delete(file string) error

	Get(file string) (*entity.File, error)
	GetAll(dir string) ([]entity.File, error)
}
