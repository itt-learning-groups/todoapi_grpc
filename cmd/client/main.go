package main

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	todoapiv1 "github.com/itt-learning-groups/proto-contracts/todoapi/gen/go/v1"
	"github.com/itt-learning-groups/todoapi-grpc/internal/interceptors"
	"google.golang.org/grpc"
	"log"
)

func main() {
	serverPort := 50051
	serverAddress := fmt.Sprintf("localhost:%d", serverPort)
	
	//conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptors.PrintHeaders),
	)

	defer conn.Close()
	if err != nil {
		log.Fatalf("failed to connect to server at %v: %v", serverAddress, err)
	}
	
	client := todoapiv1.NewTodoapiServiceClient(conn)
	
	req := todoapiv1.CreateTodoRequest{
		Name:        faker.Word(),
		Description: faker.Sentence(),
	}
	
	res, err := client.CreateTodo(context.Background(), &req)
	if err != nil {
		log.Printf("failed to create new Todo: %v", err)
	}

	log.Printf("CreateTodo response: %+v", res)

	log.Println("client exiting")
}
