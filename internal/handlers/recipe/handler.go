// internal/handlers/recipe/handler.go

package recipe

import (
	"fmt"
	"go_recipe_app/internal/models"
	"go_recipe_app/internal/storage"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TemplateData is a struct that holds the data for the template
// This is a common design pattern in Go to pass data to templates
type TemplateData struct {
	Template string
	Data     interface{}
}

// RecipeHandler holds all dependencies for recipe handling
// Struct is like a class in OOP
type RecipeHandler struct {
	tmpl   *template.Template
	logger *log.Logger
	Router *mux.Router // capitalize the first letter to export it
	store  storage.RecipeStore
}

// new creates a new RecipeHandler
// This is a constructor function that initializes the RecipeHandler struct with the necessary dependencies
func New(tmpl *template.Template, store storage.RecipeStore) *RecipeHandler {
	h := &RecipeHandler{
		tmpl:   tmpl,
		logger: log.New(os.Stdout, "[RECIPE] ", log.LstdFlags),
		Router: mux.NewRouter(),
		store:  store,
	}
	h.setupRoutes()
	return h
}

// setupRoutes registers all routes with the recipe handler
func (h *RecipeHandler) setupRoutes() {
	// Root route currently redirects to /recipes
	h.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/recipes", http.StatusSeeOther)
	}).Methods("GET")
	h.Router.HandleFunc("/recipes", h.listRecipes).Methods("GET")          // list all recipes
	h.Router.HandleFunc("/recipes/new", h.createRecipeForm).Methods("GET") // Show create form
	h.Router.HandleFunc("/recipes", h.createRecipe).Methods("POST")        // Handle form submission
	h.Router.HandleFunc("/recipes/{id}", h.getRecipe).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}/edit", h.editRecipeForm).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}", h.updateRecipe).Methods("PUT")
	h.Router.HandleFunc("/recipes/{id}", h.deleteRecipe).Methods("DELETE")
}

// Basic handler for listing recipes
func (h *RecipeHandler) listRecipes(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling list recipes request")

	recipes, err := h.store.List()
	if err != nil {
		h.logger.Printf("Error listing recipes: %v", err)
		http.Error(w, "Error getting recipes", http.StatusInternalServerError)
		return
	}
	// h.logger.Printf("Found %d recipe", len(recipes))

	data := TemplateData{
		Template: "list",
		Data:     recipes,
	}

	err = h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Println("Successfully rendered list of test recipes")
}

// Basic handler for getting a single recipe
func (h *RecipeHandler) getRecipe(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling get recipe request")

	// Get recipe ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Printf("Looking for recipe with ID: %s", id)

	// Get recipe from store
	recipe, err := h.store.Get(id)
	if err != nil {
		h.logger.Printf("Error getting recipe: %v", err)
		http.Error(w, "Error getting recipe", http.StatusInternalServerError)
		return
	}

	// Render recipe
	data := TemplateData{
		Template: "view",
		Data:     recipe,
	}

	// Execute template
	err = h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Printf("Successfully rendered recipe: %s", recipe.Title)
}

// Show the create recipe form
func (h *RecipeHandler) createRecipeForm(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Template: "create",
		Data:     nil,
	}

	err := h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle the form submission
func (h *RecipeHandler) createRecipe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Printf("Error parsing form: %v", err)
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	// Generate a unique ID (we'll improve this later)
	id := fmt.Sprintf("recipe-%d", time.Now().Unix())

	// Parse form values with error handling
	prepTime, err := time.ParseDuration(r.FormValue("prep_time") + "m")
	if err != nil {
		h.logger.Printf("Invalid prep time: %v", err)
		http.Error(w, "Invalid prep time", http.StatusBadRequest)
		return
	}

	cookTime, err := time.ParseDuration(r.FormValue("cook_time") + "m")
	if err != nil {
		h.logger.Printf("Invalid cook time: %v", err)
		http.Error(w, "Invalid cook time", http.StatusBadRequest)
		return
	}

	servings, err := strconv.Atoi(r.FormValue("servings"))
	if err != nil {
		h.logger.Printf("Invalid servings: %v", err)
		http.Error(w, "Invalid servings", http.StatusBadRequest)
		return
	}

	recipe := models.Recipe{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		PrepTime:    prepTime,
		CookTime:    cookTime,
		Servings:    int32(servings),
		// We'll add ingredients and instructions later
	}

	// Validate required fields
	if recipe.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Store the recipe
	if err := h.store.Create(recipe); err != nil {
		h.logger.Printf("Error creating recipe: %v", err)
		http.Error(w, "Error saving recipe", http.StatusInternalServerError)
		return
	}

	// Redirect to the new recipe
	http.Redirect(w, r, "/recipes/"+recipe.ID, http.StatusSeeOther)
}

// Show the edit form
func (h *RecipeHandler) editRecipeForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	recipe, err := h.store.Get(id)
	if err != nil {
		h.logger.Printf("Error getting recipe to edit: %v", err)
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	data := TemplateData{
		Template: "edit",
		Data:     recipe,
	}

	err = h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle the update
func (h *RecipeHandler) updateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := r.ParseForm(); err != nil {
		h.logger.Printf("Error parsing form: %v", err)
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	// Parse form values with error handling
	prepTime, err := time.ParseDuration(r.FormValue("prep_time") + "m")
	if err != nil {
		h.logger.Printf("Invalid prep time: %v", err)
		http.Error(w, "Invalid prep time", http.StatusBadRequest)
		return
	}

	cookTime, err := time.ParseDuration(r.FormValue("cook_time") + "m")
	if err != nil {
		h.logger.Printf("Invalid cook time: %v", err)
		http.Error(w, "Invalid cook time", http.StatusBadRequest)
		return
	}

	servings, err := strconv.Atoi(r.FormValue("servings"))
	if err != nil {
		h.logger.Printf("Invalid servings: %v", err)
		http.Error(w, "Invalid servings", http.StatusBadRequest)
		return
	}

	recipe := models.Recipe{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		PrepTime:    prepTime,
		CookTime:    cookTime,
		Servings:    int32(servings),
	}

	if err := h.store.Update(recipe); err != nil {
		h.logger.Printf("Error updating recipe: %v", err)
		http.Error(w, "Error updating recipe", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/recipes/"+recipe.ID, http.StatusSeeOther)
}

// Delete recipe handler
func (h *RecipeHandler) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Printf("Attempting to delete recipe: %s", id)

	// Check if recipe exists
	_, err := h.store.Get(id)
	if err != nil {
		h.logger.Printf("Recipe not found for deletion: %v", err)
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Delete the recipe
	if err := h.store.Delete(id); err != nil {
		h.logger.Printf("Error deleting recipe: %v", err)
		http.Error(w, "Error deleting recipe", http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully deleted recipe: %s", id)
	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}
