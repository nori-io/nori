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
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/secure2work/nori/core/endpoint"
)

type Server struct {
	e            endpoint.Endpoint
	decode       DecodeRequest
	encode       EncodeResponse
	before       []ServerBeforeFunc
	after        []ServerAfterFunc
	errorEncoder ErrorEncoder
	logger       *logrus.Logger
}

type ServerOption func(*Server)

func NewServer(
	e endpoint.Endpoint,
	decode DecodeRequest,
	encode EncodeResponse,
	logger *logrus.Logger,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:            e,
		decode:       decode,
		encode:       encode,
		errorEncoder: DefaultErrorEncoder,
		logger:       logger,
	}
	for _, option := range options {
		option(s)
	}
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	for _, f := range s.before {
		ctx = f(ctx, r)
	}

	request, err := s.decode(ctx, r)
	if err != nil {
		s.logger.Error(err)
		s.errorEncoder(ctx, err, w)
		return
	}

	response, err := s.e(ctx, request)

	if err != nil {
		s.logger.Error(err)
		s.errorEncoder(ctx, err, w)
		return
	}

	for _, f := range s.after {
		ctx = f(ctx, w)
	}

	if err := s.encode(ctx, w, response); err != nil {
		s.logger.Error(err)
		s.errorEncoder(ctx, err, w)
		return
	}
}
