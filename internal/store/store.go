package store

import (
	"database/sql"
	"errors"
	"strconv"

	"todo/internal/types"
)

type Store interface {
	CreateTask(t *types.Task) (*types.Task, error)
	GetTasks() ([]*types.Task, error)
	GetTask(string) (*types.Task, error)
	PutTask(*types.Task) error
	DeleteTask(string) error
	UpdateTaskDate(nextDate string, id string) error
}

type Todo struct {
	Store
}

type Storage struct {
	db *sql.DB
}

func StorageSQLite(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func NewStore(db *sql.DB) *Todo {
	return &Todo{
		Store: StorageSQLite(db),
	}
}

func (s *Storage) CreateTask(t *types.Task) (*types.Task, error) {
	rows, err := s.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		t.Date, t.Title, t.Comment, t.Repeat)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = strconv.Itoa(int(id))
	return t, nil
}

func (s *Storage) GetTasks() ([]*types.Task, error) {
	const limit = 50
	rows, err := s.db.Query(`SELECT * FROM scheduler ORDER BY date ASC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		p, err := scanRowsIntoScheduler(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) GetTask(id string) (*types.Task, error) {
	row := s.db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
	var task types.Task
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task no found")
		} else {
			return nil, err
		}
	}

	return &task, nil
}

func (s *Storage) PutTask(t *types.Task) error {
	rows, err := s.db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return err
}

func (s *Storage) UpdateTaskDate(nextDate string, id string) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
	return err
}

func scanRowsIntoScheduler(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)

	err := rows.Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}
