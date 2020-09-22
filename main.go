package main

import (
	"github.com/shoichiimamura/grpc-experiment/rpc"
	"github.com/shoichiimamura/grpc-experiment/server"
	"log"
)

func main() {
	params := rpc.Params{Addr: ":50051"}
	s1 := rpc.New(params)
	mux := server.NewMux([]server.Server{s1})
	if err := mux.Serve(); err != nil {
		log.Fatalln(err)
	}
}
