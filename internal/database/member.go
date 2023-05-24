package database

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"microauth.io/core/internal/member"
)

type MemberRow struct {
	ID             string      `db:"id"`
	OrganizationID string      `db:"organization_id"`
	UserID         string      `db:"user_id"`
	Role           member.Role `db:"role"`
	AppRole        string      `db:"app_role"`
}

var (
	FetchMemberFailed  = errors.New("unable to fetch member")
	MemberCreateFailed = errors.New("unable to create member")
	MemberUpdateFailed = errors.New("unable to update member")
	MemberUpdated      = "member updated"
	MemberDeleteFailed = errors.New("unable to delete member")
	MemberDeleted      = "member deleted"
)

func (db *Database) FetchAllMembers(ctx context.Context, organizationID string) ([]member.Member, error) {
	query := `
	SELECT id, organization_id, user_id, role, app_role
	FROM members
	WHERE organization_id = $1
	`

	rows, err := db.client.QueryxContext(ctx, query, organizationID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	members := make([]member.Member, 0)

	for rows.Next() {
		var mem MemberRow
		err := rows.StructScan(&mem)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		members = append(members, member.Member{
			ID:             mem.ID,
			OrganizationID: mem.OrganizationID,
			UserID:         mem.UserID,
			Role:           mem.Role,
			AppRole:        mem.AppRole,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return members, nil
}

func (db *Database) FetchMemberByID(ctx context.Context, organizationID string, userID string) (member.Member, error) {
	query := `
		SELECT * FROM members
		WHERE organization_id = $1 AND user_id = $2
	`

	var memberRow MemberRow

	err := db.client.GetContext(ctx, &memberRow, query, organizationID, userID)
	if err != nil {
		return member.Member{}, err
	}

	// Convert MemberRow to Member struct
	member := member.Member{
		ID:             memberRow.ID,
		OrganizationID: memberRow.OrganizationID,
		UserID:         memberRow.UserID,
		Role:           memberRow.Role,
		AppRole:        memberRow.AppRole,
	}

	return member, nil
}

func (db *Database) InsertMember(ctx context.Context, organizationID string, userID string, role member.Role, appRole string) (string, error) {
	memberID := uuid.New().String()

	// Create a new MemberRow instance with the provided data
	member := MemberRow{
		ID:             memberID,
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           role,
		AppRole:        appRole,
	}

	// Prepare the SQL query
	query := `
		INSERT INTO members (id, organization_id, user_id, role, app_role)
		VALUES (:id, :organization_id, :user_id, :role, :app_role)
	`

	// Execute the SQL query using named parameters
	_, err := db.client.NamedExecContext(ctx, query, &member)
	if err != nil {
		return "", MemberCreateFailed
	}

	return memberID, nil

}

func (db *Database) UpdateMember(ctx context.Context, organizationID string, userID string, role member.Role, appRole string) (string, error) {
	// Prepare the SQL query
	query := `
		UPDATE members
		SET role = $1, app_role = $2
		WHERE organization_id = $3 AND user_id = $4
	`

	// Execute the SQL query
	_, err := db.client.ExecContext(ctx, query, role, appRole, organizationID, userID)
	if err != nil {
		return "", MemberUpdateFailed
	}

	return MemberUpdated, nil
}

func (db *Database) DeleteMember(ctx context.Context, organizationID string, userID string) (string, error) {
	// Prepare the SQL query
	query := `
		DELETE FROM members
		WHERE organization_id = $1 AND user_id = $2
	`

	// Execute the SQL query
	result, err := db.client.ExecContext(ctx, query, organizationID, userID)
	if err != nil {
		return "", MemberDeleteFailed
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", MemberDeleteFailed
	}

	if rowsAffected == 0 {
		return "", MemberDeleteFailed
	}

	return MemberDeleted, nil

}
