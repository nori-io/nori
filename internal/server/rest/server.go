/*
Copyright 2019-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rest

import (
	"sync"

	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/plugins"
	"github.com/nori-io/nori/internal/server"
)

type restServer struct {
	chShutdown chan struct{}
	pm         *plugins.Manager
	secure     bool
	log        logger.Logger
	wg         *sync.WaitGroup
}

func NewServer(pm *plugins.Manager, log logger.Logger, wg *sync.WaitGroup) server.Server {
	return &restServer{
		chShutdown: make(chan struct{}, 1),
		pm:         pm,
		log:        log,
		wg:         wg,
	}
}

func (s *restServer) Start() error {
	go func(s *restServer) {
		for {
			select {
			case <-s.chShutdown:
				// GracefulStop()
			}
		}
	}(s)
	return nil
}

func (s *restServer) Stop() error {
	close(s.chShutdown)
	return nil
}
