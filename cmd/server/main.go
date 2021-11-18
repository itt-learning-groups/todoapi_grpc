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
	grpcServerPort := 50051
	grpcServerAddress := fmt.Sprintf("localhost:%d", grpcServerPort)

	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatalf("cannot listen to address %s: %v", grpcServerAddress, err)
	}

	grpcServer := grpc.NewServer()

	todoService := &services.TodoapiServer{}

	todoapiv1.RegisterTodoapiServiceServer(grpcServer, todoService)

	reflection.Register(grpcServer)

	log.Printf("serving TodoapiServer (grpc) on %s", grpcServerAddress)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
