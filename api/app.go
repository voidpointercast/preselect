package api

import (
	"preselect/business"
	"preselect/data"
	"strings"
)

// Config holds application configuration options.
type Config struct {
	// Placeholder for configuration fields
}

// App wires configuration with the business layer.
type App struct {
	scanner business.Scanner
	cfg     Config
}

// New creates a new App instance.
func New(cfg Config) *App {
	source := data.NewLoader(strings.NewReader(""), nil)
	return &App{
		scanner: business.NewScanner(source),
		cfg:     cfg,
	}
}

// Run starts the application.
func (a *App) Run() error {
	return a.scanner.Scan(nil)
}
