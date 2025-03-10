// internal/handlers/recipe/handler.go

package recipe

import (
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go_recipe_app/internal/models"
	"os"
)

// RecipeHandler holds all dependencies for recipe handling
type RecipeHandler struct {
	tmpl *template.Template
	logger *log.Logger
	router *mux.Router
	// store will be added when we implement storage
}

// new creates a new RecipeHandler
// This is a constructor function that initializes the RecipeHandler struct with the necessary dependencies
func New(tmpl *template.Template, logger *log.Logger, router *mux.Router) *RecipeHandler {
	h := &RecipeHandler{
		tmpl: tmpl,
		logger: log.New(os.Stdout, "[RECIPE] ", log.LstdFlags),
		router: mux.NewRouter(),
	}
	h.setupRoutes()
	return h
}

// setupRoutes registers all routes with the recipe handler
func (h *RecipeHandler) setupRoutes() {
	h.router.HandleFunc("/recipes", h.listRecipes).Methods("GET")
	h.router.HandleFunc("/recipes/{id}", h.getRecipe).Methods("GET")
	// More routes will be added later
}

// Basic handler for listing recipes
func (h *RecipeHandler) listRecipes(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handlign list recipes request")
	// TODO: Implement the logic to list recipes
}

// Basic handler for getting a single recipe
func (h *RecipeHandler) getRecipe(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling get recipe request")
	// TODO: Implement the logic to get a single recipe
}