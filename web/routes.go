package web

/*
ARCHITECTURAL PSEUDOCODE - ROUTES SETUP

FUNCTION SetupRoutes(authService, recipeService, logger)
    Create main router

    // Register static file handlers
    router.PathPrefix("/static/").Handler(...)

    // Register auth routes
    authHandler := auth.NewHandler(authService, templates, logger)
    authHandler.RegisterRoutes(router)

    // Register public recipe routes
    recipeHandler := recipes.NewHandler(recipeService, templates, logger)
    recipeHandler.RegisterPublicRoutes(router)

    // Register protected recipe routes
    protectedRouter := router.PathPrefix("/").Subrouter()
    protectedRouter.Use(authService.AuthMiddleware)
    recipeHandler.RegisterProtectedRoutes(protectedRouter)

    // 404 handler
    router.NotFoundHandler = customNotFoundHandler()

    Return router
END FUNCTION

FUNCTION Route Mapping
    AUTH ROUTES:
    - GET /login        -> Login form
    - POST /login       -> Process login
    - GET /register     -> Registration form
    - POST /register    -> Process registration
    - GET /logout       -> Process logout

    PUBLIC RECIPE ROUTES:
    - GET /             -> Home page
    - GET /recipes      -> List public recipes
    - GET /recipes/{id} -> View public recipe

    PROTECTED RECIPE ROUTES:
    - GET /my-recipes           -> List user's recipes
    - GET /recipes/new          -> Create recipe form
    - POST /recipes/create      -> Process recipe creation
    - GET /recipes/{id}/edit    -> Edit recipe form
    - POST /recipes/{id}/update -> Process recipe update
    - POST /recipes/{id}/delete -> Process recipe deletion
END FUNCTION
*/

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go_recipe_app/auth"
	"go_recipe_app/recipes"
	"go_recipe_app/recipes/db"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *mux.Router, authService auth.AuthService, recipeService recipes.RecipeService, logger log.Logger) {
	// Define template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"CurrentYear": func() int {
			return time.Now().Year()
		},
	}

	// Parse templates with the function map
	logger.Log("msg", "Parsing templates from web/templates/*.html")
	tmpl := template.New("").Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseGlob("web/templates/*.html"))

	// Log the templates that were found
	for _, t := range tmpl.Templates() {
		logger.Log("msg", "Template loaded", "name", t.Name())
	}

	// Add a test route for debugging
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("msg", "Attempting to render test template")
		err := tmpl.ExecuteTemplate(w, "test.html", nil)
		if err != nil {
			logger.Log("msg", "Error rendering test template", "err", err)
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	})

	// Add a direct debug route that bypasses templates
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Debug Page</h1><p>This bypasses template rendering</p></body></html>"))
	})

	// Add debug route for layout template
	router.HandleFunc("/debug-layout", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("msg", "Attempting to render layout template")

		// Create some sample recipe data
		sampleRecipes := []db.Recipe{
			{
				ID:          "1",
				Title:       "Test Recipe",
				Description: "A sample recipe for testing",
				PrepTime:    time.Duration(15) * time.Minute,
				CookTime:    time.Duration(30) * time.Minute,
				Servings:    4,
			},
		}

		// Create template data
		data := map[string]interface{}{
			"Recipes":         sampleRecipes,
			"IsAuthenticated": true,
			"CurrentYear":     time.Now().Year(),
			"CurrentTemplate": "home.html",
		}

		// Try rendering the layout with content
		err := tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			logger.Log("msg", "Error rendering layout template", "err", err)
			http.Error(w, fmt.Sprintf("Layout template error: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Add debug route specifically for home template
	router.HandleFunc("/debug-home", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("msg", "Attempting to render home template")

		// Create some sample recipe data
		sampleRecipes := []db.Recipe{
			{
				ID:          "1",
				Title:       "Test Recipe",
				Description: "A sample recipe for testing",
				PrepTime:    time.Duration(15) * time.Minute,
				CookTime:    time.Duration(30) * time.Minute,
				Servings:    4,
			},
		}

		// Create template data
		data := map[string]interface{}{
			"Recipes":         sampleRecipes,
			"IsAuthenticated": true,
		}

		// Try rendering just the home template
		err := tmpl.ExecuteTemplate(w, "home.html", data)
		if err != nil {
			logger.Log("msg", "Error rendering home template", "err", err)
			http.Error(w, fmt.Sprintf("Home template error: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Add debug routes that output only basic HTML
	router.HandleFunc("/debug-simple", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Simple Debug Page</h1><p>This page bypasses the template system entirely</p></body></html>"))
	})

	// Add a debug route for testing login form
	router.HandleFunc("/debug-login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
		<html>
		<body>
			<h1>Debug Login Form</h1>
			<form method="POST" action="/login">
				<div>
					<label for="username">Username</label>
					<input type="text" id="username" name="username" required>
				</div>
				<div>
					<label for="password">Password</label>
					<input type="password" id="password" name="password" required>
				</div>
				<div>
					<button type="submit">Login</button>
				</div>
			</form>
		</body>
		</html>
		`
		w.Write([]byte(html))
	})

	// Add route to test content templates directly
	router.HandleFunc("/debug-content/{template}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		templateName := vars["template"]
		logger.Log("msg", "Attempting to render content template directly", "template", templateName)

		data := map[string]interface{}{
			"CurrentYear":     time.Now().Year(),
			"Debug":           true,
			"CurrentURL":      r.URL.Path,
			"CurrentTemplate": templateName,
		}

		// Try rendering the template
		err := tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			logger.Log("msg", "Error rendering content template directly", "template", templateName, "err", err)
			http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Create handlers
	authHandler := auth.NewHandler(authService, tmpl, logger)
	recipeHandler := recipes.NewHandler(recipeService, tmpl, logger)

	// Setup auth routes
	router.HandleFunc("/login", authHandler.LoginForm()).Methods("GET")
	router.HandleFunc("/login", authHandler.Login()).Methods("POST")
	router.HandleFunc("/register", authHandler.RegisterForm()).Methods("GET")
	router.HandleFunc("/register", authHandler.Register()).Methods("POST")
	router.HandleFunc("/logout", authHandler.Logout())

	// Setup public recipe routes
	router.HandleFunc("/", recipeHandler.Home())
	router.HandleFunc("/recipes", recipeHandler.ListPublicRecipes()).Methods("GET")
	router.HandleFunc("/recipes/{id}", recipeHandler.ViewRecipe()).Methods("GET")

	// Setup protected recipe routes with auth middleware
	// Create a subrouter with specific paths for protected routes
	protected := router.PathPrefix("").Subrouter()
	protected.Use(auth.AuthMiddleware(authService, logger))

	// Add specific routes to the protected subrouter
	protected.HandleFunc("/my-recipes", recipeHandler.ListUserRecipes()).Methods("GET")
	protected.HandleFunc("/recipes/new", recipeHandler.CreateRecipeForm()).Methods("GET")
	protected.HandleFunc("/recipes/create", recipeHandler.CreateRecipe()).Methods("POST")
	protected.HandleFunc("/recipes/{id}/edit", recipeHandler.EditRecipeForm()).Methods("GET")
	protected.HandleFunc("/recipes/{id}/update", recipeHandler.UpdateRecipe()).Methods("POST")
	protected.HandleFunc("/recipes/{id}/delete", recipeHandler.DeleteRecipe()).Methods("POST")

	// Setup 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		tmpl.ExecuteTemplate(w, "404.html", nil)
	})
}
