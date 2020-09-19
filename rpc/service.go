package rpc

import (
	"github.com/shoichiimamura/grpc-experiment/rpc/helloworld"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	addr       string
	grpcServer *grpc.Server
	errCh      chan error
}

type Params struct {
	Addr string
}

func New(params Params) *server {
	gs := grpc.NewServer()
	for _, svc := range newServices() {
		svc.Register(gs)
	}
	return &server{
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

func (s *server) Start() error {
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

func (s *server) AsyncErr() <-chan error {
	return s.errCh
}
