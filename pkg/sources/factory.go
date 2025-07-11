package sources

import (
	"fmt"

	"github.com/adhaamehab/meval.ai/pkg/config"
)

// DefaultFactory implements the Factory interface
type DefaultFactory struct{}

// NewDefaultFactory creates a new source factory
func NewDefaultFactory() *DefaultFactory {
	return &DefaultFactory{}
}

// CreateSource creates a source based on format and configuration
func (f *DefaultFactory) CreateSource(cfg map[string]interface{}, format string, schema config.SchemaConfig) (Source, error) {
	switch format {
	case "json":
		return NewJSONSource(cfg, schema)
	case "csv":
		return nil, fmt.Errorf("CSV source not yet implemented")
	case "parquet":
		return nil, fmt.Errorf("Parquet source not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}