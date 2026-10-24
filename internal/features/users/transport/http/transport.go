package users_transport_http

import (
	"context"
	"net/http"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_http_server "github.com/guv41k/todo_app/internal/core/transport/http/server"
)

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

type UsersHTTPHandler struct {
	usersService UsersService
}

func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
