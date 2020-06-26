package repository

import "github.com/nori-io/nori/internal/domain/entity"

type FileRepository interface {
	// Get returns *entity.Plugin for provided plugin file, return error if file is not a nori Plugin
	Get(file entity.File) (*entity.Plugin, error)
	// GetAll returns entity.Plugin structure for each plugin found in array of file paths
	GetAll(path []entity.File) ([]*entity.Plugin, error)
	// File returns File entity with path to .so file
	File(path string) (*entity.File, error)
	// Dir returns files found in dir
	Dir(dir string) ([]*entity.File, error)
}
