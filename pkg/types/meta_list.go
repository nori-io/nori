// Copyright © 2018-2019 Nori info@nori.io
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
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/pkg/errors"
)

type MetaList []meta.Meta

func (ml MetaList) Find(id meta.ID) (meta.Meta, error) {
	for _, m := range ml {
		if m.Id() == id {
			return m, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (ml *MetaList) Add(m meta.Meta) {
	*ml = append(*ml, m)
}

func (ml *MetaList) Remove(id meta.ID) {
	for i, v := range *ml {
		if v.Id() == id {
			*ml = append((*ml)[:i], (*ml)[i+1:]...)
		}
	}
}
