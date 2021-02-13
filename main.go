package main

import (
	"log"

	"github.com/immrshc/grpc-experiment/rpc"
	"github.com/immrshc/grpc-experiment/server"
)

func main() {
	params := rpc.Params{Addr: ":50051"}
	s1 := rpc.New(params)
	mux := server.NewMux([]server.Server{s1})
	if err := mux.Serve(); err != nil {
		log.Fatalln(err)
	}
}
