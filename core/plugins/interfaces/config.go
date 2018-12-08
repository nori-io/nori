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

package interfaces

type Config interface {
	Bool(key string) bool
	Get(key string) interface{}
	Float(key string) float64
	Int(key string) int
	IsSet(key string) bool
	UInt(key string) uint
	Unmarshal(v interface{}, prefix string) error
	SetDefault(key string, val interface{})
	Slice(key, delimiter string) []interface{}
	String(key string) string
	StringMap(key string) map[string]interface{}
}
