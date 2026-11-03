package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/guv41k/todo_app/internal/core/logger"
	core_http_request "github.com/guv41k/todo_app/internal/core/transport/http/request"
	core_http_response "github.com/guv41k/todo_app/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)
	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		response.ErrorResponse(err, "failed to get task id")
		return
	}
	task, err := h.tasksService.GetTask(ctx, id)
	if err != nil {
		response.ErrorResponse(err, "failed to get task")
		return
	}
	response.JSONResponse(taskDTOFromDomain(task), http.StatusOK)
}
