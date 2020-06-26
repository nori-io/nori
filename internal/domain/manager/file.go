package manager

import "github.com/nori-io/nori/internal/domain/entity"

type File interface {
	Get(file entity.File) (*entity.Plugin, error)
	GetAll(files []*entity.File) ([]*entity.Plugin, error)

	File(path string) (*entity.File, error)
	Dir(dir string) ([]*entity.File, error)
	Dirs(dir []string) ([]*entity.File, error)
}
