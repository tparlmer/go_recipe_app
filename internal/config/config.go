package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server settings
	Port    int
	Env     string
	BaseURL string

	// Database settings
	DBPath    string
	DBTimeout time.Duration

	// Logging settings
	LogDir    string
	LogLevel  string
	LogFormat string // json or text

	// Security settings
	AllowedOrigins []string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnvWithDefault("RECIPE_APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	readTimeout, err := time.ParseDuration(getEnvWithDefault("RECIPE_APP_READ_TIMEOUT", "15s"))
	if err != nil {
		return nil, fmt.Errorf("invalid read timeout: %v", err)
	}

	writeTimeout, err := time.ParseDuration(getEnvWithDefault("RECIPE_APP_WRITE_TIMEOUT", "15s"))
	if err != nil {
		return nil, fmt.Errorf("invalid write timeout: %v", err)
	}

	config := &Config{
		// Server settings
		Port:    port,
		Env:     getEnvWithDefault("RECIPE_APP_ENV", "development"),
		BaseURL: getEnvWithDefault("RECIPE_APP_BASE_URL", "http://localhost:8080"),

		// Database settings
		DBPath:    getEnvWithDefault("RECIPE_APP_DB_PATH", "data/recipes.db"),
		DBTimeout: time.Second,

		// Logging settings
		LogDir:    getEnvWithDefault("RECIPE_APP_LOG_DIR", "logs"),
		LogLevel:  getEnvWithDefault("RECIPE_APP_LOG_LEVEL", "info"),
		LogFormat: getEnvWithDefault("RECIPE_APP_LOG_FORMAT", "text"),

		// Security settings
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		AllowedOrigins: []string{
			getEnvWithDefault("RECIPE_APP_ALLOWED_ORIGIN", "*"),
		},
	}

	return config, config.validate()
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(c.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	return nil
}
