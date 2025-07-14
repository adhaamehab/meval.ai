package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adhaamehab/meval.ai/pkg/config"
	"github.com/adhaamehab/meval.ai/pkg/sources"
)

func main() {
	// Example 1: Read JSON array file
	fmt.Println("=== Example 1: JSON Array ===")
	readJSONArray()

	// Example 2: Read JSON lines file
	fmt.Println("\n=== Example 2: JSON Lines ===")
	readJSONLines()

	// Example 3: Read with wildcards
	fmt.Println("\n=== Example 3: Wildcards ===")
	readWithWildcards()
}

func readJSONArray() {
	// Configuration for JSON source
	cfg := map[string]interface{}{
		"path": "./data/predictions.json",
		"mode": "array",
	}

	// Define schema
	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	// Create source
	source, err := sources.NewJSONSource(cfg, schema)
	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}
	defer source.Close()

	// Read records
	records, err := source.Read(context.Background())
	if err != nil {
		log.Fatalf("Failed to read records: %v", err)
	}

	fmt.Printf("Read %d records\n", len(records))
	for i, record := range records {
		fmt.Printf("Record %d: %v\n", i+1, record)
	}
}

func readJSONLines() {
	cfg := map[string]interface{}{
		"path": "./data/predictions.jsonl",
		"mode": "lines",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := sources.NewJSONSource(cfg, schema)
	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}
	defer source.Close()

	records, err := source.Read(context.Background())
	if err != nil {
		log.Fatalf("Failed to read records: %v", err)
	}

	fmt.Printf("Read %d records from JSON lines\n", len(records))
}

func readWithWildcards() {
	// Read all JSON files matching pattern
	cfg := map[string]interface{}{
		"path": "./data/batch_*.json",
		"mode": "array",
	}

	schema := config.SchemaConfig{
		Fields: []config.FieldConfig{
			{Name: "text", Type: "string"},
			{Name: "predicted_sentiment", Type: "string"},
		},
	}

	source, err := sources.NewJSONSource(cfg, schema)
	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}
	defer source.Close()

	records, err := source.Read(context.Background())
	if err != nil {
		log.Fatalf("Failed to read records: %v", err)
	}

	fmt.Printf("Read %d total records from multiple files\n", len(records))
}