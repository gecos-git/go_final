package service

import (
	"todo/internal/store"
	"todo/internal/types"
)

type ServiceStore interface {
	CreateTask(t *types.Task) (*types.Task, error)
	GetTasks() ([]*types.Task, error)
	GetTask(string) (*types.Task, error)
	PutTask(*types.Task) error
	DoneTask(string) error
	DeleteTask(string) error
	UpdateTaskDate(nextDate string, id string) error
}

type Service struct {
	ServiceStore
}

func NewService(repos *store.Todo) *Service {
	return &Service{
		ServiceStore: NewTaskService(repos.Store),
	}
}
