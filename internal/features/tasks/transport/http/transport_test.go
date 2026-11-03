package tasks_transport_http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_logger "github.com/guv41k/todo_app/internal/core/logger"
	core_http_middleware "github.com/guv41k/todo_app/internal/core/transport/http/middleware"
	core_http_server "github.com/guv41k/todo_app/internal/core/transport/http/server"
)

type taskServiceStub struct {
	created domain.Task
	userID  *int
}

func (s *taskServiceStub) CreateTask(_ context.Context, task domain.Task) (domain.Task, error) {
	s.created = task
	task.ID = 7
	task.Version = 1
	return task, nil
}
func (s *taskServiceStub) GetTasks(_ context.Context, userID, _, _ *int) ([]domain.Task, error) {
	s.userID = userID
	return []domain.Task{}, nil
}
func (s *taskServiceStub) GetTask(_ context.Context, _ int) (domain.Task, error) {
	return domain.Task{}, nil
}
func (s *taskServiceStub) PatchTask(_ context.Context, _ int, _ domain.TaskPatch) (domain.Task, error) {
	return domain.Task{}, nil
}
func (s *taskServiceStub) DeleteTask(_ context.Context, _ int) error { return nil }

func testRouter(t *testing.T, service *taskServiceStub) http.Handler {
	t.Helper()
	logger, err := core_logger.NewLogger(core_logger.Config{Level: "error", Folder: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(logger.Close)
	router := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	router.RegisterRoutes(NewTasksHTTPHandler(service).Routes()...)
	return core_http_middleware.Logger(logger)(router)
}

func TestCreateTaskRouteAndResponse(t *testing.T) {
	service := &taskServiceStub{}
	router := testRouter(t, service)
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Buy milk","author_user_id":3}`))
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rw.Code, rw.Body.String())
	}
	if service.created.Title != "Buy milk" || service.created.AuthorUserID != 3 {
		t.Fatalf("incorrect task passed to service: %+v", service.created)
	}
	var response TaskDTOResponse
	if err := json.Unmarshal(rw.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.ID != 7 || response.Version != 1 || response.Completed {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestGetTasksFilterAndInvalidCreate(t *testing.T) {
	service := &taskServiceStub{}
	router := testRouter(t, service)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, httptest.NewRequest(http.MethodGet, "/tasks?user_id=3", nil))
	if rw.Code != http.StatusOK || service.userID == nil || *service.userID != 3 || strings.TrimSpace(rw.Body.String()) != "[]" {
		t.Fatalf("filter response = %d %s, userID = %v", rw.Code, rw.Body.String(), service.userID)
	}
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"","author_user_id":3}`)))
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d", rw.Code)
	}
}
