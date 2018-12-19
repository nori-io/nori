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

package meta

import "github.com/secure2work/nori/version"

type Interface int

const (
	Custom Interface = iota
	Auth
	Authorize
	Cache
	Config
	HTTP
	HTTPTransport
	Mail
	PubSub
	Session
	SQL
	Templates
	Transport
)

var interfaceNames = [...]string{
	"Custom",
	"Auth",
	"Authorize",
	"Cache",
	"Config",
	"HTTP",
	"HTTPTransport",
	"Mail",
	"PubSub",
	"Session",
	"SQL",
	"Templates",
	"Transport",
}

func (i Interface) String() string {
	if i < 0 || int(i) >= len(interfaceNames) {
		return ""
	}
	return interfaceNames[i]
}

func (i Interface) Dependency() Dependency {
	return Dependency{
		ID:         PluginID(i.String()),
		Constraint: "~" + version.NoriMajorVersion,
		Interface:  i,
	}
}
