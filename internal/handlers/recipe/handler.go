// internal/handlers/recipe/handler.go

package recipe

import (
	"fmt"
	"go_recipe_app/internal/models"
	"go_recipe_app/internal/storage"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TemplateData is a struct that holds the data for the template
// This is a common design pattern in Go to pass data to templates
type TemplateData struct {
	Recipe      *models.Recipe
	Recipes     []models.Recipe
	CanEdit     bool
	Error       string
	Success     string
	CurrentPage int
	TotalPages  int
	CurrentYear int
	SearchQuery string
	SortOption  string
}

// RecipeHandler holds all dependencies for recipe handling
// Struct is like a class in OOP
type RecipeHandler struct {
	tmpl   *template.Template
	logger *slog.Logger
	Router *mux.Router // capitalize the first letter to export it
	store  storage.RecipeStore
}

// new creates a new RecipeHandler
// This is a constructor function that initializes the RecipeHandler struct with the necessary dependencies
func New(tmpl *template.Template, store storage.RecipeStore, logger *slog.Logger) *RecipeHandler {
	h := &RecipeHandler{
		tmpl:   tmpl,
		logger: logger,
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
	h.Router.HandleFunc("/recipes/create", h.createRecipe).Methods("POST") // Handle form submission
	h.Router.HandleFunc("/recipes/{id}", h.getRecipe).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}/edit", h.editRecipeForm).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}/edit", h.updateRecipe).Methods("POST")
	h.Router.HandleFunc("/recipes/{id}/delete", h.deleteRecipe).Methods("POST")
}

// Basic handler for listing recipes
func (h *RecipeHandler) listRecipes(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling list recipes request")

	// Get search and sort parameters
	searchQuery := r.URL.Query().Get("search")
	sortOption := r.URL.Query().Get("sort")

	// For now, we'll just list all recipes - in a real app you'd filter based on search
	recipes, err := h.store.List()
	if err != nil {
		h.logger.Error("Error listing recipes", slog.Any("error", err))
		http.Error(w, "Error getting recipes", http.StatusInternalServerError)
		return
	}

	// For now, simple pagination logic - we'll improve this later
	currentPage := 1
	totalPages := 1
	if len(recipes) > 0 {
		totalPages = (len(recipes) + 9) / 10 // 10 recipes per page
	}

	// Get the current year for footer
	currentYear := time.Now().Year()

	// Create template data
	data := TemplateData{
		Recipes:     recipes,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		CurrentYear: currentYear,
		SearchQuery: searchQuery,
		SortOption:  sortOption,
	}

	err = h.tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully rendered list of recipes")
}

// Basic handler for getting a single recipe
func (h *RecipeHandler) getRecipe(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling get recipe request")

	// Get recipe ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Looking for recipe with ID", slog.String("id", id))

	// Get recipe from store
	recipe, err := h.store.Get(id)
	if err != nil {
		h.logger.Error("Error getting recipe", slog.Any("error", err))
		http.Error(w, "Error getting recipe", http.StatusInternalServerError)
		return
	}

	// For demonstration purposes, we'll assume the user can edit their own recipes
	// In a real app, you'd check if the current user is the owner
	canEdit := true
	currentYear := time.Now().Year()

	// Render recipe using the new recipe-detail.html template
	data := TemplateData{
		Recipe:      &recipe,
		CanEdit:     canEdit,
		CurrentYear: currentYear,
	}

	// Execute template
	err = h.tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully rendered recipe", slog.String("title", recipe.Title))
}

// Show the create recipe form
func (h *RecipeHandler) createRecipeForm(w http.ResponseWriter, r *http.Request) {
	currentYear := time.Now().Year()

	data := TemplateData{
		Recipe:      &models.Recipe{}, // Empty recipe for the form
		CurrentYear: currentYear,
	}

	err := h.tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle the form submission
func (h *RecipeHandler) createRecipe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Error parsing form", slog.Any("error", err))
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	// Generate a unique ID (we'll improve this later)
	id := fmt.Sprintf("recipe-%d", time.Now().Unix())

	// Parse form values with error handling
	prepTimeStr := r.FormValue("prep_time")
	prepTimeInt, err := strconv.Atoi(prepTimeStr)
	if err != nil {
		h.logger.Error("Invalid prep time", slog.Any("error", err))
		http.Error(w, "Invalid prep time", http.StatusBadRequest)
		return
	}
	prepTime := time.Duration(prepTimeInt) * time.Minute

	cookTimeStr := r.FormValue("cook_time")
	cookTimeInt, err := strconv.Atoi(cookTimeStr)
	if err != nil {
		h.logger.Error("Invalid cook time", slog.Any("error", err))
		http.Error(w, "Invalid cook time", http.StatusBadRequest)
		return
	}
	cookTime := time.Duration(cookTimeInt) * time.Minute

	servings, err := strconv.Atoi(r.FormValue("servings"))
	if err != nil {
		h.logger.Error("Invalid servings", slog.Any("error", err))
		http.Error(w, "Invalid servings", http.StatusBadRequest)
		return
	}

	// Parse ingredients
	names := r.Form["ingredient_names[]"]
	amounts := r.Form["ingredient_amounts[]"]
	units := r.Form["ingredient_units[]"]

	ingredients := make([]models.Ingredient, len(names))
	for i := range names {
		if names[i] == "" {
			continue
		}

		amount, err := strconv.ParseFloat(amounts[i], 64)
		if err != nil {
			h.logger.Error("Invalid amount for ingredient", slog.Any("error", err))
			http.Error(w, "Invalid ingredient amount", http.StatusBadRequest)
			return
		}

		ingredients[i] = models.Ingredient{
			ID:       fmt.Sprintf("ing-%d", i),
			Name:     names[i],
			Amount:   amount,
			Unit:     units[i],
			Position: i,
		}
	}

	// Parse instructions
	instructionSteps := r.Form["instructions[]"]
	instructions := make([]models.Instruction, len(instructionSteps))
	for i, step := range instructionSteps {
		if step == "" {
			continue
		}
		instructions[i] = models.Instruction{
			ID:       fmt.Sprintf("step-%d", i),
			Step:     step,
			Position: i,
		}
	}

	recipe := models.Recipe{
		ID:           id,
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
		PrepTime:     prepTime,
		CookTime:     cookTime,
		Servings:     int32(servings),
		Ingredients:  ingredients,
		Instructions: instructions,
	}

	// Validate required fields
	if recipe.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Store the recipe
	if err := h.store.Create(recipe); err != nil {
		h.logger.Error("Error creating recipe", slog.Any("error", err))
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
		h.logger.Error("Error getting recipe to edit", slog.Any("error", err))
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	currentYear := time.Now().Year()

	data := TemplateData{
		Recipe:      &recipe,
		CurrentYear: currentYear,
	}

	err = h.tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Update a recipe
func (h *RecipeHandler) updateRecipe(w http.ResponseWriter, r *http.Request) {
	// Get recipe ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Error parsing form", slog.Any("error", err))
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	// Get existing recipe
	existingRecipe, err := h.store.Get(id)
	if err != nil {
		h.logger.Error("Error getting recipe to update", slog.Any("error", err))
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Parse form values
	prepTimeStr := r.FormValue("prep_time")
	prepTimeInt, err := strconv.Atoi(prepTimeStr)
	if err != nil {
		h.logger.Error("Invalid prep time", slog.Any("error", err))
		http.Error(w, "Invalid prep time", http.StatusBadRequest)
		return
	}
	prepTime := time.Duration(prepTimeInt) * time.Minute

	cookTimeStr := r.FormValue("cook_time")
	cookTimeInt, err := strconv.Atoi(cookTimeStr)
	if err != nil {
		h.logger.Error("Invalid cook time", slog.Any("error", err))
		http.Error(w, "Invalid cook time", http.StatusBadRequest)
		return
	}
	cookTime := time.Duration(cookTimeInt) * time.Minute

	servings, err := strconv.Atoi(r.FormValue("servings"))
	if err != nil {
		h.logger.Error("Invalid servings", slog.Any("error", err))
		http.Error(w, "Invalid servings", http.StatusBadRequest)
		return
	}

	// Parse ingredients
	names := r.Form["ingredient_names[]"]
	amounts := r.Form["ingredient_amounts[]"]
	units := r.Form["ingredient_units[]"]

	ingredients := make([]models.Ingredient, len(names))
	for i := range names {
		if names[i] == "" {
			continue
		}

		amount, err := strconv.ParseFloat(amounts[i], 64)
		if err != nil {
			h.logger.Error("Invalid amount for ingredient", slog.Any("error", err))
			http.Error(w, "Invalid ingredient amount", http.StatusBadRequest)
			return
		}

		ingredients[i] = models.Ingredient{
			ID:       fmt.Sprintf("ing-%d", i),
			Name:     names[i],
			Amount:   amount,
			Unit:     units[i],
			Position: i,
		}
	}

	// Parse instructions
	instructionSteps := r.Form["instructions[]"]
	instructions := make([]models.Instruction, len(instructionSteps))
	for i, step := range instructionSteps {
		if step == "" {
			continue
		}
		instructions[i] = models.Instruction{
			ID:       fmt.Sprintf("step-%d", i),
			Step:     step,
			Position: i,
		}
	}

	// Update recipe fields
	existingRecipe.Title = r.FormValue("title")
	existingRecipe.Description = r.FormValue("description")
	existingRecipe.PrepTime = prepTime
	existingRecipe.CookTime = cookTime
	existingRecipe.Servings = int32(servings)
	existingRecipe.Ingredients = ingredients
	existingRecipe.Instructions = instructions

	// Validate required fields
	if existingRecipe.Title == "" {
		data := TemplateData{
			Recipe: &existingRecipe,
			Error:  "Title is required",
		}
		h.tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	// Update the recipe
	if err := h.store.Update(existingRecipe); err != nil {
		h.logger.Error("Error updating recipe", slog.Any("error", err))
		http.Error(w, "Error saving recipe", http.StatusInternalServerError)
		return
	}

	// Redirect to the updated recipe
	http.Redirect(w, r, "/recipes/"+existingRecipe.ID, http.StatusSeeOther)
}

// Delete a recipe
func (h *RecipeHandler) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	// Get recipe ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete the recipe
	if err := h.store.Delete(id); err != nil {
		h.logger.Error("Error deleting recipe", slog.Any("error", err))
		http.Error(w, "Error deleting recipe", http.StatusInternalServerError)
		return
	}

	// Redirect to the recipes list
	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}
