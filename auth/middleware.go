package auth

/*
ARCHITECTURAL PSEUDOCODE - AUTH MIDDLEWARE

FUNCTION authMiddleware(request, next_handler)
    Extract token from request (cookie or header)
    IF token is invalid OR expired THEN
        Redirect to login page
        RETURN
    END IF

    Decode user info from token
    Add user to request context
    Call next handler with enhanced request
END FUNCTION

FUNCTION requireRole(roles []string) MIDDLEWARE
    RETURN FUNCTION(request, next_handler)
        Get user from request context
        IF user is nil OR user doesn't have required role THEN
            Return unauthorized error
        END IF

        Call next handler
    END FUNCTION
END FUNCTION

FUNCTION getUserFromContext(context)
    Extract and return user object from context
    IF no user found in context THEN
        Return nil
    END IF
    Return user
END FUNCTION
*/

// This file will contain the authentication middleware implementation
// It will handle token validation and user context propagation

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go_recipe_app/auth/db"

	"github.com/go-kit/log"
)

// User context key to prevent collisions in context values
type contextKey string

const UserContextKey contextKey = "user"

// User information stored in context
type UserContext struct {
	UserID    string
	Username  string
	FirstName string
	LastName  string
	Roles     []string
}

// AuthMiddlewareService interface for middleware integration
type AuthMiddlewareService interface {
	ValidateToken(token string) (*AuthClaims, error)
	GetUserByID(userID string) (*db.User, error)
}

// AuthMiddleware creates a middleware that validates JWT tokens and adds user info to request context
func AuthMiddleware(service AuthMiddlewareService, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			logger.Log("msg", "Auth middleware processing request", "path", path)

			// Skip authentication for login and register pages
			if path == "/login" || path == "/register" {
				logger.Log("msg", "Skipping auth check for login/register page", "path", path)
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from cookie or Authorization header
			token, err := extractToken(r)
			if err != nil {
				logger.Log("msg", "No authentication token found, redirecting to login", "path", path, "err", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Validate token
			claims, err := service.ValidateToken(token)
			if err != nil {
				logger.Log("msg", "Invalid authentication token, redirecting to login", "path", path, "err", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Get user details
			user, err := service.GetUserByID(claims.UserID)
			if err != nil {
				logger.Log("msg", "User not found for valid token, redirecting to login", "path", path, "userID", claims.UserID, "err", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			logger.Log("msg", "User authenticated successfully", "path", path, "username", user.Username)

			// Create user context
			userCtx := &UserContext{
				UserID:    user.UserID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Roles:     user.Roles,
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), UserContextKey, userCtx)

			// Call next handler with enhanced context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(roles []string, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				logger.Log("msg", "Unauthenticated access attempt to protected resource")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range user.Roles {
					if requiredRole == userRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				logger.Log("msg", "Unauthorized access attempt", "user", user.Username, "requiredRoles", strings.Join(roles, ","))
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user information from request context
func GetUserFromContext(ctx context.Context) *UserContext {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

// Helper function to extract token from request
func extractToken(r *http.Request) (string, error) {
	// First try to get from cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		return cookie.Value, nil
	}

	// Then try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Usually in format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	return "", errors.New("no authentication token found")
}
