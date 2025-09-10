package business

// Entry represents a single piece of data from a source with its origin path.
type Entry struct {
	Value string
	Path  []string
}

// DataSource abstracts the retrieval of sequential entries.
type DataSource interface {
	Next() (Entry, error)
}
