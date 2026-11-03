package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_http_server "github.com/guv41k/todo_app/internal/core/transport/http/server"
)

type TasksService interface {
	CreateTask(context.Context, domain.Task) (domain.Task, error)
	GetTasks(context.Context, *int, *int, *int) ([]domain.Task, error)
	GetTask(context.Context, int) (domain.Task, error)
	PatchTask(context.Context, int, domain.TaskPatch) (domain.Task, error)
	DeleteTask(context.Context, int) error
}

type TasksHTTPHandler struct {
	tasksService TasksService
}

func NewTasksHTTPHandler(service TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: service}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/tasks", Handler: h.CreateTask},
		{Method: http.MethodGet, Path: "/tasks", Handler: h.GetTasks},
		{Method: http.MethodGet, Path: "/tasks/{id}", Handler: h.GetTask},
		{Method: http.MethodPatch, Path: "/tasks/{id}", Handler: h.PatchTask},
		{Method: http.MethodDelete, Path: "/tasks/{id}", Handler: h.DeleteTask},
	}
}
