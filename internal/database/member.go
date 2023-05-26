package database

import (
	"context"
	"errors"
	"log"
	"time"

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

type MemberInviteRow struct {
	ID             string `db:"id"`
	Email          string `db:"email"`
	Code           string `db:"code"`
	OrganizationID string `db:"organization_id"`
	CreatedAt      int    `db:"created_at"`
	UpdatedAt      int    `db:"updated_at"`
	ExpiresAt      int    `db:"expires_at"`
}

var (
	FetchMemberInviteFailed  = errors.New("unable to get member invite")
	InsertMemberInviteFailed = errors.New("unable to insert member invite")
	DeleteMemberInviteFailed = errors.New("unable to delete member invite")
	FetchMemberFailed        = errors.New("unable to fetch member")
	MemberCreateFailed       = errors.New("unable to create member")
	MemberUpdateFailed       = errors.New("unable to update member")
	MemberDeleteFailed       = errors.New("unable to delete member")
	MemberUpdated            = "member updated"
	MemberDeleted            = "member deleted"
	MemberInviteDeleted      = "member invite deleted"
)

func (db *Database) GetMemberInvite(ctx context.Context, email string, organizationID string) (member.MemberInvite, error) {
	query := `
	SELECT id, email, code, organization_id, created_at, updated_at, expires_at
	FROM member_invite
	WHERE email = $1 AND organization_id = $2
	LIMIT 1
	`

	row := db.client.QueryRowxContext(ctx, query, email, organizationID)

	var invite MemberInviteRow
	err := row.StructScan(&invite)
	if err != nil {
		return member.MemberInvite{}, err
	}

	return member.MemberInvite{
		ID:             invite.ID,
		Email:          invite.Email,
		OrganizationID: invite.OrganizationID,
		Code:           invite.Code,
		CreatedAt:      invite.CreatedAt,
		UpdatedAt:      invite.UpdatedAt,
		ExpiresAt:      invite.ExpiresAt,
	}, nil
}

func (db *Database) InsertMemberInvite(ctx context.Context, email string, organizationID string, code string, expiresAt int) (string, error) {
	query := `
	INSERT INTO member_invite (id, email, code, organization_id, created_at, updated_at, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	id := uuid.New().String()
	createdAt := int(time.Now().Unix())
	updatedAt := createdAt

	_, err := db.client.ExecContext(ctx, query, id, email, code, organizationID, createdAt, updatedAt, expiresAt)
	if err != nil {
		log.Println(err)
		return "", InsertMemberInviteFailed
	}

	return id, nil
}

func (db *Database) DeleteMemberInvite(ctx context.Context, email string, organizationID string) (string, error) {
	query := `
	DELETE FROM member_invite
	WHERE email = $1 AND organization_id = $2
	`

	_, err := db.client.ExecContext(ctx, query, email, organizationID)
	if err != nil {
		return "", DeleteMemberInviteFailed
	}

	return MemberInviteDeleted, nil
}

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
		log.Println(err)
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
