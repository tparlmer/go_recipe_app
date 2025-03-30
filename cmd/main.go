// This file contains the entry point for the application

package main

import (
	"go_recipe_app/internal/handlers/recipe"
	"go_recipe_app/internal/models"
	"go_recipe_app/internal/storage/memory"
	"html/template" // Go's built-in template package
	"log"           // For logging messages and errors
	"net/http"      // Go's web server package
	"time"
)

var tmpl *template.Template

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	err := tmpl.ExecuteTemplate(w, "layout.html", nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// Currently setup for Recipe testing
func main() {
	// Parse templates
	log.Println("Starting template parsing...")
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	log.Println("Templates parsed successfully")

	// Create memory store
	log.Println("Initializing memory store...")
	store := memory.New()

	// Add test recipes
	log.Println("Adding test recipes to store...")
	testRecipes := []models.Recipe{
		{
			ID:          "1",
			Title:       "Refried Beans",
			Description: "Canned beans cooked down with spices",
			PrepTime:    5 * time.Minute,
			CookTime:    10 * time.Minute,
			Servings:    4,
		},
		{
			ID:          "2",
			Title:       "Falafel",
			Description: "Deep fried falafel balls",
			PrepTime:    10 * time.Minute,
			CookTime:    10 * time.Minute,
			Servings:    4,
		},
	}

	// Store test recipes
	for _, recipe := range testRecipes {
		if err := store.Create(recipe); err != nil {
			log.Printf("Error creating test recipe: %v", err)
		}
	}
	log.Println("Test recipes added successfully")

	// Create new handler instance
	log.Println("Initializing recipe handler...")
	recipeHandler := recipe.New(tmpl, store)
	log.Println("Recipe handler initialized")

	log.Println("Server is running on http://localhost:8080")
	// This is a common design pattern in Go below, where something intended to run continuously is nested inside of log.Fatal() so that if it crashes it returns an error message
	log.Fatal(http.ListenAndServe(":8080", recipeHandler.Router))
}
