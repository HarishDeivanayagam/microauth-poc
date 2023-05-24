package database

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"microauth.io/core/internal/user"
)

type UserRow struct {
	ID              string `db:"id"`
	FirstName       string `db:"first_name"`
	LastName        string `db:"last_name"`
	Email           string `db:"email"`
	IsEmailVerified bool   `db:"is_email_verified"`
	Password        string `db:"password"`
	CreatedAt       int    `db:"created_at"`
	UpdatedAt       int    `db:"updated_at"`
	ResetOtp        string `db:"reset_otp"`
	ResetExpiry     int    `db:"reset_expiry"`
}

func (db *Database) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	userRow := UserRow{}
	err := db.client.GetContext(ctx, &userRow, "SELECT id, first_name, last_name, email, is_email_verified, password, created_at, updated_at FROM public.users WHERE email=$1 LIMIT 1", email)
	if err != nil {
		log.Println(err)
		return user.User{}, err
	}
	return user.User{
		ID:              userRow.ID,
		FirstName:       userRow.FirstName,
		LastName:        userRow.LastName,
		Email:           userRow.Email,
		IsEmailVerified: userRow.IsEmailVerified,
		Password:        userRow.Password,
		CreatedAt:       userRow.CreatedAt,
		UpdatedAt:       userRow.UpdatedAt,
		ResetOtp:        userRow.ResetOtp,
		ResetExpiry:     userRow.ResetExpiry,
	}, nil
}

func (db *Database) InsertUser(ctx context.Context, firstName string, lastName string, email string, password string, isEmailVerified bool) (string, error) {
	userID := uuid.New().String()
	_, err := db.client.ExecContext(ctx, "INSERT INTO public.users (id, first_name, last_name, email, is_email_verified, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", userID, firstName, lastName, email, isEmailVerified, password, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		log.Println(err)
		return "", err
	}
	return userID, nil
}
