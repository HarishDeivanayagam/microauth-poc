package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"microauth.io/core/internal/member"
)

type MemberService interface {
	FetchAllMembers(context.Context, string, string) ([]member.Member, error)
	FetchMember(context.Context, string, string) (member.Member, error)
	AddMember(context.Context, string, string, member.Role, string) (string, error)
	UpdateMember(context.Context, string, string, member.Role, string) (string, error)
	DeleteMember(context.Context, string, string) (string, error)
}

type MemberResponse struct {
	ID             string      `json:"id"`
	OrganizationID string      `json:"organization_id"`
	UserID         string      `json:"user_id"`
	Role           member.Role `json:"role"`
	AppRole        string      `json:"app_role"`
}

func (h *Http) FetchAllMembersHandler(ctx echo.Context) error {
	members, err := h.memberService.FetchAllMembers(ctx.Request().Context(), ctx.Param("organizationID"), ctx.Get("UserID").(string))
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InternalServerError)
	}
	result := make([]MemberResponse, len(members))
	for i, mem := range members {
		result[i] = MemberResponse{
			ID:             mem.ID,
			OrganizationID: mem.OrganizationID,
			UserID:         mem.UserID,
			Role:           mem.Role,
			AppRole:        mem.AppRole,
		}
	}

	return ctx.JSON(http.StatusOK, result)

}

func (h *Http) FetchMemberHandler(ctx echo.Context) error {
	member, err := h.memberService.FetchMember(ctx.Request().Context(), ctx.Param("organizationID"), ctx.Get("UserID").(string))
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusBadRequest, InternalServerError)
	}
	result := MemberResponse{
		ID:             member.ID,
		OrganizationID: member.OrganizationID,
		UserID:         member.UserID,
		Role:           member.Role,
		AppRole:        member.AppRole,
	}
	return ctx.JSON(http.StatusOK, result)
}
