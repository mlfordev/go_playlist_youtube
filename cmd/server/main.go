package main

import (
	"grpc_test/pkg/api"
	"grpc_test/pkg/playlist"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	srv := &playlist.GRPCServer{}
	api.RegisterPlaylistServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
