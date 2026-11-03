package tasks_transport_http

import (
	"net/http"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_logger "github.com/guv41k/todo_app/internal/core/logger"
	core_http_request "github.com/guv41k/todo_app/internal/core/transport/http/request"
	core_http_response "github.com/guv41k/todo_app/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"gt=0"`
}

func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)
	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "failed to decode and validate request")
		return
	}
	task := domain.NewTaskUninitialized(request.Title, request.Description, request.AuthorUserID)
	created, err := h.tasksService.CreateTask(ctx, task)
	if err != nil {
		response.ErrorResponse(err, "failed to create task")
		return
	}
	response.JSONResponse(taskDTOFromDomain(created), http.StatusCreated)
}
