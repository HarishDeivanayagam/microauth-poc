package http

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"microauth.io/core/internal/member"
)

type MemberService interface {
	InviteMember(context.Context, string, string, string) (string, error)
	FetchAllMembers(context.Context, string, string) ([]member.Member, error)
	FetchMember(context.Context, string, string) (member.Member, error)
	AddMember(context.Context, string, string, member.Role, string) (string, error)
	UpdateMember(context.Context, string, string, member.Role, string) (string, error)
	DeleteMember(context.Context, string, string) (string, error)
	AcceptInvite(ctx context.Context, email string, code string, organizationID string, firstName string, lastName string, password string) (string, error)
}

type MemberResponse struct {
	ID             string      `json:"id"`
	OrganizationID string      `json:"organization_id"`
	UserID         string      `json:"user_id"`
	Role           member.Role `json:"role"`
	AppRole        string      `json:"app_role"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
}

type AcceptInviteRequest struct {
	Email          string `json:"email"`
	Code           string `json:"code"`
	OrganizationID string `json:"organization_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
}

func (h *Http) AcceptInviteHandler(ctx echo.Context) error {
	// Parse the request body
	var request AcceptInviteRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}

	// Invoke the service to accept the invitation
	response, err := h.memberService.AcceptInvite(
		ctx.Request().Context(),
		request.Email,
		request.Code,
		request.OrganizationID,
		request.FirstName,
		request.LastName,
		request.Password,
	)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, InternalServerError)
	}

	return ctx.String(http.StatusOK, response)
}

func (h *Http) InviteMemberHandler(ctx echo.Context) error {
	// Get the organization ID from path parameter
	organizationID := ctx.Param("organizationID")

	// Get the user ID from path parameter
	userID := ctx.Get("UserID").(string)

	// Parse the request body
	var request InviteMemberRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.String(http.StatusBadRequest, InvalidRequestBody)
	}

	// Invoke the service to invite a member
	result, err := h.memberService.InviteMember(ctx.Request().Context(), request.Email, userID, organizationID)
	if err != nil {
		log.Println(err)
		return ctx.String(http.StatusInternalServerError, InternalServerError)
	}

	return ctx.String(http.StatusOK, result)
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
