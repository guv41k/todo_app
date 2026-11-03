package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/guv41k/todo_app/internal/core/logger"
	core_http_request "github.com/guv41k/todo_app/internal/core/transport/http/request"
	core_http_response "github.com/guv41k/todo_app/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), rw)
	userID, err := core_http_request.GetIntQueryParam(r, "user_id")
	if err != nil {
		response.ErrorResponse(err, "failed to get user_id")
		return
	}
	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		response.ErrorResponse(err, "failed to get limit")
		return
	}
	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		response.ErrorResponse(err, "failed to get offset")
		return
	}
	tasks, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		response.ErrorResponse(err, "failed to get tasks")
		return
	}
	response.JSONResponse(tasksDTOFromDomains(tasks), http.StatusOK)
}
