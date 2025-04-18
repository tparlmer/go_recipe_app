package auth

import (
	"errors"
	"time"

	"go_recipe_app/auth/db"

	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt/v5"
)

// Auth token lasts for 5 minutes
const (
	tokenDuration = 5 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// OMIT Address

// OMIT User profile

type User struct {
	UserID string `json:"userID"`
	Username string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	PasswordHash string `json:"passwordHash"`
	Email string `json:"email"`
	Roles []string `json:"roles"` // I will need an admin role to manually adjust other users recipes if necessary
	EmailVerified bool `json:"emailverified"`
}

type AuthClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	UserID string `json:"user_id"`
	Roles []string `json:"roles"`
}

// AuthService provides user authentication operations
type AuthService interface {
	// Authentication
	Login()
	RefreshToken()
	ValidateToken()

	// Registration
	Register()

	// Role Management
	AddRoleToUser()
	RemoveRoleFromuser()
	GetRolesForUser()

	// User Management
	GetUserProfile()
	UpdateUserProfile()
	GetUserByUsername()
	GetUserById()
	DeleteUser()
	ChangeUsername()

	// Password Management
	ChangePassword()

	// ResetPasswordInitiate takes an email address and returns a reset token that contains the username to be checked in ResetPasswordComplete
	ResetPasswordInitiate()

	// ResetPasswordComplete takes a reset token, username, and new password. It validates the token, checks that the username matches the token, and updates the user's password.
	ResetPasswordComplete()

	// Uses go-kit logs
	Close(logger log.Logger) error
}

// WHAT IS ??
type authServiceImpl struct {
	JWT_SECRET_KEY string
	logger log.Logger
	repo db.AuthRepository
}