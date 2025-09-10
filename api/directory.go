package api

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"preselect/api/processor"
	"preselect/business"
)

// SourceFactory creates a business.DataSource from an io.Reader. The factory is
// used to wire file extensions to concrete data source implementations.
type SourceFactory func(io.Reader) business.DataSource

// ScanDirectory walks the provided directory recursively and, for every file
// whose extension matches an entry in extMap, creates a Scanner. Each Scanner
// runs concurrently and pushes entries to a shared channel that the supplied
// Processor consumes.
func ScanDirectory(root string, extMap map[string]SourceFactory, proc processor.Processor) error {
	entries := make(chan string)
	var scanners sync.WaitGroup

	// Start the processor consuming entries.
	var procErr error
	done := make(chan struct{})
	go func() {
		defer close(done)
		for e := range entries {
			if _, err := proc.Process(e); err != nil && procErr == nil {
				procErr = err
			}
		}
	}()

	// Discover files and launch scanners.
	walkErr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(d.Name())), ".")
		factory, ok := extMap[ext]
		if !ok {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		source := factory(f)
		scanner := business.NewScanner(source)

		scanners.Add(1)
		go func(file *os.File) {
			defer scanners.Done()
			defer file.Close()
			_ = scanner.ScanTo(entries)
		}(f)
		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	scanners.Wait()
	close(entries)
	<-done
	return procErr
}
