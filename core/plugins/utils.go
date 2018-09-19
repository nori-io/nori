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
	"github.com/secure2work/nori/core/entities"
	"sort"
)

type pluginList []PluginEntry

func CheckDependencies(plugins map[string]PluginEntry) []error {
	var err []error
	for _, p := range plugins {
		e := p.checkDependencies(plugins)
		err = append(err, e...)
	}
	return err
}

func SortPlugins(plugins map[string]PluginEntry) []PluginEntry {
	pl := pluginList{}
	for _, p := range plugins {
		p.calcWeight(plugins)
		pl = append(pl, p)
	}
	sort.Sort(pl)
	return pl
}

func (i pluginList) Len() int      { return len(i) }
func (i pluginList) Swap(x, y int) { i[x], i[y] = i[y], i[x] }
func (i pluginList) Less(x, y int) bool {
	// more weight = less priority
	if i[x].getWeight() < i[y].getWeight() {
		return true
	}

	// sorting by name
	if i[x].getWeight() == i[y].getWeight() && i[x].Plugin().GetMeta().GetId() < i[y].Plugin().GetMeta().GetId() {
		return true
	}

	return false
}

func isInstalled(meta entities.PluginMeta, installed []entities.PluginMeta) bool {
	for _, m := range installed {
		// @todo add hash check or not (?)
		if m.GetId() == meta.GetId() && m.GetVersion() == meta.GetVersion() {
			return true
		}
	}
	return false
}
