package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"todo/storedb"
	"todo/types"
	"todo/utils"
)

var webDir = "./web/"

type APIServer struct {
	addr  string
	store storedb.Store
}

func NewAPIServer(addr string, store storedb.Store) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	// Домашняя страница
	router.Handle("/", http.FileServer(http.Dir(webDir)))

	// Создание задачи
	router.HandleFunc("POST /api/task", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		defer r.Body.Close()

		var task *types.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			return
		}

		if err := validateTaskPayload(task); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		t, err := s.store.CreateTask(task)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, t)
	})

	// Получение списка задач
	router.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := s.store.GetTasks()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
	})

	// Получение задачи
	router.HandleFunc("GET /api/task", func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		task, err := s.store.GetTask(taskID)
		if err != nil {
			if err.Error() == "task not found" {
				utils.WriteError(w, http.StatusInternalServerError, err)
			} else {
				utils.WriteError(w, http.StatusInternalServerError, err)
			}
			return
		}

		utils.WriteJSON(w, http.StatusOK, task)
	})

	// Редактирование задачи
	router.HandleFunc("PUT /api/task", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		defer r.Body.Close()

		var task *types.Task
		err = json.Unmarshal(body, &task)
		if err != nil {
			return
		}

		if err := validateTaskUpdate(task); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.PutTask(task); err != nil {
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
			}
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{})
	})

	// Отметить выполненой задачу
	router.HandleFunc("POST /api/task/done", func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
			return
		}
		if err := s.store.DoneTask(taskID); err != nil {
			if err.Error() == "task not found" {
				utils.WriteError(w, http.StatusInternalServerError, err)
			} else {
				utils.WriteError(w, http.StatusInternalServerError, err)
			}
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{})
	})

	// Удалить задачу
	router.HandleFunc("DELETE /api/task", func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		if _, err := strconv.Atoi(taskID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.DeleteTask(taskID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{})
	})

	router.HandleFunc("GET /api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		nowStr := r.FormValue("now")
		dateStr := r.FormValue("date")
		repeatStr := r.FormValue("repeat")

		now, err := time.Parse("20060102", nowStr)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		nextDate, err := utils.NextDate(now, dateStr, repeatStr)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if _, err := w.Write([]byte(nextDate)); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
	})

	server := http.Server{
		Addr:    s.addr,
		Handler: RequestLoggerMiddleware(router),
	}

	log.Printf("Система Т800 запущена...")
	log.Printf("Сервер слушает порт %s", s.addr)

	return server.ListenAndServe()
}
