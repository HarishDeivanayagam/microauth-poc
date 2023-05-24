package member

import (
	"context"
	"errors"
)

type Role string

const (
	Admin Role = "admin"
	Staff Role = "staff"
	User  Role = "user"
)

var (
	AdminPermissionFailed = errors.New("you aren't a administrator")
	FetchMemberFailed     = errors.New("unable to fetch member")
	MemberCreateFailed    = errors.New("unable to create member")
	MemberAdded           = "member added"
	MemberUpdateFailed    = errors.New("unable to update member")
	MemberUpdated         = "member updated"
	MemberDeleteFailed    = errors.New("unable to delete member")
	MemberDeleted         = "member deleted"
)

type Member struct {
	// member id must not be used only for reference
	ID string

	OrganizationID string
	UserID         string
	Role           Role
	AppRole        string
}

type MemberStore interface {
	FetchMemberByID(context.Context, string, string) (Member, error)
	FetchAllMembers(context.Context, string) ([]Member, error)
	InsertMember(context.Context, string, string, Role, string) (string, error)
	UpdateMember(context.Context, string, string, Role, string) (string, error)
	DeleteMember(context.Context, string, string) (string, error)
}

type Service struct {
	store MemberStore
}

func New(store MemberStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) FetchAllMembers(ctx context.Context, organizationID string, userID string) ([]Member, error) {
	// check access
	member, err := s.store.FetchMemberByID(ctx, organizationID, userID)
	if err != nil {
		return []Member{}, FetchMemberFailed
	}
	if member.Role != Admin {
		return []Member{}, AdminPermissionFailed
	}

	members, err := s.store.FetchAllMembers(ctx, organizationID)
	if err != nil {
		return []Member{}, FetchMemberFailed
	}
	return members, nil
}

func (s *Service) FetchMember(ctx context.Context, organizationID string, userID string) (Member, error) {
	member, err := s.store.FetchMemberByID(ctx, organizationID, userID)
	if err != nil {
		return Member{}, FetchMemberFailed
	}
	return member, nil
}

func (s *Service) AddMember(ctx context.Context, organizationID string, userID string, role Role, appRole string) (string, error) {
	memberID, err := s.store.InsertMember(ctx, organizationID, userID, role, appRole)
	if err != nil {
		return "", MemberCreateFailed
	}
	return memberID, nil
}

func (s *Service) UpdateMember(ctx context.Context, organizationID string, userID string, role Role, appRole string) (string, error) {
	_, err := s.store.UpdateMember(ctx, organizationID, userID, role, appRole)
	if err != nil {
		return "", MemberUpdateFailed
	}
	return MemberUpdated, nil
}

func (s *Service) DeleteMember(ctx context.Context, organizationID string, userID string) (string, error) {
	_, err := s.store.DeleteMember(ctx, organizationID, userID)
	if err != nil {
		return "", MemberDeleteFailed
	}
	return MemberDeleted, nil
}
