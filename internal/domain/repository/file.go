package repository

import "github.com/nori-io/nori/internal/domain/entity"

type FileRepository interface {
	// Find returns *entity.Plugin for provided plugin file, return error if file is not a nori Plugin
	Find(file string) (*entity.File, error)
	// FindAll returns entity.Plugin structure for each plugin found in array of file paths
	FindAll(path ...string) ([]entity.File, error)
}
