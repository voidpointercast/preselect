package business

import "io"

// Scanner coordinates similarity checks on data provided by a DataSource.
type Scanner struct {
	source DataSource
}

// NewScanner creates a new Scanner instance.
func NewScanner(source DataSource) Scanner {
	return Scanner{source: source}
}

// Scan searches for keywords in the provided source.
func (s Scanner) Scan(keywords []string) error {
	for {
		entry, err := s.source.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_ = entry // placeholder for processing logic
	}
	_ = keywords // reference to avoid unused warning
	return nil
}

// ScanTo reads entries from the underlying DataSource and pushes their values
// to the provided channel until the source is exhausted. The channel is not
// closed by this method.
func (s Scanner) ScanTo(ch chan<- string) error {
	for {
		entry, err := s.source.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		ch <- entry.Value
	}
}
