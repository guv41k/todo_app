package tasks_transport_http

import (
	"net/http"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_logger "github.com/guv41k/todo_app/internal/core/logger"
	core_http_request "github.com/guv41k/todo_app/internal/core/transport/http/request"
	core_http_response "github.com/guv41k/todo_app/internal/core/transport/http/response"
	core_http_types "github.com/guv41k/todo_app/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)
	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		response.ErrorResponse(err, "failed to get task id")
		return
	}
	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "failed to decode and validate request")
		return
	}
	patch := domain.NewTaskPatch(request.Title.ToDomain(), request.Description.ToDomain(), request.Completed.ToDomain())
	task, err := h.tasksService.PatchTask(ctx, id, patch)
	if err != nil {
		response.ErrorResponse(err, "failed to patch task")
		return
	}
	response.JSONResponse(taskDTOFromDomain(task), http.StatusOK)
}
