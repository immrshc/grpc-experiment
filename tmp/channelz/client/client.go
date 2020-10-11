package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

func main() {
	// channelzのRPC用のサーバーを起動する
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	// 確認用にgRPC Serviceの情報を返せるようにする
	// c.f. https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
	reflection.Register(s)
	go s.Serve(lis)
	defer s.Stop()

	// 三つのサーバーにラウンドロビンするための名前解決の設定
	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()
	state := resolver.State{Addresses: []resolver.Address{{Addr: ":10001"}, {Addr: ":10002"}, {Addr: ":10003"}}}
	r.InitialState(state)
	// サーバーへのコネクションを設定する
	conn, err := grpc.Dial(
		r.Scheme()+":///test.server",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// サーバーへRPCするクライアントの設定
	c := pb.NewGreeterClient(conn)
	// 100回RPCし、150msをタイムアウトの閾値とする
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
		if err != nil {
			log.Printf("could not greet: %v", err)
		} else {
			log.Printf("Greeting: %s", r.Message)
		}
	}

	// CTRL+Cでexitするまで待つことで、channelzの情報を保持しておける
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
