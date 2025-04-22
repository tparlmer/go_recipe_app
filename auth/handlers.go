package auth

/*
ARCHITECTURAL PSEUDOCODE - AUTH HANDLERS

FUNCTION NewHandler(authService, templates, logger)
    Return initialized handler with auth service and templates
END FUNCTION

FUNCTION RegisterRoutes(router)
    router.Handle("/login", handler.LoginForm()).Methods("GET")
    router.Handle("/login", handler.Login()).Methods("POST")
    router.Handle("/register", handler.RegisterForm()).Methods("GET")
    router.Handle("/register", handler.Register()).Methods("POST")
    router.Handle("/logout", handler.Logout())
END FUNCTION

FUNCTION LoginForm()
    RETURN FUNCTION(response, request)
        Render login template
    END FUNCTION
END FUNCTION

FUNCTION Login()
    RETURN FUNCTION(response, request)
        Parse form values (username, password)

        Call auth service to validate credentials
        IF error THEN
            Render login form with error message
            RETURN
        END IF

        Generate JWT token
        Set token in cookie
        Redirect to home page
    END FUNCTION
END FUNCTION

FUNCTION RegisterForm()
    RETURN FUNCTION(response, request)
        Render registration template
    END FUNCTION
END FUNCTION

FUNCTION Register()
    RETURN FUNCTION(response, request)
        Parse form values (username, password, email, etc)
        Validate input data

        Call auth service to create user
        IF error (e.g., username already exists) THEN
            Render registration form with error message
            RETURN
        END IF

        Redirect to login page with success message
    END FUNCTION
END FUNCTION

FUNCTION Logout()
    RETURN FUNCTION(response, request)
        Clear auth cookie
        Redirect to home page
    END FUNCTION
END FUNCTION
*/

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-kit/log"
)

// Handler manages HTTP requests for auth operations
type Handler struct {
	service AuthService
	tmpl    *template.Template
	logger  log.Logger
}

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Error           string
	Success         string
	User            *UserContext
	CurrentURL      string
	CurrentYear     int
	CurrentTemplate string
	Debug           bool
}

// NewHandler creates a new auth handler
func NewHandler(service AuthService, tmpl *template.Template, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		tmpl:    tmpl,
		logger:  logger,
	}
}

// LoginForm renders the login form
func (h *Handler) LoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Log("msg", "Rendering login form template")

		// Check for success message from registration
		success := ""
		if r.URL.Query().Get("success") == "1" {
			success = "Registration successful! Please log in."
		}

		data := TemplateData{
			CurrentURL:      r.URL.Path,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "login.html",
			Success:         success,
			Debug:           true,
		}

		// Render using layout template
		err := h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering login form", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Login processes login form submissions
func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Log("msg", "Error parsing login form", "err", err)
			data := TemplateData{
				Error:           "Invalid form submission",
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "login.html",
			}
			err = h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering login form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Attempt login
		token, expiry, userID, firstName, lastName, err := h.service.Login(username, password, nil)
		if err != nil {
			h.logger.Log("msg", "Login failed", "username", username, "err", err)
			data := TemplateData{
				Error:           "Invalid username or password",
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "login.html",
			}
			err = h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering login form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Set auth cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Expires:  time.Now().Add(time.Duration(expiry) * time.Second),
			HttpOnly: true,
			Path:     "/",
		})

		h.logger.Log("msg", "User logged in", "username", username, "userID", userID, "firstName", firstName, "lastName", lastName)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// RegisterForm renders the registration form
func (h *Handler) RegisterForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Log("msg", "Rendering register form template")

		data := TemplateData{
			CurrentURL:      r.URL.Path,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "register.html",
			Debug:           true,
		}

		// Render using layout template
		err := h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering registration form", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Register processes registration form submissions
func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			h.logger.Log("msg", "Error parsing registration form", "err", err)
			data := TemplateData{
				Error:           "Invalid form submission",
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "register.html",
			}
			err = h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering registration form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		email := r.FormValue("email")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")

		// Basic validation
		if username == "" || password == "" || email == "" {
			data := TemplateData{
				Error:           "Username, password, and email are required",
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "register.html",
			}
			err := h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering registration form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if password != confirmPassword {
			data := TemplateData{
				Error:           "Passwords do not match",
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "register.html",
			}
			err := h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering registration form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Register the user
		err := h.service.Register(username, password, email, firstName, lastName, []string{"user"})
		if err != nil {
			h.logger.Log("msg", "Registration failed", "username", username, "err", err)
			data := TemplateData{
				Error:           "Registration failed: " + err.Error(),
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "register.html",
			}
			err = h.tmpl.ExecuteTemplate(w, "layout", data)
			if err != nil {
				h.logger.Log("msg", "Error rendering registration form with error", "err", err)
				http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		h.logger.Log("msg", "User registered", "username", username)

		// Redirect to login with success message
		http.Redirect(w, r, "/login?success=1", http.StatusSeeOther)
	}
}

// Logout logs out a user
func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear auth cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Path:     "/",
		})

		h.logger.Log("msg", "User logged out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
