package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	stdplugin "plugin"

	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/pkg/errors"
)

type FileRepository struct {
	plugins []*entity.Plugin
}

func (r *FileRepository) Get(file entity.File) (*entity.Plugin, error) {
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return nil, errors.FileDoesNotExist{
			Path: file.Path,
			Err:  err,
		}
	}

	f, err := stdplugin.Open(file.Path)
	if err != nil {
		e := errors.FileOpenError{
			Path: file.Path,
			Err:  err,
		}
		return nil, e
	}

	instance, err := f.Lookup("Plugin")
	if err != nil {
		e := errors.LookupError{
			Path: file.Path,
			Err:  err,
		}
		return nil, e
	}

	p, ok := instance.(plugin.Plugin)
	if !ok {
		e := errors.NoPluginInterfaceError{
			Path: file.Path,
		}
		return nil, e
	}

	return entity.NewPlugin(file, p), nil
}

func (r *FileRepository) GetAll(files []entity.File) ([]*entity.Plugin, error) {
	items := []*entity.Plugin{}
	for _, file := range files {
		p, err := r.Get(file)
		if err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *FileRepository) File(path string) (*entity.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.FileDoesNotExist{
			Path: path,
			Err:  err,
		}
	}

	if filepath.Ext(path) != ".so" {
		return nil, errors.FileExtError{
			Path: path,
		}
	}

	return &entity.File{
		Path: path,
	}, nil
}

func (r *FileRepository) Dir(dir string) ([]*entity.File, error) {
	var (
		err   error
		dirs  []os.FileInfo
		files []*entity.File
	)
	if dirs, err = ioutil.ReadDir(dir); err != nil {
		return nil, err
	}
	for _, d := range dirs {
		if d.IsDir() {
			continue
		}
		if filepath.Ext(d.Name()) != ".so" {
			continue
		}
		file, err := r.File(filepath.Join(dir, d.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}
