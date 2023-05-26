package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserCreationFailed   = errors.New("unable to create new user")
	UnableToFindUser     = errors.New("unable to find user")
	MethodNotImplemented = errors.New("method not implemented")
	InvalidPassword      = errors.New("wrong password")
	PasswordHashFailed   = errors.New("failed while hashing password")
	TokenGenFailed       = errors.New("unable to generate token")
	InvalidRefreshToken  = errors.New("invalid refresh token")
	UserCreated          = "user created"
)

type User struct {
	ID              string
	FirstName       string
	LastName        string
	Email           string
	IsEmailVerified bool
	Password        string
	ResetOtp        string
	ResetExpiry     int
	CreatedAt       int
	UpdatedAt       int
}

type UserStore interface {
	GetUserByEmail(context.Context, string) (User, error)
	InsertUser(context.Context, string, string, string, string, bool) (string, error)
}

type Service struct {
	store UserStore
}

func New(store UserStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return User{}, UnableToFindUser
	}
	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, firstName string, lastName string, email string, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", PasswordHashFailed
	}
	userID, err := s.store.InsertUser(ctx, firstName, lastName, email, string(hashedPassword), false)
	if err != nil {
		log.Println(err)
		return "", UserCreationFailed
	}
	return userID, nil
}

func (s *Service) GenerateAccessToken(ctx context.Context, refreshToken string) (string, string, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Provide the JWT secret key for token verification
		return []byte("refreshtokensecret"), nil
	})

	refreshTokenClaims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", "", InvalidRefreshToken
	}

	if err != nil {
		log.Println(err)
		return "", "", InvalidRefreshToken
	}

	if !token.Valid {
		return "", "", InvalidRefreshToken
	}

	userID, ok := refreshTokenClaims["id"].(string)
	if !ok {
		return "", "", InvalidRefreshToken
	}

	// Generate a new access token
	currentTime := time.Now()
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": currentTime.Add(time.Hour).Unix(),
		"iat": currentTime.Unix(),
	})
	accessToken, err := accessTokenClaims.SignedString([]byte("accesstokensecret")) // Replace with your actual access token secret key
	if err != nil {
		return "", "", TokenGenFailed
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return "", "", UnableToFindUser
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return "", "", InvalidPassword
	}

	currentTime := time.Now()

	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   currentTime.Add(time.Hour).Unix(),
		"iat":   currentTime.Unix(),
	})

	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Date(currentTime.Year(), currentTime.Month()+1, currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix(),
		"iat": currentTime.Unix(),
	})

	accessToken, err := accessTokenClaims.SignedString([]byte("accesstokensecret"))
	if err != nil {
		log.Println(err)
		return "", "", TokenGenFailed
	}
	refreshToken, err := refreshTokenClaims.SignedString([]byte("refreshtokensecret"))
	if err != nil {
		log.Println(err)
		return "", "", TokenGenFailed
	}

	return accessToken, refreshToken, nil
}
