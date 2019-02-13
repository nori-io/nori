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

package storage

import (
	"github.com/nori-io/nori-common/meta"
)

type none struct {
	plugins Plugins
}

type nonePlugins struct {
	metas map[meta.ID]meta.Meta
}

func getNoneStorage() (Storage, error) {
	return none{
		plugins: &nonePlugins{
			metas: map[meta.ID]meta.Meta{},
		},
	}, nil
}

func (n none) Plugins() Plugins {
	return n.plugins
}

func (n *nonePlugins) All() ([]meta.Meta, error) {
	var metas []meta.Meta
	for _, meta := range n.metas {
		metas = append(metas, meta)
	}
	return []meta.Meta{}, nil
}

func (n *nonePlugins) Save(meta meta.Meta) error {
	n.metas[meta.Id()] = meta
	return nil
}

func (n *nonePlugins) Delete(id meta.ID) error {
	delete(n.metas, id)
	return nil
}
