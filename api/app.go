package api

import (
	"io"

	"preselect/api/processor"
	"preselect/business"
	"preselect/data"
)

// Config holds application configuration options.
type Config struct {
	Root      string
	ExtMap    map[string]SourceFactory
	Processor processor.Processor
}

// App wires configuration with the business layer.
type App struct {
	cfg Config
}

// New creates a new App instance.
func New(cfg Config) *App {
	return &App{cfg: cfg}
}

// Run starts the application.
func (a *App) Run() error {
	extMap := make(map[string]SourceFactory)
	for k, v := range a.cfg.ExtMap {
		extMap[k] = v
	}
	if _, ok := extMap["txt"]; !ok {
		extMap["txt"] = func(r io.Reader) business.DataSource {
			return data.NewLoader(r, nil)
		}
	}

	root := a.cfg.Root
	if root == "" {
		root = "."
	}

	proc := a.cfg.Processor
	if proc == nil {
		proc = noopProcessor{}
	}

	return ScanDirectory(root, extMap, proc)
}

type noopProcessor struct{}

func (noopProcessor) Process(string) (bool, error) { return false, nil }
