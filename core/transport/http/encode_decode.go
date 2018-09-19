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
	"encoding/json"
	"net/http"
)

type StatusCoder interface {
	StatusCode() int
}

type Headerer interface {
	Headers() http.Header
}

type DecodeRequest func(context.Context, *http.Request) (request interface{}, err error)

type EncodeResponse func(context.Context, http.ResponseWriter, interface{}) error

// This function can be called before main
func EncodeWrapper(
	ctx context.Context,
	w http.ResponseWriter,
	response interface{},
	encode EncodeResponse,
) error {
	// Headerer interface
	if headerer, ok := response.(Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}

	// StatusCoder interface
	code := http.StatusOK
	if sc, ok := response.(StatusCoder); ok {
		code = sc.StatusCode()
		w.WriteHeader(code)
	}
	if code == http.StatusNoContent {
		return nil
	}

	return encode(ctx, w, response)
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return EncodeWrapper(ctx, w, response, func(context.Context, http.ResponseWriter, interface{}) error {
		return json.NewEncoder(w).Encode(response)
	})
}
