/*
Copyright 2018-2020 The Nori Authors.
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

package grpc

import (
	"sync"

	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/domain/handler"
	"github.com/nori-io/nori/internal/domain/manager"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	chShutdown chan struct{}
	grpcServer *grpc.Server
	pm         *manager.Plugin
	secure     bool
	log        logger.Logger
	wg         *sync.WaitGroup
}

func NewServer(pm *manager.Plugin, log logger.Logger, wg *sync.WaitGroup) handler.Handler {
	return &gRPCServer{
		chShutdown: make(chan struct{}, 1),
		pm:         pm,
		log:        log,
		wg:         wg,
	}
}

func (s *gRPCServer) Start() error {
	go func(s *gRPCServer) {
		for {
			select {
			case <-s.chShutdown:
				s.grpcServer.GracefulStop()
			}
		}
	}(s)
	return nil
}

func (s *gRPCServer) Stop() error {
	close(s.chShutdown)
	return nil
}
