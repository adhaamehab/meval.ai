package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ReaderInterface defines the interface for configuration readers
type ReaderInterface interface {
	Read(r io.Reader) (*Config, error)
	ReadFile(path string) (*Config, error)
}

// ValidatorInterface defines the interface for configuration validation
type ValidatorInterface interface {
	Validate(config *Config) error
}

// Reader implements configuration reading functionality
type Reader struct{}

// NewReader creates a new config reader
func NewReader() *Reader {
	return &Reader{}
}

// Read reads configuration from an io.Reader
func (r *Reader) Read(reader io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode yaml: %w", err)
	}

	return &config, nil
}

// ReadFile reads configuration from a file path
func (r *Reader) ReadFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	return r.Read(file)
}