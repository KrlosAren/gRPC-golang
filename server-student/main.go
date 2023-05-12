package main

import (
	"krlosaren/go/grpc/database"
	"krlosaren/go/grpc/server"
	"krlosaren/go/grpc/studentpbf"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRespository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")
	server := server.NewStudentServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	studentpbf.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
