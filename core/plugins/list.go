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
	"github.com/secure2work/nori-common/meta"
	"github.com/secure2work/nori-common/plugin"
	"github.com/secure2work/nori/core/errors"
)

type PluginList []plugin.Plugin

func (pl PluginList) Find(id meta.ID) (plugin.Plugin, error) {
	for _, p := range pl {
		if p.Meta().Id() == id {
			return p, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (pl PluginList) FindByPluginID(id meta.PluginID) PluginList {
	list := PluginList{}
	for _, p := range pl {
		if p.Meta().Id().ID == id {
			list = append(list, p)
		}
	}
	return list
}

func (pl PluginList) Resolve(dependency meta.Dependency) plugin.Plugin {
	cons, err := dependency.GetConstraint()
	if err != nil {
		return nil
	}

	for _, p := range pl {
		if dependency.ID != p.Meta().Id().ID {
			continue
		}

		v, _ := p.Meta().Id().GetVersion()

		if cons.Check(v) {
			return p
		}
	}
	return nil
}

func (pl *PluginList) Add(p plugin.Plugin) {
	if p, _ := pl.Find(p.Meta().Id()); p != nil {
		return
	}
	*pl = append(*pl, p)
}

func (pl *PluginList) Remove(p plugin.Plugin) {
	for i, v := range *pl {
		if v == p {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}
