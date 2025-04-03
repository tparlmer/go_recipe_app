// This file contains the entry point for the application

package main

import (
	"fmt"
	"go_recipe_app/internal/config"
	"go_recipe_app/internal/handlers/recipe"
	"go_recipe_app/internal/storage/boltdb"
	"html/template" // Go's built-in template package
	"log"           // For logging messages and errors
	"net/http"      // Go's web server package
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure logging
	log.Printf("Starting application in %s mode", cfg.Env)
	log.Printf("Log directory: %s", cfg.LogDir)

	// Parse templates
	log.Println("Starting template parsing...")
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	log.Println("Templates parsed successfully")

	// Initialize store with configured path
	store, err := boltdb.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer store.Close()
	log.Println("BoltDB store initialized successfully")

	// Create new handler instance
	log.Println("Initializing recipe handler...")
	recipeHandler := recipe.New(tmpl, store)
	log.Println("Recipe handler initialized")

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server is running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, recipeHandler.Router))
}
