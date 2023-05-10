package http

import (
	"github.com/labstack/echo/v4"
)

type Http struct {
	server      *echo.Echo
	userService UserService
}

func New(userService UserService) *Http {
	return &Http{
		userService: userService,
		server:      echo.New(),
	}
}

// start the echo server
func (h *Http) Start(port string) {
	h.server.Logger.Fatal(h.server.Start(":" + port))
}

// register all the routes here
func (h *Http) RegisterHandlers() {
	h.server.POST("/api/v1/users/signup", h.SignupHandler)
	h.server.POST("/api/v1/users/login", h.LoginHandler)
}
