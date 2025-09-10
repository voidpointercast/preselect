package processor

// SimilarityMetric defines a function to measure the similarity between two strings.
type SimilarityMetric func(a, b string) float64

// Config contains options for setting up a Processor.
type Config struct {
	Metric    SimilarityMetric
	Keywords  []string
	Threshold float64
}

// Processor evaluates entries using a similarity metric and stores those
// that meet a configured threshold. Implementations determine how kept
// entries are stored.
type Processor interface {
	// Process evaluates the entry. It returns true if the entry was kept
	// and stored, along with any error encountered during processing.
	Process(entry string) (bool, error)
}
