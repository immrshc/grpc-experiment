package helloworld

import (
	"context"
	"log"

	"github.com/immrshc/grpc-experiment/proto"
	"google.golang.org/grpc"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (h *server) Register(gs *grpc.Server) {
	proto.RegisterGreeterServer(gs, h)
}

func (h *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloReply{Message: "hello" + in.GetName()}, nil
}
