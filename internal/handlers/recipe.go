// OLD TEST CODE - CURRENTLY DEFUNCT

package handlers

import (
	"html/template"
	"net/http"
	"go_recipe_app/internal/models"
	"time"
	"log"
)

// This is like a class in other languages
type RecipeHandler struct {
	tmpl *template.Template // Store templates
	// Later we will add more stored resources
}

// Constructor function - Go's pattern for creating new instances
func NewRecipeHandler(tmpl *template.Template) *RecipeHandler {
	return &RecipeHandler{
		tmpl: tmpl,
	}
}

// Helper method to create a slice of test data
func createTestRecipe() models.Recipe {
	return models.Recipe{
		ID:          "test-recipe-1",
		Title:       "Test Recipe",
		Description: "A simple test recipe",
		PrepTime:    15 * time.Minute,
		CookTime:    30 * time.Minute,
		Servings:    2,
		Ingredients: []models.Ingredient{
			{
				Quantity: 2.5,
				Unit:     "cups",
				Name:     "flour",
			},
		},
		Instructions: []models.Instruction{
			{
				Step:             1,
				StepInstructions: "Mix ingredients",
			},
		},
	}
}

// Method with receiver (r *RecipeHandler)
func (r *RecipeHandler) TestRecipe(w http.ResponseWriter, req *http.Request) {
	log.Println("=== TestRecipe handler called ===")

	recipe := createTestRecipe()
	log.Printf("Test recipe created: %+v", recipe)

	log.Println("attempting to execute layout.html")
	err := r.tmpl.ExecuteTemplate(w, "layout.html", recipe)
	if err != nil {
		log.Printf("ERROR executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Template executed successfully")
}