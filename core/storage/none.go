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
	"errors"

	"github.com/nori-io/nori-common/meta"
)

type none struct {
}

func getNoneStorage() (Storage, error) {
	return none{}, nil
}

func (n none) GetPluginMetas() ([]meta.Meta, error) {
	return []meta.Meta{}, nil
}

func (n none) SavePluginMeta(meta meta.Meta) error {
	return errors.New("Can't save to None storage")
}

func (n none) DeletePluginMeta(id meta.ID) error {
	return errors.New("Can't delete from None storage")
}
