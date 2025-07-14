package sources

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/adhaamehab/meval.ai/pkg/config"
)

// JSONSource implements Source interface for JSON files
type JSONSource struct {
	path       string
	mode       string // "array" or "lines"
	schema     config.SchemaConfig
	isWritable bool
	writer     io.WriteCloser
}

// NewJSONSource creates a new JSON source
func NewJSONSource(cfg map[string]interface{}, schema config.SchemaConfig) (*JSONSource, error) {
	path, ok := cfg["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path is required for JSON source")
	}

	mode, _ := cfg["mode"].(string)
	if mode == "" {
		mode = "array" // default to JSON array
	}

	if mode != "array" && mode != "lines" {
		return nil, fmt.Errorf("unsupported mode: %s (must be 'array' or 'lines')", mode)
	}

	return &JSONSource{
		path:   path,
		mode:   mode,
		schema: schema,
	}, nil
}

// Read reads records from JSON files
func (j *JSONSource) Read(ctx context.Context) ([]Record, error) {
	files, err := j.findFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found matching pattern: %s", j.path)
	}

	var allRecords []Record

	for _, file := range files {
		select {
		case <-ctx.Done():
			return allRecords, ctx.Err()
		default:
			records, err := j.readFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", file, err)
			}
			allRecords = append(allRecords, records...)
		}
	}

	return allRecords, nil
}

// Write writes records to a JSON file
func (j *JSONSource) Write(ctx context.Context, records []Record) error {
	if j.writer == nil {
		// Ensure directory exists
		dir := filepath.Dir(j.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		file, err := os.Create(j.path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		j.writer = file
		j.isWritable = true

		if j.mode == "array" {
			// Write opening bracket for array mode
			if _, err := j.writer.Write([]byte("[\n")); err != nil {
				return err
			}
		}
	}

	encoder := json.NewEncoder(j.writer)
	encoder.SetIndent("", "  ")

	for i, record := range records {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if j.mode == "array" && i > 0 {
				if _, err := j.writer.Write([]byte(",\n")); err != nil {
					return err
				}
			}

			// Validate record against schema
			if err := j.validateRecord(record); err != nil {
				return fmt.Errorf("record validation failed: %w", err)
			}

			if err := encoder.Encode(record); err != nil {
				return fmt.Errorf("failed to encode record: %w", err)
			}
		}
	}

	return nil
}

// Close closes the source
func (j *JSONSource) Close() error {
	if j.writer != nil {
		if j.mode == "array" {
			// Write closing bracket for array mode
			if _, err := j.writer.Write([]byte("\n]")); err != nil {
				return err
			}
		}
		return j.writer.Close()
	}
	return nil
}

// findFiles finds all files matching the path pattern
func (j *JSONSource) findFiles() ([]string, error) {
	// Check if path contains wildcards
	if filepath.IsAbs(j.path) && !containsWildcard(j.path) {
		// Direct file path
		if _, err := os.Stat(j.path); err != nil {
			return nil, err
		}
		return []string{j.path}, nil
	}

	// Use glob for wildcard patterns
	matches, err := filepath.Glob(j.path)
	if err != nil {
		return nil, err
	}

	// Filter only files (not directories)
	var files []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

// readFile reads records from a single JSON file
func (j *JSONSource) readFile(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if j.mode == "array" {
		return j.readJSONArray(file)
	}
	return j.readJSONLines(file)
}

// readJSONArray reads a JSON array file
func (j *JSONSource) readJSONArray(reader io.Reader) ([]Record, error) {
	var rawRecords []json.RawMessage
	decoder := json.NewDecoder(reader)
	
	if err := decoder.Decode(&rawRecords); err != nil {
		return nil, fmt.Errorf("failed to decode JSON array: %w", err)
	}

	records := make([]Record, 0, len(rawRecords))
	for i, raw := range rawRecords {
		var record Record
		if err := json.Unmarshal(raw, &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal record %d: %w", i, err)
		}

		if err := j.validateRecord(record); err != nil {
			return nil, fmt.Errorf("record %d validation failed: %w", i, err)
		}

		records = append(records, record)
	}

	return records, nil
}

// readJSONLines reads a JSON lines file (one JSON object per line)
func (j *JSONSource) readJSONLines(reader io.Reader) ([]Record, error) {
	var records []Record
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		var record Record
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal line %d: %w", lineNum, err)
		}

		if err := j.validateRecord(record); err != nil {
			return nil, fmt.Errorf("line %d validation failed: %w", lineNum, err)
		}

		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return records, nil
}

// validateRecord validates a record against the schema
func (j *JSONSource) validateRecord(record Record) error {
	for _, field := range j.schema.Fields {
		value, exists := record[field.Name]
		if !exists {
			return fmt.Errorf("missing required field: %s", field.Name)
		}

		if err := validateFieldType(value, field.Type); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	return nil
}

// validateFieldType validates that a value matches the expected type
func validateFieldType(value interface{}, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "number":
		switch value.(type) {
		case float64, float32, int, int32, int64:
			// Valid number types
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	default:
		return fmt.Errorf("unsupported type: %s", expectedType)
	}
	return nil
}

// containsWildcard checks if a path contains wildcard characters
func containsWildcard(path string) bool {
	return filepath.Base(path) != path || 
		   filepath.Dir(path) == "." ||
		   containsAny(path, "*?[")
}

// containsAny checks if string contains any of the given characters
func containsAny(s string, chars string) bool {
	for _, char := range chars {
		for _, c := range s {
			if c == char {
				return true
			}
		}
	}
	return false
}