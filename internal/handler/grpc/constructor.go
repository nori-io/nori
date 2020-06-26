package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	listener *net.Listener
}

func New() (*Server, error) {
	// todo: load addr:port from config
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	// todo: register grpc serber

	return &Server{
		Server:   server,
		listener: &listener,
	}, nil
}

func (s *Server) Start() {
	//
}

func (s *Server) Stop() {
	//
}
