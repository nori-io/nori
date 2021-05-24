package file

import (
	"bytes"
	"errors"
	"os"
	"path"

	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/env"
)

type Service struct {
	Env            *env.Env
	FileRepository repository.FileRepository
}

func (s *Service) Open(name string, src bytes.Buffer) (*entity.File, error) {
	if s.Env.Config.Nori.Plugins.Dir == "" {
		// todo: error
		return nil, errors.New("plugin directory not defined")
	}

	name = path.Join(s.Env.Config.Nori.Plugins.Dir, name)

	err := os.WriteFile(name, src.Bytes(), 0644)
	if err != nil {
		return nil, err
	}

	return s.Get(name)
}

func (s *Service) Delete(file string) error {
	// todo: name must be inside plugins.dir
	return os.Remove(file)
}

func (s *Service) Get(file string) (*entity.File, error) {
	return s.FileRepository.Find(file)
}

func (s *Service) GetAll(dir string) ([]entity.File, error) {
	return s.FileRepository.FindAll(dir)
}
