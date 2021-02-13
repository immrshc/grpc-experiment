package main

import (
	"context"
	"log"
	"time"

	"github.com/immrshc/grpc-experiment/proto"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "hello"})
	if err != nil {
		log.Fatalf("cound not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
