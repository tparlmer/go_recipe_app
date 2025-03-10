// This file contains the entry point for the application

package main

import (
	"go_recipe_app/internal/handlers"
	"html/template" // Go's built-in template package
	"log"           // For logging messages and errors
	"net/http"      // Go's web server package
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

	// Create new handler instance
	log.Println("Initializing recipe handler...")
	recipeHandler := handlers.NewRecipeHandler(tmpl)
	log.Println("Recipe handler initialized")

	// Register routes
	// http.HandleFunc("/", homeHandler)
	http.HandleFunc("/test", recipeHandler.TestRecipe)
	log.Println("Routes registered: / and test")

	log.Println("Server is running on http://localhost:8080")
	// This is a common design pattern in Go below, where something intended to run continuously is nested inside of log.Fatal() so that if it crashes it returns an error message
	log.Fatal(http.ListenAndServe(":8080", nil))
}
