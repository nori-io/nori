package rest

import (
	"sync"

	"github.com/nori-io/nori-common/logger"
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
