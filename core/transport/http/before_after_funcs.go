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

package http

import (
	"context"
	"net/http"
)

type ServerBeforeFunc func(context.Context, *http.Request) context.Context

type ServerAfterFunc func(context.Context, http.ResponseWriter) context.Context

func ServerBefore(before ...ServerBeforeFunc) ServerOption {
	return func(s *Server) { s.before = append(s.before, before...) }
}

func ServerAfter(after ...ServerAfterFunc) ServerOption {
	return func(s *Server) { s.after = append(s.after, after...) }
}
