package services

import (
	"context"
	"github.com/google/uuid"
	todoapiv1 "github.com/itt-learning-groups/proto-contracts/todoapi/gen/go/v1"
)

type TodoapiServer struct{}

func (tas *TodoapiServer) CreateTodo(ctx context.Context, req *todoapiv1.CreateTodoRequest) (*todoapiv1.CreateTodoResponse, error) {
	id := uuid.New().String()
	return &todoapiv1.CreateTodoResponse{
		Todo: &todoapiv1.Todo{
			Id:          id,
			Name:        req.Name,
			Description: req.Description,
		},
	}, nil
}
