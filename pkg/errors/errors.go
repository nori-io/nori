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

package errors

import (
	"fmt"
	"strings"

	"github.com/nori-io/nori-common/v2/meta"
)

type NotFound struct {
	ID meta.ID
}

func (e NotFound) Error() string {
	return fmt.Sprintf("plugin [%s] not found", e.ID.String())
}

type AlreadyExists struct {
	ID   meta.ID
	Path string
	Err  error
}

func (e AlreadyExists) Error() string {
	return fmt.Sprintf("plugin [%s] already exists {%s}", e.ID.String())
}

type FileDoesNotExist struct {
	Path string
	Err  error
}

func (e FileDoesNotExist) Error() string {
	return fmt.Sprintf("file [%s] does not exist", e.Path)
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

type NoPluginInterfaceError struct {
	Path string
}

func (e NoPluginInterfaceError) Error() string {
	return fmt.Sprintf("plugin [%s] does not implement Plugin interface", e.Path)
}

type NonInstallablePlugin struct {
	ID   meta.ID
	Path string
}

func (e NonInstallablePlugin) Error() string {
	return fmt.Sprintf("non-installable plugin [%s] in %s", e.ID.String(), e.Path)
}

type IncompatibleCoreVersion struct {
	ID                 meta.ID
	NeededCoreVersion  string
	CurrentCoreVersion string
}

func (e IncompatibleCoreVersion) Error() string {
	return fmt.Sprintf("Plugin [%s] requires Nori [%s], running Nori [%s]",
		e.ID.String(), e.NeededCoreVersion, e.CurrentCoreVersion)
}

type InterfaceNotFound struct {
	Interface meta.Interface
}

func (e InterfaceNotFound) Error() string {
	return fmt.Sprintf("Interface %s is nil", e.Interface)
}

type DependencyNotFound struct {
	Dependency meta.Dependency
}

func (e LoopVertexFound) Error() string {
	return fmt.Sprintf("LoopVertex [%s][%s] found",
		e.Dependency.Interface, e.Dependency.Constraint)
}

type LoopVertexFound struct {
	Dependency meta.Dependency
}

func (e DependencyNotFound) Error() string {
	return fmt.Sprintf("Dependency [%s][%s] not found",
		e.Dependency.Interface, e.Dependency.Constraint)
}

type DependenciesNotFound struct {
	Dependencies map[meta.ID][]meta.Dependency
}

func (e DependenciesNotFound) Add(id meta.ID, dep meta.Dependency) {
	if e.Dependencies == nil {
		e.Dependencies = map[meta.ID][]meta.Dependency{}
	}
	e.Dependencies[id] = append(e.Dependencies[id], dep)
}

func (e DependenciesNotFound) HasErrors() bool {
	return len(e.Dependencies) > 0
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
