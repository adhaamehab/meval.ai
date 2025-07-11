package config

import (
	"strings"
	"testing"
)

func TestReader_Read(t *testing.T) {
	yamlContent := `experiment:
  name: test-evaluation
  version: 0.1
  metadata:
    author: test
    environment: development

inputs:
  - id: predictions
    format: json
    config:
      path: ./data/input.json
    schema:
      fields:
        - name: text
          type: string
        - name: predicted_sentiment
          type: string

outputs:
  - id: eval-results
    format: json
    config:
      path: ./data/output.json
    schema:
      fields:
        - name: text
          type: string
        - name: evaluation_sentiment
          type: string

evaluation:
  provider: gemini
  model: text-bison@001
  params:
    temperature: 0.2
    max_tokens: 64
  auth:
    api_key_env: GEMINI_API_KEY
  strategy: classification
  prompt: |
    Given the text below, output exactly one word:
    positive, negative, or neutral.
    
    Text: {{text}}
  mappings:
    input:
      text: text
    output:
      evaluation_sentiment: $.label

controls:
  concurrency: 4
  on_error: retry
`

	reader := NewReader()
	config, err := reader.Read(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Validate experiment fields
	if config.Experiment.Version != "0.1" {
		t.Errorf("Expected experiment version 0.1, got %s", config.Experiment.Version)
	}

	if config.Experiment.Name != "test-evaluation" {
		t.Errorf("Expected experiment name test-evaluation, got %s", config.Experiment.Name)
	}

	if config.Experiment.Metadata["author"] != "test" {
		t.Errorf("Expected metadata author test, got %v", config.Experiment.Metadata["author"])
	}

	// Validate inputs
	if len(config.Inputs) != 1 {
		t.Errorf("Expected 1 input, got %d", len(config.Inputs))
	}

	if config.Inputs[0].ID != "predictions" {
		t.Errorf("Expected input id predictions, got %s", config.Inputs[0].ID)
	}

	// Validate evaluation
	if config.Evaluation.Provider != "gemini" {
		t.Errorf("Expected provider gemini, got %s", config.Evaluation.Provider)
	}

	// Validate controls
	if config.Controls.Concurrency != 4 {
		t.Errorf("Expected concurrency 4, got %d", config.Controls.Concurrency)
	}

	// Validate with validator
	validator := NewValidator()
	if err := validator.Validate(config); err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
}