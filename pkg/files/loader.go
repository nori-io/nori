// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package files

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	stdplugin "plugin"

	"github.com/nori-io/nori/pkg/types"

	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/pkg/errors"
)

/**
PluginFilesLoader get path to plugin files from provided list of directories
and get plugins from found files
*/
type FilesLoader interface {
	// Get returns *File structure for provided plugin file, return error if file is not a nori Plugin
	Get(file string) (*types.File, error)
	// GetAll returns *File structure for each nori Plugin provided in path array
	GetAll(path []string) ([]*types.File, error)
	// Files returns file paths to .so file found in provided dirs
	Files(dirs []string) ([]string, error)
}

func NewFilesLoader() FilesLoader {
	return &filesLoader{}
}

type filesLoader struct{}

func (pl *filesLoader) Get(path string) (*types.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.FileDoesNotExist{
			Path: path,
			Err:  err,
		}
	}

	file, err := stdplugin.Open(path)
	if err != nil {
		e := errors.FileOpenError{
			Path: path,
			Err:  err,
		}
		return nil, e
	}

	instance, err := file.Lookup("Plugin")
	if err != nil {
		e := errors.LookupError{
			Path: path,
			Err:  err,
		}
		return nil, e
	}

	p, ok := instance.(plugin.Plugin)
	if !ok {
		e := errors.NoPluginInterfaceError{
			Path: path,
		}
		return nil, e
	}

	return &types.File{
		Plugin: p,
		Path:   path,
	}, nil
}

func (pl *filesLoader) GetAll(paths []string) ([]*types.File, error) {
	files := []*types.File{}
	for _, path := range paths {
		f, err := pl.Get(path)
		if err != nil {
			return []*types.File{}, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (pl *filesLoader) Files(dirs []string) ([]string, error) {
	var err error
	var ret []string
	for _, dir := range dirs {
		var dirs []os.FileInfo
		if dirs, err = ioutil.ReadDir(dir); err != nil {
			return ret, err
		}
		for _, d := range dirs {
			if d.IsDir() {
				continue
			}
			if path.Ext(d.Name()) != ".so" {
				continue
			}
			filePath := filepath.Join(dir, d.Name())
			ret = append(ret, filePath)
		}
	}
	return ret, nil
}
