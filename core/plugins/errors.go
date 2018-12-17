package plugins

import (
	"fmt"

	"github.com/secure2work/nori/core/plugins/meta"
)

type NotFound struct {
	ID meta.ID
}

func (e NotFound) Error() string {
	return fmt.Sprintf("plugin [%s] not found", e.ID.String())
}

type FileOpenError struct {
	Path string
	Err  error
}

func (e FileOpenError) Error() string {
	return fmt.Sprintf("error on file [%s] opening: %s", e.Path, e.Err.Error())
}

type LookupError struct {
	Path string
	Err  error
}

func (e LookupError) Error() string {
	return fmt.Sprintf("error on lookup in [%s]: %s", e.Path, e.Err.Error())
}

type TypeAssertError struct {
	Path string
}

func (e TypeAssertError) Error() string {
	return fmt.Sprintf("plugin [%s] does not implement Plugin interface", e.Path)
}

type UnknownInterface struct {
	Path string
}

func (e UnknownInterface) Error() string {
	return fmt.Sprintf("plugin [%s] implements unknown interface", e.Path)
}
