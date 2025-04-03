package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type LogConfig struct {
	Level   string
	Format  string // json or text
	LogPath string
}

func NewLogger(cfg LogConfig) (*slog.Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0755); err != nil {
		return nil, err
	}

	// Open log file
	f, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create multi-writer for both file and stdout
	mw := io.MultiWriter(os.Stdout, f)

	// Set log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Configure handler
	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(mw, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(mw, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler), nil
}
