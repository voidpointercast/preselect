package business

import (
	"encoding/csv"
	"os"

	"preselect/api/processor"
)

// CSVProcessor implements processor.Processor by storing kept entries in a CSV file.
type CSVProcessor struct {
	cfg    processor.Config
	writer *csv.Writer
	file   *os.File
}

var _ processor.Processor = (*CSVProcessor)(nil)

// NewCSVProcessor creates a CSVProcessor writing kept entries to the given file path.
func NewCSVProcessor(cfg processor.Config, filePath string) (*CSVProcessor, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &CSVProcessor{
		cfg:    cfg,
		writer: csv.NewWriter(f),
		file:   f,
	}, nil
}

// Process evaluates an entry against configured keywords using the provided
// similarity metric. Entries meeting or exceeding the threshold are written to
// the CSV file. It returns true if the entry was kept.
func (p *CSVProcessor) Process(entry string) (bool, error) {
	for _, k := range p.cfg.Keywords {
		if p.cfg.Metric(entry, k) >= p.cfg.Threshold {
			if err := p.writer.Write([]string{entry}); err != nil {
				return false, err
			}
			p.writer.Flush()
			if err := p.writer.Error(); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

// Close releases the underlying file resources.
func (p *CSVProcessor) Close() error {
	p.writer.Flush()
	if err := p.writer.Error(); err != nil {
		_ = p.file.Close()
		return err
	}
	return p.file.Close()
}
