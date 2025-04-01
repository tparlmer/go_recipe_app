// This file contains the entry point for the application

package main

import (
	"go_recipe_app/internal/handlers/recipe"
	"go_recipe_app/internal/storage/boltdb"
	"html/template" // Go's built-in template package
	"log"           // For logging messages and errors
	"net/http"      // Go's web server package
	"os"
)

func main() {
	// Parse templates
	log.Println("Starting template parsing...")
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	log.Println("Templates parsed successfully")

	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/recipes.db" // fallback to current behavior
	}

	// Initialize store with configurable path
	store, err := boltdb.New(dbPath)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer store.Close()
	log.Println("BoltDB store initialized successfully")

	// Create new handler instance
	log.Println("Initializing recipe handler...")
	recipeHandler := recipe.New(tmpl, store)
	log.Println("Recipe handler initialized")

	log.Println("Server is running on http://localhost:8080")
	// This is a common design pattern in Go below, where something intended to run continuously is nested inside of log.Fatal() so that if it crashes it returns an error message
	log.Fatal(http.ListenAndServe(":8080", recipeHandler.Router))
}
