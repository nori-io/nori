package grpc

import (
	"sync"

	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/plugins"
	"github.com/nori-io/nori/internal/server"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	chShutdown chan struct{}
	grpcServer *grpc.Server
	pm         *plugins.Manager
	secure     bool
	log        logger.Logger
	wg         *sync.WaitGroup
}

func NewServer(pm *plugins.Manager, log logger.Logger, wg *sync.WaitGroup) server.Server {
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
