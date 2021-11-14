package main

import (
	"fmt"
	todoapiv1 "github.com/itt-learning-groups/proto-contracts/todoapi/gen/go/v1"
	"github.com/itt-learning-groups/todoapi-grpc/cmd/server/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	serverPort := 50051
	serverAddress := fmt.Sprintf("localhost:%d", serverPort)

	conn, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", serverAddress)
	}

	grpcServer := grpc.NewServer()

	todoService := &services.TodoapiServer{}

	todoapiv1.RegisterTodoapiServiceServer(grpcServer, todoService)

	reflection.Register(grpcServer)

	log.Printf("serving TodoapiServer on %s", serverAddress)
	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
