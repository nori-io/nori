package grpc

import (
	"fmt"
	"net"

	"github.com/nori-io/nori-grpc/pkg/api/proto"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type Params struct {
	dig.In

	Env     *env.Env
	Handler *Handler
}

type Server struct {
	*grpc.Server
	listener *net.Listener
}

func New(params Params) (*Server, error) {
	if params.Env.Config.Nori.GRPC.Port == 0 {
		return nil, fmt.Errorf("gRPC port isn't defined")
	}

	addr := fmt.Sprintf("%s:%d", params.Env.Config.Nori.GRPC.Host, params.Env.Config.Nori.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	params.Env.Logger.Info("gRPC server started on " + addr)

	server := grpc.NewServer()
	proto.RegisterNoriServer(server, params.Handler)

	return &Server{
		Server:   server,
		listener: &listener,
	}, nil
}

func (s *Server) Start() error {
	if err := s.Serve(*s.listener); err != nil {
		if err != grpc.ErrServerStopped {
			return err
		}
	}
	return nil
}

func (s *Server) Shutdown() {
	s.GracefulStop()
}
