package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	ports = []string{":10001", ":10002", ":10003"}
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

type slowServer struct {
	pb.UnimplementedGreeterServer
}

func (s *slowServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// Delay 100ms ~ 200ms before replying
	time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// Set up the server channelz service.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	defer s.Stop()

	// Start three GreeterServers(with one of them to be the slowServer).
	var listeners []net.Listener
	var svrs []*grpc.Server
	for i := 0; i < 3; i++ {
		lis, err := net.Listen("tcp", ports[i])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		listeners = append(listeners, lis)
		s := grpc.NewServer()
		svrs = append(svrs, s)
		if i == 2 {
			pb.RegisterGreeterServer(s, &slowServer{})
		} else {
			pb.RegisterGreeterServer(s, &server{})
		}
		go s.Serve(lis)
	}

	// Wait for CTRL+C to exit.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received.
	<-ch
	for i := 0; i < 3; i++ {
		svrs[i].Stop()
		listeners[i].Close()
	}
}
