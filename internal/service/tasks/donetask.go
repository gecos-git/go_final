package service

import (
	"database/sql"
	"time"

	"todo/internal/nextdate"
	"todo/internal/store"
	"todo/internal/types"
)

type ItemService struct {
	stor store.Store
}

func NewTaskService(stor store.Store) *ItemService {
	return &ItemService{stor: stor}
}

func (s *ItemService) CreateTask(t *types.Task) (*types.Task, error) {
	return s.stor.CreateTask(t)
}

func (s *ItemService) DeleteTask(id string) error {
	return s.stor.DeleteTask(id)
}

func (s *ItemService) GetTasks() ([]*types.Task, error) {
	return s.stor.GetTasks()
}

func (s *ItemService) GetTask(id string) (*types.Task, error) {
	return s.stor.GetTask(id)
}

func (s *ItemService) PutTask(t *types.Task) error {
	return s.stor.PutTask(t)
}

func (s *ItemService) UpdateTaskDate(nextDate string, id string) error {
	return s.stor.UpdateTaskDate(nextDate, id)
}

func (s *ItemService) DoneTask(id string) error {
	task, err := s.stor.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	if task.Repeat == "" {
		if err := s.stor.DeleteTask(id); err != nil {
			return err
		}
	} else {
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if err := s.stor.UpdateTaskDate(nextDate, id); err != nil {
			return err
		}
	}

	return nil
}
