package api

import (
	"log"
	"net/http"

	"todo/internal/service/tasks"
)

var webDir = "./web/"

type APIServer struct {
	addr  string
	store tasks.Store
}

func NewAPIServer(addr string, store tasks.Store) *APIServer {
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
	router.HandleFunc("POST /api/task", s.CreateTask)

	// Получение списка задач
	router.HandleFunc("GET /api/tasks", s.ListTasks)

	// Получение задачи
	router.HandleFunc("GET /api/task", s.GetTask)

	// Редактирование задачи
	router.HandleFunc("PUT /api/task", s.EditTask)

	// Отметить выполненой задачу
	router.HandleFunc("POST /api/task/done", s.DoneTask)

	// Удалить задачу
	router.HandleFunc("DELETE /api/task", s.DeleteTask)

	// Да
	router.HandleFunc("GET /api/nextdate", s.NextDate)

	server := http.Server{
		Addr:    s.addr,
		Handler: RequestLoggerMiddleware(router),
	}

	log.Printf("Система Т800 запущена...")
	log.Printf("Сервер слушает порт %s", s.addr)

	return server.ListenAndServe()
}
