package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"todo/internal/nextdate"
	service "todo/internal/service/tasks"
	"todo/internal/types"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTask(w, r)
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPut:
		h.EditTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		log.Println("неверный метод")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "неверный метод"}); err != nil {
			log.Println(err)
		}
		return
	}
}

// Создание задачи
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	t, err := h.service.CreateTask(task)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

// Получение списка задач
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

// Получение задачи
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := h.service.GetTask(taskID)
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
func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.PutTask(task); err != nil {
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{})
}

// Отметить задачу выполненой
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}
	if err := h.service.DoneTask(taskID); err != nil {
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
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
		return
	}

	if _, err := strconv.Atoi(taskID); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteTask(taskID); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{})
}

// Да
func (h *Handler) NextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeatStr)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Println("Ошибка при записи даты:", err)
	}
}
