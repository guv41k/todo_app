package users_transport_http

import (
	"net/http"

	core_http_server "github.com/guv41k/todo_app/internal/core/transport/http/server"
)

type UsersService interface {
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
			Method: http.MethodPost,
			Path: "/users",
			Handler: h.CreateUser,
		},
	}
}