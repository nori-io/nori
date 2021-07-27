/*
Copyright 2018-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package file_test

import (
	"path"
	"runtime"
	"testing"

	"github.com/nori-io/nori/internal/repository/file"
	nori_errors "github.com/nori-io/nori/pkg/nori/domain/errors"
	"github.com/stretchr/testify/assert"
)

// todo: build plugin before test
func TestPluginExtractor_Get(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := file.New()
	file, err := pe.Find(path.Join(dir, "testdata/plugin.so"))

	a.Nil(err)
	a.NotNil(file)
}

func TestPluginExtractor_Get_FileDoesNotExist(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := file.New()
	filePath := path.Join(dir, "testdata/no_file.so")
	file, err := pe.Find(filePath)

	a.Error(err)
	a.IsType(nori_errors.FileDoesNotExist{}, err)
	errTyped := err.(nori_errors.FileDoesNotExist)
	a.Equal(filePath, errTyped.Path)
	a.Nil(file)
}

func TestPluginExtractor_Get_FileOpenError(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := file.New()
	// todo: create empty plugin file before test
	file, err := pe.Find(path.Join(dir, "testdata/empty.so"))

	a.Error(err)
	a.IsType(nori_errors.FileOpenError{}, err)
	a.Nil(file)
}

func TestPluginExtractor_Get_LookupError(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := file.New()
	// todo: build plugin before test
	file, err := pe.Find(path.Join(dir, "testdata/no_variable.so"))

	a.Error(err)
	a.IsType(nori_errors.LookupError{}, err)
	a.Nil(file)
}

func TestPluginExtractor_Get_NoPluginInterfaceError(t *testing.T) {
	a := assert.New(t)

	dir := getTestDataDirPath()
	pe := file.New()
	// todo: build plugin before test
	file, err := pe.Find(path.Join(dir, "/testdata/no_interface.so"))

	a.Error(err)
	a.IsType(nori_errors.NoPluginInterfaceError{}, err)
	a.Nil(file)
}

func getTestDataDirPath() string {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	return dir
}
