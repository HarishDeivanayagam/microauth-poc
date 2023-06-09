package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	UserCreated          = "user created successfully"
	InvalidEmailPassword = "invalid email address or password"
	UnableSingup         = "unable to signup user"
)

type UserService interface {
	GenerateAccessToken(context.Context, string) (string, string, error)
	Login(context.Context, string, string) (string, string, error)
	CreateUser(context.Context, string, string, string, string) (string, error)
}

// login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// login request
type SignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// access token and refresh tokens to be returned
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func (h *Http) RefreshTokenHandler(ctx echo.Context) error {
	body := RefreshTokenRequest{}
	err := ctx.Bind(&body)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}

	// Generate a new access token and refresh token
	newAccessToken, newRefreshToken, err := h.userService.GenerateAccessToken(ctx.Request().Context(), body.AccessToken)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, InternalServerError)
	}

	tokens := Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	return ctx.JSON(http.StatusOK, tokens)
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
	_, err = h.userService.CreateUser(ctx.Request().Context(), body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, UnableSingup)
	}
	return ctx.JSON(http.StatusCreated, "account created")
}
