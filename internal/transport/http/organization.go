package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"microauth.io/core/internal/member"
	"microauth.io/core/internal/organization"
)

type OrganizationService interface {
	FetchUserOrganizations(context.Context, string) ([]organization.Organization, error)
	GetOrganization(context.Context, string) (organization.Organization, error)
	CreateOrganization(context.Context, string, string) (string, error)
	DeleteOrganization(context.Context, string) (string, error)
	EditOrganization(context.Context, string, string, string) (string, error)
}

type CreateOrganizationRequest struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	AppRole string `json:"app_role"`
}

type OrganizationsResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

var (
	CreateOrganizationFailed   = "create organization failed"
	UnableToFetchOrganizations = "unable to fetch organizations"
)

func (h *Http) FetchOrganizationsHandler(ctx echo.Context) error {
	organizations, err := h.organizationService.FetchUserOrganizations(ctx.Request().Context(), ctx.Get("UserID").(string))
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, UnableToFetchOrganizations)
	}
	response := make([]OrganizationsResponse, len(organizations))

	for i, org := range organizations {
		response[i] = OrganizationsResponse{
			ID:        org.ID,
			Name:      org.Name,
			Domain:    org.Domain,
			CreatedAt: org.CreatedAt,
			UpdatedAt: org.UpdatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Http) CreateOrganizationHandler(ctx echo.Context) error {
	body := CreateOrganizationRequest{}
	err := ctx.Bind(&body)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}

	orgID, err := h.organizationService.CreateOrganization(ctx.Request().Context(), body.Name, body.Domain)

	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, CreateOrganizationFailed)
	}

	_, err = h.memberService.AddMember(ctx.Request().Context(), orgID, ctx.Get("UserID").(string), member.Admin, body.AppRole)

	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, CreateOrganizationFailed)
	}

	return ctx.String(http.StatusOK, "create organization success")
}
