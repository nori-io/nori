package files_test

import (
	"path"
	"runtime"
	"testing"

	"github.com/nori-io/nori/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// todo: build plugin before test
func TestPluginExtractor_Get(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := NewFilesLoader()
	file, err := pe.Get(path.Join(dir, "testdata/plugin.so"))

	a.Nil(err)
	a.NotNil(file)
}

func TestPluginExtractor_Get_FileDoesNotExist(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := NewFilesLoader()
	filePath := path.Join(dir, "testdata/no_file.so")
	file, err := pe.Get(filePath)

	a.Error(err)
	a.IsType(errors.FileDoesNotExist{}, err)
	errTyped := err.(errors.FileDoesNotExist)
	a.Equal(filePath, errTyped.Path)
	a.Nil(file)
}

func TestPluginExtractor_Get_FileOpenError(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := NewFilesLoader()
	// todo: create empty plugin file before test
	file, err := pe.Get(path.Join(dir, "testdata/empty.so"))

	a.Error(err)
	a.IsType(errors.FileOpenError{}, err)
	a.Nil(file)
}

func TestPluginExtractor_Get_LookupError(t *testing.T) {
	a := assert.New(t)
	dir := getTestDataDirPath()

	pe := NewFilesLoader()
	// todo: build plugin before test
	file, err := pe.Get(path.Join(dir, "testdata/no_variable.so"))

	a.Error(err)
	a.IsType(errors.LookupError{}, err)
	a.Nil(file)
}

func TestPluginExtractor_Get_NoPluginInterfaceError(t *testing.T) {
	a := assert.New(t)

	dir := getTestDataDirPath()
	pe := NewFilesLoader()
	// todo: build plugin before test
	file, err := pe.Get(path.Join(dir, "/testdata/no_interface.so"))

	a.Error(err)
	a.IsType(errors.NoPluginInterfaceError{}, err)
	a.Nil(file)
}

func getTestDataDirPath() string {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	return dir
}
