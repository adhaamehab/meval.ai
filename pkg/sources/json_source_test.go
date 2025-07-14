package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adhaamehab/meval.ai/pkg/config"
)

func TestJSONSource_ReadArray(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	testData := []map[string]interface{}{
		{
			"text":               "This is positive",
			"predicted_sentiment": "positive",
		},
		{
			"text":               "This is negative",
			"predicted_sentiment": "negative",
		},
	}

	// Write test data
	data, _ := json.MarshalIndent(testData, "", "  ")
	if err := os.WriteFile(testFile, data, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create JSON source
	cfg := map[string]interface{}{
		"path": testFile,
		"mode": "array",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := NewJSONSource(cfg, schema)
	if err != nil {
		t.Fatalf("Failed to create JSON source: %v", err)
	}

	// Read records
	records, err := source.Read(context.Background())
	if err != nil {
		t.Fatalf("Failed to read records: %v", err)
	}

	// Validate results
	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	if records[0]["text"] != "This is positive" {
		t.Errorf("Expected first text to be 'This is positive', got %v", records[0]["text"])
	}

	if records[1]["predicted_sentiment"] != "negative" {
		t.Errorf("Expected second sentiment to be 'negative', got %v", records[1]["predicted_sentiment"])
	}
}

func TestJSONSource_ReadLines(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl")

	testData := `{"text": "Line 1", "predicted_sentiment": "positive"}
{"text": "Line 2", "predicted_sentiment": "negative"}
{"text": "Line 3", "predicted_sentiment": "neutral"}`

	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create JSON source
	cfg := map[string]interface{}{
		"path": testFile,
		"mode": "lines",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := NewJSONSource(cfg, schema)
	if err != nil {
		t.Fatalf("Failed to create JSON source: %v", err)
	}

	// Read records
	records, err := source.Read(context.Background())
	if err != nil {
		t.Fatalf("Failed to read records: %v", err)
	}

	// Validate results
	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	if records[2]["text"] != "Line 3" {
		t.Errorf("Expected third text to be 'Line 3', got %v", records[2]["text"])
	}
}

func TestJSONSource_Wildcards(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	
	// Create multiple JSON files
	for i := 1; i <= 3; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("data%d.json", i))
		testData := []map[string]interface{}{
			{
				"text":               fmt.Sprintf("Text from file %d", i),
				"predicted_sentiment": "positive",
			},
		}
		data, _ := json.MarshalIndent(testData, "", "  ")
		if err := os.WriteFile(testFile, data, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create JSON source with wildcard
	cfg := map[string]interface{}{
		"path": filepath.Join(tmpDir, "data*.json"),
		"mode": "array",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := NewJSONSource(cfg, schema)
	if err != nil {
		t.Fatalf("Failed to create JSON source: %v", err)
	}

	// Read records
	records, err := source.Read(context.Background())
	if err != nil {
		t.Fatalf("Failed to read records: %v", err)
	}

	// Validate results
	if len(records) != 3 {
		t.Errorf("Expected 3 records from 3 files, got %d", len(records))
	}
}

func TestJSONSource_SchemaValidation(t *testing.T) {
	// Create temporary test file with invalid data
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	testData := []map[string]interface{}{
		{
			"text": "This is text",
			// Missing predicted_sentiment field
		},
	}

	data, _ := json.MarshalIndent(testData, "", "  ")
	if err := os.WriteFile(testFile, data, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create JSON source
	cfg := map[string]interface{}{
		"path": testFile,
		"mode": "array",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := NewJSONSource(cfg, schema)
	if err != nil {
		t.Fatalf("Failed to create JSON source: %v", err)
	}

	// Read should fail due to schema validation
	_, err = source.Read(context.Background())
	if err == nil {
		t.Error("Expected schema validation error, got nil")
	}
}

func TestJSONSource_Write(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.json")

	// Create JSON source
	cfg := map[string]interface{}{
		"path": testFile,
		"mode": "array",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "result", Type: "string"},
		},
	}

	source, err := NewJSONSource(cfg, schema)
	if err != nil {
		t.Fatalf("Failed to create JSON source: %v", err)
	}

	// Write records
	records := []Record{
		{
			"text":   "Test 1",
			"result": "positive",
		},
		{
			"text":   "Test 2",
			"result": "negative",
		},
	}

	if err := source.Write(context.Background(), records); err != nil {
		t.Fatalf("Failed to write records: %v", err)
	}

	if err := source.Close(); err != nil {
		t.Fatalf("Failed to close source: %v", err)
	}

	// Verify written file
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	var written []Record
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("Failed to unmarshal written data: %v", err)
	}

	if len(written) != 2 {
		t.Errorf("Expected 2 written records, got %d", len(written))
	}

	if written[0]["text"] != "Test 1" {
		t.Errorf("Expected first text to be 'Test 1', got %v", written[0]["text"])
	}
}

func TestJSONSource_TypeValidation(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		fieldType   string
		shouldError bool
	}{
		{"valid string", "hello", "string", false},
		{"invalid string", 123, "string", true},
		{"valid number", 42.5, "number", false},
		{"valid int", 42, "number", false},
		{"invalid number", "not a number", "number", true},
		{"valid boolean", true, "boolean", false},
		{"invalid boolean", "true", "boolean", true},
		{"valid array", []interface{}{1, 2, 3}, "array", false},
		{"invalid array", "not an array", "array", true},
		{"valid object", map[string]interface{}{"key": "value"}, "object", false},
		{"invalid object", "not an object", "object", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldType(tt.value, tt.fieldType)
			if (err != nil) != tt.shouldError {
				t.Errorf("validateFieldType() error = %v, shouldError = %v", err, tt.shouldError)
			}
		})
	}
}