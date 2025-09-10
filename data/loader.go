package data

import (
	"io"

	"preselect/business"
)

// Loader retrieves file data entry by entry.
type Loader struct{}

// NewLoader creates a new Loader instance.
func NewLoader() Loader {
	return Loader{}
}

// Next returns the next entry from the source.
func (l Loader) Next() (business.Entry, error) {
	// Placeholder for loader logic
	return business.Entry{}, io.EOF
}
