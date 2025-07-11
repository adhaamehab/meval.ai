package sources

import "context"

// Source defines the interface for input/output data sources
type Source interface {
	// Read reads data from the source
	Read(ctx context.Context) ([]Record, error)
	// Write writes data to the source
	Write(ctx context.Context, records []Record) error
	// Close closes the source
	Close() error
}

// Record represents a single data record
type Record map[string]interface{}

// Factory creates sources based on configuration
type Factory interface {
	CreateSource(config map[string]interface{}, format string) (Source, error)
}