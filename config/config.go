package config

/*
ARCHITECTURAL PSEUDOCODE - ENHANCED CONFIGURATION

FUNCTION LoadConfig()
    Load configuration from environment variables
    Apply defaults for missing values

    Config structure:
    - Server settings (port, environment, timeouts)
    - Database settings (paths for recipes and users)
    - Logging settings (level, format, paths)
    - Auth settings (JWT secret, token expiration)
    - Security settings (allowed origins, CSRF)

    Validate configuration values
    Return complete configuration object
END FUNCTION

Auth Settings Required:
- JWT_SECRET: Secret key for signing JWT tokens (generate secure random string)
- TOKEN_EXPIRATION: How long tokens are valid for (e.g., 24h)
- COOKIE_SECURE: Whether cookies should be secure-only (true in production)
- COOKIE_DOMAIN: Domain for cookies
- PASSWORD_MIN_LENGTH: Minimum password length requirement

Database Settings Enhancements:
- Separate paths for recipe and user databases
- Backup directory configuration
- Connection pool settings
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	// Server settings
	Port    int
	Env     string
	BaseURL string

	// Database settings
	DataDir   string
	DBTimeout time.Duration

	// Logging settings
	LogDir    string
	LogLevel  string
	LogFormat string // json or text
	LogPath   string

	// Auth settings
	JWTSecret       string
	TokenExpiration time.Duration
	CookieSecure    bool
	CookieDomain    string
	PasswordMinLen  int

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

	tokenExpiration, err := time.ParseDuration(getEnvWithDefault("RECIPE_APP_TOKEN_EXPIRATION", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid token expiration: %v", err)
	}

	passwordMinLen, err := strconv.Atoi(getEnvWithDefault("RECIPE_APP_PASSWORD_MIN_LENGTH", "8"))
	if err != nil {
		return nil, fmt.Errorf("invalid password minimum length: %v", err)
	}

	cookieSecure := getEnvWithDefault("RECIPE_APP_COOKIE_SECURE", "false") == "true"

	dataDir := getEnvWithDefault("RECIPE_APP_DATA_DIR", "data")

	config := &Config{
		// Server settings
		Port:    port,
		Env:     getEnvWithDefault("RECIPE_APP_ENV", "development"),
		BaseURL: getEnvWithDefault("RECIPE_APP_BASE_URL", "http://localhost:8080"),

		// Database settings
		DataDir:   dataDir,
		DBTimeout: time.Second,

		// Logging settings
		LogDir:    getEnvWithDefault("RECIPE_APP_LOG_DIR", "logs"),
		LogLevel:  getEnvWithDefault("RECIPE_APP_LOG_LEVEL", "info"),
		LogFormat: getEnvWithDefault("RECIPE_APP_LOG_FORMAT", "text"),
		LogPath: filepath.Join(
			getEnvWithDefault("RECIPE_APP_LOG_DIR", "logs"),
			"app.log",
		),

		// Auth settings
		JWTSecret:       getEnvWithDefault("RECIPE_APP_JWT_SECRET", "your-secret-key-change-in-production"),
		TokenExpiration: tokenExpiration,
		CookieSecure:    cookieSecure,
		CookieDomain:    getEnvWithDefault("RECIPE_APP_COOKIE_DOMAIN", ""),
		PasswordMinLen:  passwordMinLen,

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

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(c.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Warn if using default JWT secret in production
	if c.Env == "production" && c.JWTSecret == "your-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret in production environment!")
	}

	return nil
}
