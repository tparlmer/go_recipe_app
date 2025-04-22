package recipes

// ALL THESE HANDLERS ARE PASSED TO THE ROUTES IN setupRoutes() in main.go
// TODO: Need to implement auth_middleware in recipes package
// For this to be a hypermedia API I need to return fragments of html and links that maintain application state

/*
ARCHITECTURAL PSEUDOCODE - RECIPE HANDLERS WITH AUTH INTEGRATION

FUNCTION NewHandler(recipeService, templates, logger)
    Return initialized handler with services and templates
END FUNCTION

FUNCTION RegisterRoutes(router)
    // Public routes
    router.Handle("/recipes", handler.ListPublicRecipes())
    router.Handle("/recipes/{id}", handler.GetRecipe())

    // Protected routes
    protectedRouter := router.PathPrefix("/").Subrouter()
    protectedRouter.Use(authMiddleware)

    protectedRouter.Handle("/my-recipes", handler.ListUserRecipes())
    protectedRouter.Handle("/recipes/new", handler.CreateRecipeForm())
    protectedRouter.Handle("/recipes/create", handler.CreateRecipe())
    protectedRouter.Handle("/recipes/{id}/edit", handler.EditRecipeForm())
    protectedRouter.Handle("/recipes/{id}/update", handler.UpdateRecipe())
    protectedRouter.Handle("/recipes/{id}/delete", handler.DeleteRecipe())
END FUNCTION

FUNCTION ListPublicRecipes()
    RETURN FUNCTION(response, request)
        Get all public recipes from service
        Render template with recipes
    END FUNCTION
END FUNCTION

FUNCTION ListUserRecipes()
    RETURN FUNCTION(response, request)
        Get user from request context
        IF user not authenticated THEN
            Redirect to login
            RETURN
        END IF

        Get user's recipes from service
        Render template with recipes
    END FUNCTION
END FUNCTION

FUNCTION GetRecipe()
    RETURN FUNCTION(response, request)
        Extract recipe ID from URL
        Get user from request context (may be nil if not logged in)

        Get recipe with visibility check
        IF error THEN
            Show 404 page
            RETURN
        END IF

        Render template with recipe
        // If user is owner, template will show edit/delete controls
    END FUNCTION
END FUNCTION

FUNCTION CreateRecipe()
    RETURN FUNCTION(response, request)
        Get user from request context
        IF user not authenticated THEN
            Redirect to login
            RETURN
        END IF

        Parse form data
        Create recipe object
        Set owner to current user
        Set public/private based on form input

        Save recipe
        Redirect to view new recipe
    END FUNCTION
END FUNCTION

FUNCTION UpdateRecipe()
    RETURN FUNCTION(response, request)
        Get user from request context
        IF user not authenticated THEN
            Redirect to login
            RETURN
        END IF

        Extract recipe ID from URL
        Parse form data

        Update recipe (service checks ownership)
        IF unauthorized error THEN
            Show 403 forbidden page
            RETURN
        END IF

        Redirect to view updated recipe
    END FUNCTION
END FUNCTION

FUNCTION DeleteRecipe()
    RETURN FUNCTION(response, request)
        Get user from request context
        IF user not authenticated THEN
            Redirect to login
            RETURN
        END IF

        Extract recipe ID from URL

        Delete recipe (service checks ownership)
        IF unauthorized error THEN
            Show 403 forbidden page
            RETURN
        END IF

        Redirect to recipes list
    END FUNCTION
END FUNCTION
*/

// This file will contain the HTTP handlers for recipe-related routes
// It will integrate with auth middleware to secure protected routes

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"go_recipe_app/auth"
	"go_recipe_app/recipes/db"
)

// Handler manages HTTP requests for recipe operations
type Handler struct {
	service RecipeService
	tmpl    *template.Template
	logger  log.Logger
}

// TemplateData holds data to be passed to templates
type TemplateData struct {
	User            *auth.UserContext
	Recipes         []db.Recipe
	Recipe          *db.Recipe
	Error           string
	Success         string
	IsAuthenticated bool
	SearchQuery     string
	SortOption      string
	CurrentPage     int
	TotalPages      int
	CurrentYear     int
	CurrentTemplate string
	Debug           bool
}

// NewHandler creates a new recipe handler
func NewHandler(service RecipeService, tmpl *template.Template, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		tmpl:    tmpl,
		logger:  logger,
	}
}

// Home renders the home page
func (h *Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := h.service.ListPublicRecipes()
		if err != nil {
			h.logger.Log("msg", "Error fetching public recipes for home page", "err", err)
			http.Error(w, "Error fetching recipes", http.StatusInternalServerError)
			return
		}

		// Get a limited selection of recipes for the homepage
		featuredRecipes := recipes
		if len(featuredRecipes) > 5 {
			featuredRecipes = featuredRecipes[:5]
		}

		user := auth.GetUserFromContext(r.Context())
		data := TemplateData{
			User:            user,
			Recipes:         featuredRecipes,
			IsAuthenticated: user != nil,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "home.html",
			Debug:           true,
		}

		// Try to render the layout template instead of just the content
		err = h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering home page", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListPublicRecipes lists all public recipes
func (h *Handler) ListPublicRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := h.service.ListPublicRecipes()
		if err != nil {
			h.logger.Log("msg", "Error fetching public recipes", "err", err)
			http.Error(w, "Error fetching recipes", http.StatusInternalServerError)
			return
		}

		user := auth.GetUserFromContext(r.Context())
		data := TemplateData{
			User:            user,
			Recipes:         recipes,
			IsAuthenticated: user != nil,
			CurrentPage:     1,
			TotalPages:      1,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "recipes.html",
		}

		err = h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering recipes page", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListUserRecipes lists recipes owned by the current user
func (h *Handler) ListUserRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		recipes, err := h.service.ListUserRecipes(user.UserID)
		if err != nil {
			h.logger.Log("msg", "Error fetching user recipes", "userID", user.UserID, "err", err)
			http.Error(w, "Error fetching your recipes", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			User:            user,
			Recipes:         recipes,
			IsAuthenticated: true,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "my-recipes.html",
		}

		err = h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering my-recipes page", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ViewRecipe shows a single recipe
func (h *Handler) ViewRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		user := auth.GetUserFromContext(r.Context())
		var userID string
		if user != nil {
			userID = user.UserID
		}

		recipe, err := h.service.GetRecipe(id, userID)
		if err != nil {
			if errors.Is(err, ErrRecipeNotFound) {
				http.NotFound(w, r)
				return
			}
			h.logger.Log("msg", "Error fetching recipe", "id", id, "err", err)
			http.Error(w, "Error fetching recipe", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			User:            user,
			Recipe:          recipe,
			IsAuthenticated: user != nil,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "recipe-detail.html",
		}

		err = h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering recipe detail page", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateRecipeForm renders the recipe creation form
func (h *Handler) CreateRecipeForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		data := TemplateData{
			User:            user,
			IsAuthenticated: true,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "recipe-form.html",
		}

		err := h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering create recipe form", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateRecipe processes recipe creation
func (h *Handler) CreateRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.logger.Log("msg", "Error parsing form", "err", err)
			http.Error(w, "Error processing form", http.StatusBadRequest)
			return
		}

		// Basic recipe data
		title := r.FormValue("title")
		description := r.FormValue("description")
		isPublic := r.FormValue("is_public") == "on"

		// Parse times
		prepTimeStr := r.FormValue("prep_time")
		prepTime, err := strconv.Atoi(prepTimeStr)
		if err != nil {
			h.logger.Log("msg", "Invalid prep time", "value", prepTimeStr, "err", err)
			http.Error(w, "Invalid prep time", http.StatusBadRequest)
			return
		}

		cookTimeStr := r.FormValue("cook_time")
		cookTime, err := strconv.Atoi(cookTimeStr)
		if err != nil {
			h.logger.Log("msg", "Invalid cook time", "value", cookTimeStr, "err", err)
			http.Error(w, "Invalid cook time", http.StatusBadRequest)
			return
		}

		servingsStr := r.FormValue("servings")
		servings, err := strconv.Atoi(servingsStr)
		if err != nil {
			h.logger.Log("msg", "Invalid servings", "value", servingsStr, "err", err)
			http.Error(w, "Invalid servings", http.StatusBadRequest)
			return
		}

		// Create recipe object
		recipe := &db.Recipe{
			Title:       title,
			Description: description,
			PrepTime:    time.Duration(prepTime) * time.Minute,
			CookTime:    time.Duration(cookTime) * time.Minute,
			Servings:    int32(servings),
			IsPublic:    isPublic,
			// Ingredients and instructions would be handled here
		}

		// Save recipe
		err = h.service.CreateRecipe(recipe, user.UserID)
		if err != nil {
			h.logger.Log("msg", "Error creating recipe", "err", err)
			http.Error(w, "Error saving recipe", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/recipes/"+recipe.ID, http.StatusSeeOther)
	}
}

// EditRecipeForm renders the recipe edit form
func (h *Handler) EditRecipeForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		recipe, err := h.service.GetRecipe(id, user.UserID)
		if err != nil {
			if errors.Is(err, ErrRecipeNotFound) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(err, ErrUnauthorized) {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			h.logger.Log("msg", "Error fetching recipe for edit", "id", id, "err", err)
			http.Error(w, "Error fetching recipe", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			User:            user,
			Recipe:          recipe,
			IsAuthenticated: true,
			CurrentYear:     time.Now().Year(),
			CurrentTemplate: "recipe-form.html",
		}

		err = h.tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			h.logger.Log("msg", "Error rendering edit recipe form", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// UpdateRecipe processes recipe updates
func (h *Handler) UpdateRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		if err := r.ParseForm(); err != nil {
			h.logger.Log("msg", "Error parsing form", "err", err)
			http.Error(w, "Error processing form", http.StatusBadRequest)
			return
		}

		// Basic recipe data
		title := r.FormValue("title")
		description := r.FormValue("description")
		isPublic := r.FormValue("is_public") == "on"

		// Parse times
		prepTimeStr := r.FormValue("prep_time")
		prepTime, err := strconv.Atoi(prepTimeStr)
		if err != nil {
			h.logger.Log("msg", "Invalid prep time", "value", prepTimeStr, "err", err)
			http.Error(w, "Invalid prep time", http.StatusBadRequest)
			return
		}

		cookTimeStr := r.FormValue("cook_time")
		cookTime, err := strconv.Atoi(cookTimeStr)
		if err != nil {
			h.logger.Log("msg", "Invalid cook time", "value", cookTimeStr, "err", err)
			http.Error(w, "Invalid cook time", http.StatusBadRequest)
			return
		}

		servingsStr := r.FormValue("servings")
		servings, err := strconv.Atoi(servingsStr)
		if err != nil {
			h.logger.Log("msg", "Invalid servings", "value", servingsStr, "err", err)
			http.Error(w, "Invalid servings", http.StatusBadRequest)
			return
		}

		// Create recipe object
		recipe := &db.Recipe{
			ID:          id,
			Title:       title,
			Description: description,
			PrepTime:    time.Duration(prepTime) * time.Minute,
			CookTime:    time.Duration(cookTime) * time.Minute,
			Servings:    int32(servings),
			IsPublic:    isPublic,
			// Ingredients and instructions would be handled here
		}

		// Update recipe
		err = h.service.UpdateRecipe(recipe, user.UserID)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}
			if errors.Is(err, ErrRecipeNotFound) {
				http.NotFound(w, r)
				return
			}

			h.logger.Log("msg", "Error updating recipe", "id", id, "err", err)

			// Render the form again with the error
			data := TemplateData{
				User:            user,
				Recipe:          recipe,
				Error:           "Error updating recipe: " + err.Error(),
				IsAuthenticated: true,
				CurrentYear:     time.Now().Year(),
				CurrentTemplate: "recipe-form.html",
			}
			h.tmpl.ExecuteTemplate(w, "layout", data)
			return
		}

		http.Redirect(w, r, "/recipes/"+id, http.StatusSeeOther)
	}
}

// DeleteRecipe processes recipe deletion
func (h *Handler) DeleteRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		err := h.service.DeleteRecipe(id, user.UserID)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}
			if errors.Is(err, ErrRecipeNotFound) {
				http.NotFound(w, r)
				return
			}

			h.logger.Log("msg", "Error deleting recipe", "id", id, "err", err)
			http.Error(w, "Error deleting recipe", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/my-recipes", http.StatusSeeOther)
	}
}
