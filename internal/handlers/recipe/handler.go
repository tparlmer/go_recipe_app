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
	Template string
	Data     interface{}
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
	h.Router.HandleFunc("/recipes", h.createRecipe).Methods("POST")        // Handle form submission
	h.Router.HandleFunc("/recipes/{id}", h.getRecipe).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}/edit", h.editRecipeForm).Methods("GET")
	h.Router.HandleFunc("/recipes/{id}", h.updateRecipe).Methods("PUT")
	h.Router.HandleFunc("/recipes/{id}", h.deleteRecipe).Methods("DELETE")
}

// Basic handler for listing recipes
func (h *RecipeHandler) listRecipes(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling list recipes request")

	recipes, err := h.store.List()
	if err != nil {
		h.logger.Error("Error listing recipes", slog.Any("error", err))
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
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully rendered list of test recipes")
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

	// Render recipe
	data := TemplateData{
		Template: "view",
		Data:     recipe,
	}

	// Execute template
	err = h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully rendered recipe", slog.String("title", recipe.Title))
}

// Show the create recipe form
func (h *RecipeHandler) createRecipeForm(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Template: "create",
		Data:     nil,
	}

	err := h.tmpl.ExecuteTemplate(w, "layout.html", data)
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
	prepTime, err := time.ParseDuration(r.FormValue("prep_time") + "m")
	if err != nil {
		h.logger.Error("Invalid prep time", slog.Any("error", err))
		http.Error(w, "Invalid prep time", http.StatusBadRequest)
		return
	}

	cookTime, err := time.ParseDuration(r.FormValue("cook_time") + "m")
	if err != nil {
		h.logger.Error("Invalid cook time", slog.Any("error", err))
		http.Error(w, "Invalid cook time", http.StatusBadRequest)
		return
	}

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

	data := TemplateData{
		Template: "edit",
		Data:     recipe,
	}

	err = h.tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		h.logger.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle the update
func (h *RecipeHandler) updateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Update request - ID format check", slog.String("id", id))

	// Just check if recipe exists
	if _, err := h.store.Get(id); err != nil {
		h.logger.Error("Error getting existing recipe", slog.Any("error", err))
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Error parsing form", slog.Any("error", err))
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	// Debug logging for all form values
	h.logger.Info("All form values", slog.Any("form", r.Form))
	h.logger.Info("Form method", slog.String("method", r.Method))
	h.logger.Info("Content-Type", slog.String("content_type", r.Header.Get("Content-Type")))

	// Check if we're getting the value from Form vs PostForm
	h.logger.Info("prep_time from FormValue", slog.String("prep_time", r.FormValue("prep_time")))
	h.logger.Info("prep_time from Form", slog.String("prep_time", r.Form.Get("prep_time")))
	h.logger.Info("prep_time from PostForm", slog.String("prep_time", r.PostForm.Get("prep_time")))

	// Parse form values with error handling
	prepTimeStr := r.FormValue("prep_time")
	if prepTimeStr == "" {
		h.logger.Error("Empty prep time received")
		http.Error(w, "Prep time is required", http.StatusBadRequest)
		return
	}
	prepTime := time.Duration(mustParseFloat(prepTimeStr)) * time.Minute

	cookTimeStr := r.FormValue("cook_time")
	if cookTimeStr == "" {
		h.logger.Error("Empty cook time received")
		http.Error(w, "Cook time is required", http.StatusBadRequest)
		return
	}
	cookTime := time.Duration(mustParseFloat(cookTimeStr)) * time.Minute

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

	h.logger.Info("Method", slog.String("method", r.Method))
	h.logger.Info("Content-Type", slog.String("content_type", r.Header.Get("Content-Type")))
	h.logger.Info("Raw prep_time value", slog.String("prep_time", r.FormValue("prep_time")))
	h.logger.Info("Raw cook_time value", slog.String("cook_time", r.FormValue("cook_time")))

	if err := h.store.Update(recipe); err != nil {
		h.logger.Error("Error updating recipe", slog.Any("error", err))
		http.Error(w, "Error updating recipe", http.StatusInternalServerError)
		return
	}

	// Instead of redirecting, send a success response
	w.WriteHeader(http.StatusOK)
	// Optionally send a success message
	w.Write([]byte("Recipe updated successfully"))
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// Delete recipe handler
func (h *RecipeHandler) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Info("Attempting to delete recipe", slog.String("id", id))

	// Check if recipe exists
	_, err := h.store.Get(id)
	if err != nil {
		h.logger.Error("Recipe not found for deletion", slog.Any("error", err))
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Delete the recipe
	if err := h.store.Delete(id); err != nil {
		h.logger.Error("Error deleting recipe", slog.Any("error", err))
		http.Error(w, "Error deleting recipe", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully deleted recipe", slog.String("id", id))
	http.Redirect(w, r, "/recipes", http.StatusSeeOther)
}
