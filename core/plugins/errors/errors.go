package errors

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

type NonInstallablePlugin struct {
	Id   meta.ID
	Path string
}

func (e NonInstallablePlugin) Error() string {
	return fmt.Sprintf("non-installable plugin [%s] in %s", e.Id.String(), e.Path)
}

type IncompatibleCoreVersion struct {
	Id                 meta.ID
	NeededCoreVersion  string
	CurrentCoreVersion string
}

func (e IncompatibleCoreVersion) Error() string {
	return fmt.Sprintf("Plugin [%s] requires Nori [%s], running Nori [%s]",
		e.Id.String(), e.NeededCoreVersion, e.CurrentCoreVersion)
}

type InterfaceNotFound struct {
	Interface meta.Interface
}

func (e InterfaceNotFound) Error() string {
	return fmt.Sprintf("Interface %s is nil", e.Interface.String())
}

type DependencyNotFound struct {
	Dependency meta.Dependency
}

func (e DependencyNotFound) Error() string {
	return fmt.Sprintf("Dependency [%s][%s] not found",
		e.Dependency.ID, e.Dependency.Constraint)
}

type InterfaceAssertError struct {
	Interface meta.Interface
}

func (e InterfaceAssertError) Error() string {
	return fmt.Sprintf("can's assert %s to interface %s", e.Interface.String())
}
