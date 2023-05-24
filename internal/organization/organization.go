package organization

import (
	"context"
	"errors"
	"log"
)

type Organization struct {
	ID        string
	Name      string
	Domain    string
	CreatedAt int
	UpdatedAt int
}

var (
	FetchOrganizationFailed    = errors.New("unable to find organization")
	OrganizationCreationFailed = errors.New("unable to create organization")
	OrganizationDeleteFailed   = errors.New("unable to delete organization")
	OrganizationDeleted        = "organization deleted"
	OrganizationUpdateFailed   = errors.New("unable to update organization")
	OrganizationUpdated        = "organization updated"
)

type OrganizationStore interface {
	GetOrganizationByUserID(context.Context, string) ([]Organization, error)
	InsertOrganization(context.Context, string, string) (string, error)
	GetOrganizationByID(context.Context, string) (Organization, error)
	DeleteOrganizationByID(context.Context, string) (string, error)
	UpdateOrganization(context.Context, string, string, string) (Organization, error)
}

type Service struct {
	store OrganizationStore
}

func New(store OrganizationStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) FetchUserOrganizations(ctx context.Context, userID string) ([]Organization, error) {
	organizations, err := s.store.GetOrganizationByUserID(ctx, userID)
	if err != nil {
		log.Println(err)
		return []Organization{}, FetchOrganizationFailed
	}
	return organizations, nil
}

func (s *Service) GetOrganization(ctx context.Context, id string) (Organization, error) {
	organization, err := s.store.GetOrganizationByID(ctx, id)
	if err != nil {
		return Organization{}, FetchOrganizationFailed
	}
	return organization, nil
}

func (s *Service) CreateOrganization(ctx context.Context, name string, domain string) (string, error) {
	orgID, err := s.store.InsertOrganization(ctx, name, domain)
	if err != nil {
		return "", OrganizationCreationFailed
	}
	return orgID, nil
}

func (s *Service) DeleteOrganization(ctx context.Context, id string) (string, error) {
	_, err := s.store.DeleteOrganizationByID(ctx, id)
	if err != nil {
		return "", OrganizationDeleteFailed
	}
	return OrganizationDeleted, nil
}

func (s *Service) EditOrganization(ctx context.Context, id string, name string, domain string) (string, error) {
	_, err := s.store.UpdateOrganization(ctx, id, name, domain)
	if err != nil {
		return "", OrganizationUpdateFailed
	}
	return OrganizationUpdated, nil
}
