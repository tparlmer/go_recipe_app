package main

/*
ARCHITECTURAL PSEUDOCODE - MAIN ENTRY POINT

FUNCTION main()
    Initialize logging system
    Set up configuration
    Prepare data directory

    Initialize Auth Service:
      - Create auth repository with data directory
      - Set up auth logger
      - Configure with JWT secret
      - Register auth routes

    Initialize Recipe Service:
      - Create recipe repository with data directory
      - Set up recipe logger
      - Configure service with repository
      - Register recipe routes

    Create HTTP router with these routes:
      - Public routes (no auth required)
      - Protected routes (with auth middleware)

    Start HTTP server
END FUNCTION
*/

import (
	"fmt"
	"go_recipe_app/auth"
	"go_recipe_app/recipes"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Initialize logging service
	appLogger := setupLogger("logs/recipe-app.log")

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		appLogger.Log("msg", "Error loading .env file", "err", err)
		os.Exit(1)
	}

	// Initialize auth service
	var authSvc auth.AuthService
	authSvc = setupAuthService("data", appLogger)

	// Initialize recipe service
	var recipeSvc recipes.RecipeService
	recipeSvc = setupRecipeService("data", appLogger)

	// Initialize routes
	router := mux.NewRouter()
	handler := setupRoutes(router, authSvc, recipeSvc, appLogger)

	// Start HTTP server
	port := os.Getenv("RECIPE_APP_PORT")
	if port == "" {
		// Provide default if not found
		port = "8080"
	}
	addr := ":" + port
	err := http.ListenAndServe(addr, handler)

	// This code only executes if the server fails to start
	if err != nil {
		appLogger.Log("msg", "Failed to start HTTP server", "err", err)
		os.Exit(1)
	}
}

// setupLogger creates a logger that writes to a file and rotates the log file using the Lumberjack golang library
func setupLogger(filename string) log.Logger {
	logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1024, // Max size in MB before log is rotated
		MaxBackups: 3,    // Max number of old log files to keep
		MaxAge:     28,   // Max age in days to keep a log file
		Compress:   true, // Whether to compress log files using gzip
	}

	var cumulativeLogger log.Logger
	cumulativeLogger = log.NewLogfmtLogger(log.NewSyncWriter(logger))
	cumulativeLogger = log.With(cumulativeLogger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	return cumulativeLogger
}

// setupAuthService initializes the auth service in main
func setupAuthService(dataDir string, mainLogger log.Logger) auth.AuthService {
	authLogger := setupLogger("logs/auth.log")

	var authSvc auth.AuthService
	authSvc, err := auth.NewAuthService(os.Getenv("JWT_SECRET_KEY"), dataDir, mainLogger, authLogger)
	if err != nil {
		panic(err)
	}
	authSvc = auth.NewLoggingMiddleware(authLogger, authSvc)

	return authSvc
}

// setupRecipeService initializes the recipe service in main
func setupRecipeService(dataDir string, mainLogger log.Logger) recipes.RecipeService {
	recipesLogger := setupLogger("logs/recipes.log")

	var recipeSvc recipes.RecipeService
	recipeSvc, err := recipes.NewRecipeService(dataDir, mainLogger, recipesLogger)
	if err != nil {
		panic(err)
	}
	// TODO: Add logging middleware

	return recipeSvc
}

// Sets up all routes for the application
func setupRoutes(router *mux.Router, authService auth.AuthService, recipeService recipes.RecipeService, logger log.Logger) http.Handler {
	// Add logic
	return router
}

/* Original code:
import (
	"fmt"
	"go_recipe_app/internal/config"
	"go_recipe_app/internal/handlers/recipe"
	"go_recipe_app/internal/logging"
	"go_recipe_app/internal/storage/boltdb"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize new logger
	logger, err := logging.NewLogger(logging.LogConfig{
		Level:   cfg.LogLevel,
		Format:  cfg.LogFormat,
		LogPath: cfg.LogPath,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Parse templates
	logger.Info("parsing templates")
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		logger.Error("failed to parse templates", "error", err)
		return
	}
	logger.Info("templates parsed successfully")

	// Initialize store
	store, err := boltdb.New(cfg.DBPath)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		return
	}
	defer store.Close()
	logger.Info("database initialized", "path", cfg.DBPath)

	// Create handler
	logger.Info("initializing recipe handler")
	recipeHandler := recipe.New(tmpl, store, logger)
	logger.Info("recipe handler initialized")

	addr := fmt.Sprintf(":%d", cfg.Port)
	if cfg.Env == "development" || cfg.Env == "local" {
		logger.Info("starting development server",
			"url", fmt.Sprintf("http://localhost%s", addr),
			"env", cfg.Env,
		)
	} else {
		logger.Info("starting production server",
			"port", cfg.Port,
			"env", cfg.Env,
		)
	}

	if err := http.ListenAndServe(addr, recipeHandler.Router); err != nil {
		logger.Error("server failed", "error", err)
	}
}
*/
