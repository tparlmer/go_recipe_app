// This file contains the entry point for the application

package main

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
