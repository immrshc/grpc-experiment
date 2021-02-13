package rpc

import (
	"net"

	"github.com/immrshc/grpc-experiment/rpc/helloworld"
	"google.golang.org/grpc"
)

// Server is a server for gRPC
type Server struct {
	addr       string
	grpcServer *grpc.Server
	errCh      chan error
}

// Params is a parameters to set up Server.
type Params struct {
	Addr string
}

// New generates Server.
func New(params Params) *Server {
	gs := grpc.NewServer()
	for _, svc := range newServices() {
		svc.Register(gs)
	}
	return &Server{
		addr:       params.Addr,
		grpcServer: gs,
	}
}

type serviceImpl interface {
	Register(*grpc.Server)
}

func newServices() []serviceImpl {
	return []serviceImpl{
		helloworld.NewServer(),
	}
}

// Start starts to serve gRPC call.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.errCh = make(chan error, 1)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.errCh <- err
		}
		close(s.errCh)
	}()
	return nil
}

// AsyncErr returns the channels containing error.
func (s *Server) AsyncErr() <-chan error {
	return s.errCh
}
