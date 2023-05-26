package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"microauth.io/core/internal/organization"
)

type OrganizationRow struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Domain    string `db:"domain"`
	CreatedAt int    `db:"created_at"`
	UpdatedAt int    `db:"updated_at"`
}

var (
	FetchOrganizationFailed    = errors.New("unable to find organization")
	OrganizationCreationFailed = errors.New("unable to create organization")
	OrganizationCreated        = "organization created"
	OrganizationDeleteFailed   = errors.New("unable to delete organization")
	OrganizationDeleted        = "organization deleted"
	OrganizationUpdateFailed   = errors.New("unable to update organization")
	OrganizationUpdated        = "organization updated"
)

func (db *Database) GetOrganizationByUserID(ctx context.Context, userID string) ([]organization.Organization, error) {
	query := `
		SELECT id, name, domain, created_at, updated_at
		FROM organizations
		WHERE id IN (
			SELECT organization_id
			FROM members
			WHERE user_id = $1
		)
	`

	rows, err := db.client.QueryxContext(ctx, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	organizations := make([]organization.Organization, 0)

	for rows.Next() {
		var row OrganizationRow
		err := rows.StructScan(&row)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		organizations = append(organizations, organization.Organization{
			ID:        row.ID,
			Name:      row.Name,
			Domain:    row.Domain,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return organizations, nil
}

func (db *Database) InsertOrganization(ctx context.Context, name string, domain string) (string, error) {
	orgID := uuid.New().String()
	org := OrganizationRow{
		ID:        orgID,
		Name:      name,
		Domain:    domain,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	}
	query := `
		INSERT INTO organizations (id, name, domain, created_at, updated_at)
		VALUES (:id, :name, :domain, :created_at, :updated_at)
	`
	_, err := db.client.NamedExecContext(ctx, query, &org)
	if err != nil {
		return "", OrganizationCreationFailed
	}
	return orgID, nil
}

func (db *Database) GetOrganizationByID(ctx context.Context, id string) (organization.Organization, error) {
	org := OrganizationRow{}

	query := `
		SELECT * FROM organizations
		WHERE id = $1
	`

	err := db.client.GetContext(ctx, &org, query, id)

	if err != nil {
		return organization.Organization{}, err
	}

	return organization.Organization{
		ID:        org.ID,
		Name:      org.Name,
		Domain:    org.Domain,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}, nil
}

func (db *Database) DeleteOrganizationByID(ctx context.Context, id string) (string, error) {
	query := `
		DELETE FROM organizations
		WHERE id = $1
	`

	result, err := db.client.ExecContext(ctx, query, id)
	if err != nil {
		return "", OrganizationDeleteFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", OrganizationDeleteFailed
	}

	if rowsAffected == 0 {
		return "", FetchOrganizationFailed
	}

	return OrganizationDeleted, nil
}

func (db *Database) UpdateOrganization(ctx context.Context, id string, name string, domain string) (organization.Organization, error) {
	// Prepare the SQL query
	query := `
		UPDATE organizations
		SET name = $1, domain = $2, updated_at = $3
		WHERE id = $4
	`

	// Execute the SQL query
	_, err := db.client.ExecContext(ctx, query, name, domain, time.Now().Unix(), id)
	if err != nil {
		return organization.Organization{}, OrganizationUpdateFailed
	}

	// Retrieve the updated organization
	updatedOrg, err := db.GetOrganizationByID(ctx, id)
	if err != nil {
		return organization.Organization{}, err
	}

	return updatedOrg, nil
}
