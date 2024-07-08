package api

import (
	"log"
	"net/http"

	"todo/internal/handlers"
)

var webDir = "./web/"

type APIServer struct {
	addr     string
	handlers handlers.Handler
}

func NewAPIServer(addr string, handlers handlers.Handler) *APIServer {
	return &APIServer{
		addr:     addr,
		handlers: handlers,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	// Домашняя страница
	router.Handle("/", http.FileServer(http.Dir(webDir)))

	// Создание задачи
	router.HandleFunc("POST /api/task", s.handlers.CreateTask)

	// Получение списка задач
	router.HandleFunc("GET /api/tasks", s.handlers.ListTasks)

	// Получение задачи
	router.HandleFunc("GET /api/task", s.handlers.GetTask)

	// Редактирование задачи
	router.HandleFunc("PUT /api/task", s.handlers.EditTask)

	// Отметить выполненой задачу
	router.HandleFunc("POST /api/task/done", s.handlers.DoneTask)

	// Удалить задачу
	router.HandleFunc("DELETE /api/task", s.handlers.DeleteTask)

	// Да
	router.HandleFunc("GET /api/nextdate", s.handlers.NextDate)

	server := http.Server{
		Addr:    s.addr,
		Handler: RequestLoggerMiddleware(router),
	}

	log.Printf("Система Т800 запущена...")
	log.Printf("Сервер слушает порт %s", s.addr)

	return server.ListenAndServe()
}
