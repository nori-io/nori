package file

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/repository"
)

type Manager struct {
	fileRepository repository.FileRepository
	logger         logger.Logger
}

func (m *Manager) Get(file entity.File) (*entity.Plugin, error) {
	return m.fileRepository.Get(file)
}

func (m *Manager) GetAll(files []*entity.File) ([]*entity.Plugin, error) {
	var items []*entity.Plugin
	for _, f := range files {
		p, err := m.fileRepository.Get(*f)
		if err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (m *Manager) File(path string) (*entity.File, error) {
	return m.fileRepository.File(path)
}

func (m *Manager) Dir(dir string) ([]*entity.File, error) {
	return m.fileRepository.Dir(dir)
}

func (m *Manager) Dirs(dirs []string) ([]*entity.File, error) {
	var (
		files []*entity.File
	)
	for _, dir := range dirs {
		fs, err := m.fileRepository.Dir(dir)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	return files, nil
}
