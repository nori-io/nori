package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	stdplugin "plugin"

	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/pkg/errors"
)

type FileRepository struct{}

func (r *FileRepository) Find(file string) (*entity.File, error) {
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return nil, errors.FileDoesNotExist{
			Path: file,
			Err:  err,
		}
	}

	if stat.IsDir() {
		return nil, errors.FileOpenError{
			Path: file,
			Err:  fmt.Errorf("%s is a directory", file),
		}
	}

	if filepath.Ext(file) != ".so" {
		return nil, errors.FileExtError{
			Path: file,
		}
	}

	f, err := stdplugin.Open(file)
	if err != nil {
		e := errors.FileOpenError{
			Path: file,
			Err:  err,
		}
		return nil, e
	}

	symbol, err := f.Lookup("New")
	if err != nil {
		e := errors.LookupError{
			Path: file,
			Err:  err,
		}
		return nil, e
	}

	fn, ok := symbol.(func() plugin.Plugin)
	if !ok {
		e := errors.NoPluginInterfaceError{
			Path: file,
		}
		return nil, e
	}

	return &entity.File{
		Path: file,
		Fn:   fn,
	}, nil
}

func (r *FileRepository) FindAll(paths ...string) ([]entity.File, error) {
	var (
		files []entity.File
	)
	for _, path := range paths {
		if path == "" {
			continue
		}

		stat, err := os.Stat(path)

		if !stat.IsDir() {
			file, err := r.Find(path)
			if err != nil {
				return nil, err
			}
			files = append(files, *file)
			continue
		}

		dirs, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, d := range dirs {
			if d.IsDir() {
				nestedFiles, err := r.FindAll(filepath.Join(path, d.Name()))
				if err != nil {
					return nil, err
				}
				files = append(files, nestedFiles...)
				continue
			}
			if d.Size() == 0 {
				continue
			}
			if filepath.Ext(d.Name()) != ".so" {
				continue
			}
			file, err := r.Find(filepath.Join(path, d.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, *file)
		}
	}
	return files, nil
}
