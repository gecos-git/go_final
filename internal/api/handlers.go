package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"todo/internal/service"
	"todo/internal/types"
)

// Создание задачи
func (s *APIServer) CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	var task *types.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

// Получение списка задач
func (s *APIServer) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.GetTasks()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

// Получение задачи
func (s *APIServer) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := s.store.GetTask(taskID)
	if err != nil {
		if err.Error() == "task not found" {
			WriteError(w, http.StatusNotFound, err)
		} else {
			WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

// Редактирование задачи
func (s *APIServer) EditTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	defer r.Body.Close()

	var task *types.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateTaskUpdate(task); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := s.store.PutTask(task); err != nil {
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{})
}

// Отметить задачу выполненой
func (s *APIServer) DoneTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}
	if err := s.store.DoneTask(taskID); err != nil {
		if err.Error() == "task not found" {
			WriteError(w, http.StatusNotFound, err)
		} else {
			WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{})
}

// Удалить задачу
func (s *APIServer) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	if _, err := strconv.Atoi(taskID); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.store.DeleteTask(taskID); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{})
}

// Да
func (s *APIServer) NextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	nextDate, err := service.NextDate(now, dateStr, repeatStr)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Println("Ошибка при записи даты:", err)
	}
}
