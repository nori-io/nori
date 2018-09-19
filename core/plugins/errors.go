// Copyright © 2018 Secure2Work info@secure2work.com
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
	"errors"
	"fmt"
	"strings"
)

var PluginOpenError = errors.New("can't open plugin file")
var PluginLookupError = errors.New("can't lookup Plugin variable")
var PluginInterfaceError = errors.New("can't match Plugin variable to Plugin interface")
var PluginHashError = errors.New("can't calculate hash for plugin file")
var PluginNamespaceError = errors.New("can't parse plugin namespace")

type DependencyError struct {
	PlugName      string
	PlugVer       string
	DepName       string
	DepConstraint string
}

func (e DependencyError) Error() string {
	return fmt.Sprintf("%s:%s has unresolved dependency %s:%s", e.PlugName, e.PlugVer, e.DepName, e.DepConstraint)
}

// PluginNotFound
type PluginNotFound struct {
	PluginId string
}

func (p *PluginNotFound) Error() string {
	return fmt.Sprintf("plugin with id [%s] not found", p.PluginId)
}

// Plugin has dependent plugins
type PluginHasDependentPlugins struct {
	PluginId     string
	Dependencies []string
}

func (p *PluginHasDependentPlugins) Error() string {
	return fmt.Sprintf("plugin [%s] in dependencies in plugin(s) [%s]",
		p.PluginId,
		strings.Join(p.Dependencies, ","))
}
