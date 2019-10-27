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

package plugins

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	stdplugin "plugin"

	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/internal/errors"
)

/**
PluginExtractor get path to plugin files from provided list of directories
and extracts plugins from found files
*/
type PluginExtractor interface {
	// get plugin from file
	Get(file string) (plugin.Plugin, error)
	// get plugin file paths from provided dirs
	Files(dirs []string) ([]string, error)
}

func NewPluginExtractor() PluginExtractor {
	return &pluginExtractor{}
}

type pluginExtractor struct{}

func (pl *pluginExtractor) Get(path string) (plugin.Plugin, error) {
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

	return p, nil
}

func (pl *pluginExtractor) Files(dirs []string) ([]string, error) {
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
