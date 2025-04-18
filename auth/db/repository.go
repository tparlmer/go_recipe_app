package db

import "github.com/go-kit/log"

// OMIT Address

// OMIT UserProfile for now
// Profile represents easily changeable, non-sensitive information about a user

// User represents a user in the system.
type User struct {
	UserID string `json:"userID"`
	Username string `json:"username"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	PasswordHash string `json:"passwordHash"`
	Email string `json:"email"`
	Roles []string
	EmailVerified bool `json:"emailVerified"`

	// add optional fields here if desired
}

// Responsible for data access operations for auth
type AuthRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(userID string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]User, error)

	Close(logger log.Logger) error
}