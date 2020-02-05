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

package types

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori/pkg/errors"
)

type PluginList []plugin.Plugin

func (pl *PluginList) Add(p plugin.Plugin) error {
	if p, _ := pl.ID(p.Meta().Id()); p != nil {
		return errors.AlreadyExists{
			ID: p.Meta().Id(),
		}
	}
	*pl = append(*pl, p)
	return nil
}

func (pl *PluginList) ID(id meta.ID) (plugin.Plugin, error) {
	for _, p := range *pl {
		if p.Meta().Id() == id {
			return p, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (pl *PluginList) Interface(i meta.Interface) (plugin.Plugin, error) {
	for _, p := range *pl {
		if p.Meta().GetInterface().Equal(i) {
			return p, nil
		}
	}
	return nil, errors.InterfaceNotFound{Interface: i}
}

func (pl *PluginList) Remove(p plugin.Plugin) {
	for i, v := range *pl {
		if v == p {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}

func (pl *PluginList) Resolve(dep meta.Dependency) (plugin.Plugin, error) {
	cons, err := dep.GetConstraint()
	if err != nil {
		return nil, err
	}

	for _, p := range *pl {
		if dep.ID != p.Meta().Id().ID {
			continue
		}

		v, _ := p.Meta().Id().GetVersion()

		if cons.Check(v) {
			return p, nil
		}
	}
	return nil, errors.DependencyNotFound{Dependency: dep}
}
