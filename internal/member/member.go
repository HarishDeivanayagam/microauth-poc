package member

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"microauth.io/core/internal/user"
)

type Role string

const (
	Admin Role = "admin"
	Staff Role = "staff"
	User  Role = "user"
)

var (
	AdminPermissionFailed   = errors.New("you aren't a administrator")
	FetchMemberFailed       = errors.New("unable to fetch member")
	MemberCreateFailed      = errors.New("unable to create member")
	MemberAdded             = "member added"
	MemberUpdateFailed      = errors.New("unable to update member")
	MemberUpdated           = "member updated"
	MemberDeleteFailed      = errors.New("unable to delete member")
	MemberDeleted           = "member deleted"
	InviteSent              = "invite sent successfully"
	InviteFailed            = errors.New("unable to send invite")
	FetchMemberInviteFailed = errors.New("unable to fetch member invite")
	UserCreationFailed      = errors.New("unable to create new user")
	InvalidInviteCode       = errors.New("invalid invite code")
)

type Member struct {
	// member id must not be used only for reference
	ID string

	OrganizationID string
	UserID         string
	Role           Role
	AppRole        string
}

type MemberInvite struct {
	ID             string
	Email          string
	Code           string
	OrganizationID string
	CreatedAt      int
	UpdatedAt      int
	ExpiresAt      int
}

type MemberStore interface {
	FetchMemberByID(context.Context, string, string) (Member, error)
	FetchAllMembers(context.Context, string) ([]Member, error)
	InsertMember(context.Context, string, string, Role, string) (string, error)
	UpdateMember(context.Context, string, string, Role, string) (string, error)
	DeleteMember(context.Context, string, string) (string, error)
	GetMemberInvite(context.Context, string, string) (MemberInvite, error)
	InsertMemberInvite(context.Context, string, string, string, int) (string, error)
	DeleteMemberInvite(context.Context, string, string) (string, error)
}

type UserService interface {
	GetUserByEmail(context.Context, string) (user.User, error)
	CreateUser(context.Context, string, string, string, string) (string, error)
}

type EmailService interface {
	SendEmail(string, string, string) (string, error)
}

type Service struct {
	store        MemberStore
	userService  UserService
	emailService EmailService
}

func New(store MemberStore, userService UserService, emailService EmailService) *Service {
	return &Service{
		store:        store,
		userService:  userService,
		emailService: emailService,
	}
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())

	otpLength := 6
	otpChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	otp := make([]byte, otpLength)

	for i := 0; i < otpLength; i++ {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}

	return string(otp)
}

func (s *Service) InviteMember(ctx context.Context, email string, userID string, organizationID string) (string, error) {
	// Check if the email exists
	_, err := s.userService.GetUserByEmail(ctx, email)
	newUser := false
	if err != nil {
		newUser = true
	}

	// Check access
	member, err := s.store.FetchMemberByID(ctx, organizationID, userID)
	if err != nil {
		return "", FetchMemberFailed
	}
	if member.Role != Admin {
		return "", AdminPermissionFailed
	}

	// Generate an OTP
	otp := generateOTP()

	// Store the OTP and insert member invite with an expiry of 72 hours
	expiresAt := int(time.Now().Add(72 * time.Hour).Unix())
	_, err = s.store.InsertMemberInvite(ctx, email, organizationID, otp, expiresAt)
	if err != nil {
		return "", err
	}

	// Construct the invitation URL
	clientURL := "https://example.com" // Replace with your actual client URL
	invitationURL := fmt.Sprintf("%s/auth/login?organizationID=%s&otp=%s", clientURL, organizationID, otp)
	if newUser {
		invitationURL += "&new_user=true"
	}

	// Send the OTP in the email
	subject := "Invitation OTP"
	body := fmt.Sprintf("Your invitation OTP: %s\n\nTo accept the invitation, please click the following link:\n%s", otp, invitationURL)
	_, err = s.emailService.SendEmail(email, subject, body)
	if err != nil {
		return "", InviteFailed
	}

	return InviteSent, nil
}

func (s *Service) AcceptInvite(ctx context.Context, email string, code string, organizationID string, firstName string, lastName string, password string) (string, error) {
	// Check MemberInvite for error
	memberInvite, err := s.store.GetMemberInvite(ctx, email, organizationID)
	if err != nil {
		log.Println(err)
		return "", FetchMemberInviteFailed
	}

	// Check if the code matches with the entry in the table
	if memberInvite.Code != code {
		log.Println(err)
		return "", InvalidInviteCode
	}

	// Check if the user exists for the given email
	existingUser, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		// Create a new user
		newUserID, err := s.userService.CreateUser(ctx, firstName, lastName, email, password)
		if err != nil {
			log.Println(err)
			return "", UserCreationFailed
		}
		// Set the user ID as the newly created user
		existingUser.ID = newUserID
	}

	// Add the member with the userID and role
	defaultRole := "user"
	_, err = s.store.InsertMember(ctx, organizationID, existingUser.ID, Role(defaultRole), "")
	if err != nil {
		log.Println(err)
		return "", MemberCreateFailed
	}

	// Delete the member invite entry
	_, err = s.store.DeleteMemberInvite(ctx, email, organizationID)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return MemberAdded, nil
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
