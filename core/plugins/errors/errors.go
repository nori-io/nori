package errors

import (
	"fmt"
	"strings"

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

type DependenciesNotFound struct {
	Dependencies map[meta.ID][]meta.Dependency
}

func (e DependenciesNotFound) Error() string {
	var msg []string
	for id, deps := range e.Dependencies {
		var msgs []string
		for _, dep := range deps {
			msgs = append(msgs, DependencyNotFound{
				Dependency: dep,
			}.Error())
		}
		msg = append(msg, id.String()+"\n"+strings.Join(msgs, "\n"))
	}
	return strings.Join(msg, "\n")
}

type DependencyCycleFound struct {
	DependencyCycle []meta.Dependency
}

func (e DependencyCycleFound) Error() string {
	var deps []string
	for _, d := range e.DependencyCycle {
		deps = append(deps, d.String())
	}
	return strings.Join(deps, "\n")
}

type InterfaceAssertError struct {
	Interface meta.Interface
}

func (e InterfaceAssertError) Error() string {
	return fmt.Sprintf("can's assert %s to interface %s", e.Interface.String())
}
