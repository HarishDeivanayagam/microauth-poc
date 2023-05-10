package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	InvalidRequestBody   = "invalid request body"
	UserCreated          = "user created successfully"
	InvalidEmailPassword = "invalid email address or password"
	UnableSingup         = "unable to signup user"
)

type UserService interface {
	Login(context.Context, string, string) (string, string, error)
	Signup(context.Context, string, string, string, string) (string, error)
}

// login request
type LoginRequest struct {
	Email    string `query:"email"`
	Password string `query:"password"`
}

// login request
type SignupRequest struct {
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
	Email     string `query:"email"`
	Password  string `query:"password"`
}

// access token and refresh tokens to be returned
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Http) LoginHandler(ctx echo.Context) error {
	body := LoginRequest{}
	err := ctx.Bind(&body)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}
	accessToken, refreshToken, err := h.userService.Login(ctx.Request().Context(), body.Email, body.Password)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InvalidEmailPassword)
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, tokens)
}

func (h *Http) SignupHandler(ctx echo.Context) error {
	body := SignupRequest{}
	err := ctx.Bind(&body)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}
	_, err = h.userService.Signup(ctx.Request().Context(), body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, UnableSingup)
	}
	return ctx.JSON(http.StatusCreated, "account created")
}
