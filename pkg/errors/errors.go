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
	return fmt.Sprintf("plugin [%s] already exists {%s}", e.ID.String(), e.Path)
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

type FileExtError struct {
	Path string
	Err  error
}

func (e FileExtError) Error() string {
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

func (e DependencyNotFound) Error() string {
	return fmt.Sprintf("Dependency [%s][%s] not found",
		e.Dependency.Interface, e.Dependency.Interface.Constraint())
}

type LoopVertexFound struct {
	Dependency meta.Dependency
}

func (e LoopVertexFound) Error() string {
	return fmt.Sprintf("LoopVertex [%s][%s] found",
		e.Dependency.Interface, e.Dependency.Interface.Constraint())
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

type ConfigFormatError struct {
	Param string
}

func (e ConfigFormatError) Error() string {
	return fmt.Sprintf("config param format error [%s]", e.Param)
}
