package http

import (
	"github.com/labstack/echo/v4"
)

type Http struct {
	server              *echo.Echo
	userService         UserService
	organizationService OrganizationService
	memberService       MemberService
}

var (
	InvalidRequestBody  = "invalid request body"
	InternalServerError = "some error happened"
)

func New(userService UserService, organizationService OrganizationService, memberService MemberService) *Http {
	return &Http{
		userService:         userService,
		organizationService: organizationService,
		memberService:       memberService,
		server:              echo.New(),
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

	// authenticated requests
	authenticated := h.server.Group("/api/v1")
	authenticated.Use(h.JWTMiddleware)
	authenticated.GET("/organizations", h.FetchOrganizationsHandler)
	authenticated.POST("/organizations", h.CreateOrganizationHandler)
	authenticated.GET("/organizations/:organizationID/members/me", h.FetchAllMembersHandler)
	authenticated.GET("/organizations/:organizationID/members", h.FetchMemberHandler)
}
