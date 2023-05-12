package main

import (
	"krlosaren/go/grpc/database"
	"krlosaren/go/grpc/server"
	"krlosaren/go/grpc/testpbf"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRespository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	server := server.NewTestServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	testpbf.RegisterTestServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
