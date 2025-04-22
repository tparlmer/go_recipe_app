package auth

import (
	"errors"
	"fmt"
	"time"

	"go_recipe_app/auth/db"

	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// Auth token lasts for 24 hours
const (
	tokenDuration = 24 * time.Hour
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// OMIT Address

// OMIT User profile

// type UserProfile struct {
// 	ShippingAddress *Address `json:"address,omitempty"`
// 	OrganizationName string `json:"organizationName,omitempty"`
// 	PhoneNumber string `json:"phoneNumber,omitempty"`
// }

type User struct {
	UserID        string   `json:"userID"`
	Username      string   `json:"username"`
	Firstname     string   `json:"firstname"`
	Lastname      string   `json:"lastname"`
	PasswordHash  string   `json:"passwordHash"`
	Email         string   `json:"email"`
	Roles         []string `json:"roles"` // I will need an admin role to manually adjust other users recipes if necessary
	EmailVerified bool     `json:"emailverified"`
}

type AuthClaims struct {
	jwt.RegisteredClaims
	Username string   `json:"username"`
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
}

// AuthService provides user authentication operations
type AuthService interface {
	// Authentication
	Login(username, password string, roles []string) (string, int, string, string, string, error)
	RefreshToken(token string) (refreshedToken string, expiresInSec int, err error)
	ValidateToken(token string) (*AuthClaims, error)

	// Registration
	Register(username, password, email, firstName, lastName string, roles []string) error

	// Role Management
	// THESE FUNCTIONS/METHODS STILL NEED TO BE CREATED
	// AddRoleToUser(userID, role string) error - TODO: ADD
	// RemoveRoleFromuser(userID, role string) error - TODO: ADD
	// GetRolesForUser(userID string) ([]string, error) - TODO: ADD

	// User Management
	// GetUserProfile(userID string) (*UserProfile, error)
	// UpdateUserProfile(UserID string, profile *UserProfile) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(userID string) (*db.User, error)
	// DeleteUser() - TODO: ADD
	// ChangeUsername() - TODO: ADD

	// Password Management
	// Change password for user with userID and oldPassword
	ChangePassword(userID, oldPassword, newPassword string) error

	// ResetPasswordInitiate takes an email address and returns a reset token that contains the username to be checked in ResetPasswordComplete
	ResetPasswordInitiate(email string) (string, error)

	// ResetPasswordComplete takes a reset token, username, and new password. It validates the token, checks that the username matches the token, and updates the user's password.
	ResetPasswordComplete(resetToken string, username string, newPassword string) error

	// Uses go-kit logs
	Close(logger log.Logger) error
}

// WHAT IS ??
type authServiceImpl struct {
	JWT_SECRET_KEY string
	logger         log.Logger
	repo           db.AuthRepository
}

// Constructor function for anew instance of authService
func NewAuthService(JWTsecretKey string, dataDir string, mainLogger log.Logger, authLogger log.Logger) (AuthService, error) {
	repo, err := db.NewBoltAuthRepository(dataDir, mainLogger)
	if err != nil {
		mainLogger.Log("msg", "Could not initialize auth data repository", "err", err)
		return nil, err
	}

	return &authServiceImpl{
		JWT_SECRET_KEY: JWTsecretKey,
		logger:         authLogger,
		repo:           repo,
	}, nil
}

func (s *authServiceImpl) Close(mainLogger log.Logger) error {
	err := s.repo.Close(mainLogger)
	if err != nil {
		mainLogger.Log("msg", "Could not close auth data repository", "err", err)
		return err
	}
	mainLogger.Log("msg", "Closed auth data repository")
	return nil
}

// Login validates credentials and returns a JWT.
func (s *authServiceImpl) Login(username, password string, roles []string) (string, int, string, string, string, error) {
	s.logger.Log("msg", "Starting login", "username", username)
	var user *db.User
	// Retrieve user by supplied username
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", 0, "", "", "", err
	}
	s.logger.Log("msg", "User found", "username", username)
	if len(password) > 72 {
		return "", 0, "", "", "", bcrypt.ErrPasswordTooLong
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	s.logger.Log("msg", "Password compared", "username", username)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", 0, "", "", "", ErrInvalidCredentials
	}
	if err != nil {
		return "", 0, "", "", "", fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}
	s.logger.Log("msg", "Password matched", "username", username)
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: user.Username,
		UserID:   user.UserID,
		Roles:    user.Roles,
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s.logger.Log("msg", "Token created", "username", username, "token", token)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(s.JWT_SECRET_KEY))
	s.logger.Log("msg", "Token signed", "username", username, "tokenString", tokenString)
	expiresInSec := int(tokenDuration.Seconds())
	if err != nil {
		return "", 0, "", "", "", err
	}
	s.logger.Log("msg", "Login Successful", "username", username)
	return tokenString, expiresInSec, user.UserID, user.FirstName, user.LastName, nil
}

// Register creates a new user record.
func (s *authServiceImpl) Register(username, password, email, firstName, lastName string, roles []string) error {
	if len(password) > 72 {
		return bcrypt.ErrPasswordTooLong
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create a new user with new user ID
	user := &db.User{
		UserID:        uuid.New().String(),
		Username:      username,
		FirstName:     firstName,
		LastName:      lastName,
		PasswordHash:  string(passwordHash),
		Email:         email,
		Roles:         roles,
		EmailVerified: false,

		// Initialize optional fields with zero values
		// Profile: nil,
	}

	return s.repo.CreateUser(user)
}

// RefreshToken refreshes a JWT for a user.
func (s *authServiceImpl) RefreshToken(tokenString string) (string, int, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWT_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return "", 0, errors.New("invalid token claims")
	}
	// Create a new token with refreshed expiration time
	refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Username: claims.Username,
		UserID:   claims.UserID,
		Roles:    claims.Roles,
	})

	refreshedTokenString, err := refreshedToken.SignedString([]byte(s.JWT_SECRET_KEY))
	expiresInSec := int(tokenDuration.Seconds())
	if err != nil {
		return "", 0, err
	}

	return refreshedTokenString, expiresInSec, nil
}

// ValidateToken validates a JWT for a user and returns the decoded claims.
func (s *authServiceImpl) ValidateToken(tokenString string) (*AuthClaims, error) {
	var claims AuthClaims

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWT_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}

// GetUserProfile retrieves user profile
// func (s *authServiceImpl) GetUserProfile(userID string) (*UserProfile, error) {
// USER PROFILE NOT CURRENTLY USED
// }

// UpdateUserPRofile update's a user's profile. (cannot change anything else)
// func (s *authServiceImpl) UpdateUserProfile(userID sting, profile *UserProfile) error {
// USER PROFILE NOT CURRENTLY USED
// }

// ChangePassword changes a user's password.
func (s *authServiceImpl) ChangePassword(userID, oldPassword, newPassword string) error {
	s.logger.Log("msg", "Changing password", "userID", userID)
	if len(newPassword) > 72 {
		return bcrypt.ErrPasswordTooLong
	}

	// Get the user from the repository
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Compare the old password with the current password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidCredentials
	}
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}

	// Generate a new password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password hash
	user.PasswordHash = string(passwordHash)
	return s.repo.UpdateUser(user)
}

// ResetPasswordInitiate takes an email address and returns a reset token that contains the username to be checked in ResetPasswordComplete
func (s *authServiceImpl) ResetPasswordInitiate(email string) (string, error) {
	s.logger.Log("msg", "Reset password initiated", "email", email)

	// Get the user from the repository
	user, err := s.repo.GetUserByUsername(email)
	if err != nil {
		return "", err
	}

	// Genearte a jwt token with the username as the subject
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.UserID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(s.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ResetPasswordComplete
// It validates the token, checks that the username matches the token, and updates the user's password.
func (s *authServiceImpl) ResetPasswordComplete(resetToken string, username string, newPassword string) error {
	s.logger.Log("msg", "Reset password completion attempted.", "username", username)

	// Parse the token with claims
	token, err := jwt.Parse(resetToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWT_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Check that the username matches the token
	if claims["username"] != username {
		return errors.New("username does nto match token")
	}

	// Get the user from the repository
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	// Generate a new password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password hash
	user.PasswordHash = string(passwordHash)
	return s.repo.UpdateUser(user)
}

// GetUserByUsername returns a user by username (with sensitive fields redacted).
func (s *authServiceImpl) GetUserByUsername(username string) (*User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	svcUser := convertToServiceUser(user)
	svcUser.PasswordHash = ""
	return svcUser, nil
}

// GetUserByID returns a user by ID
func (s *authServiceImpl) GetUserByID(userID string) (*db.User, error) {
	return s.repo.GetUserByID(userID)
}

// ----------
// HELPER FUNCTIONS - MOVE TO helper.go
// ----------

// func convertToServiceProfile(profile *db.UserProfile) *UserProfile {
// 	if profile == nil {
// 		return nil
// 	}
// 	var convertedAddress *Address
// 	if profile.ShippingAddress != nil {
// 		convertedAddress = convertToServiceAddress(profile.ShippingAddress)
// 	} else {
// 		convertedAddress = nil
// 	}

// 	return &UserProfile{
// 		ShippingAddress:  convertedAddress,
// 		OrganizationName: profile.OrganizationName,
// 		PhoneNumber:      profile.PhoneNumber,
// 	}
// }

func convertToServiceUser(user *db.User) *User {
	// var convertedProfile *UserProfile
	// if user.Profile != nil {
	// 	convertedProfile = convertToServiceProfile(user.Profile)
	// } else {
	// 	convertedProfile = nil
	// }
	return &User{
		UserID:   user.UserID,
		Username: user.Username,
		// FirstName:     user.FirstName,
		// LastName:      user.LastName,
		PasswordHash:  user.PasswordHash,
		Email:         user.Email,
		Roles:         user.Roles,
		EmailVerified: user.EmailVerified,
		// Profile:       convertedProfile,
	}
}

func errorToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}