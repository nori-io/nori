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

package graph

import "github.com/nori-io/nori-common/v2/meta"

type Edge interface {
	From() meta.ID
	To() meta.ID
}

type edge struct {
	from meta.ID
	to   meta.ID
}

func (e *edge) From() meta.ID {
	return e.from
}

func (e *edge) To() meta.ID {
	return e.to
}
