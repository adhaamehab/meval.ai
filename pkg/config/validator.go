package config

import (
	"fmt"
	"strings"
)

// Validator implements configuration validation
type Validator struct{}

// NewValidator creates a new config validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates the configuration
func (v *Validator) Validate(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	// Validate experiment
	if err := v.validateExperiment(config.Experiment); err != nil {
		return err
	}

	// Validate inputs
	if len(config.Inputs) == 0 {
		return fmt.Errorf("at least one input is required")
	}

	for i, input := range config.Inputs {
		if err := v.validateInput(input, i); err != nil {
			return err
		}
	}

	// Validate outputs
	if len(config.Outputs) == 0 {
		return fmt.Errorf("at least one output is required")
	}

	for i, output := range config.Outputs {
		if err := v.validateOutput(output, i); err != nil {
			return err
		}
	}

	// Validate evaluation
	if err := v.validateEvaluation(config.Evaluation); err != nil {
		return err
	}

	// Validate controls
	if err := v.validateControls(config.Controls); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateExperiment(exp ExperimentConfig) error {
	if exp.Name == "" {
		return fmt.Errorf("experiment.name is required")
	}

	if exp.Version == "" {
		return fmt.Errorf("experiment.version is required")
	}

	return nil
}

func (v *Validator) validateInput(input InputConfig, index int) error {
	if input.ID == "" {
		return fmt.Errorf("input[%d]: id is required", index)
	}

	if input.Format == "" {
		return fmt.Errorf("input[%d]: format is required", index)
	}

	supportedFormats := []string{"json", "csv", "parquet"}
	if !contains(supportedFormats, input.Format) {
		return fmt.Errorf("input[%d]: unsupported format %s", index, input.Format)
	}

	if len(input.Config) == 0 {
		return fmt.Errorf("input[%d]: config is required", index)
	}

	if _, ok := input.Config["path"]; !ok {
		return fmt.Errorf("input[%d]: config.path is required", index)
	}

	return v.validateSchema(input.Schema, fmt.Sprintf("input[%d]", index))
}

func (v *Validator) validateOutput(output OutputConfig, index int) error {
	if output.ID == "" {
		return fmt.Errorf("output[%d]: id is required", index)
	}

	if output.Format == "" {
		return fmt.Errorf("output[%d]: format is required", index)
	}

	supportedFormats := []string{"json", "csv", "parquet"}
	if !contains(supportedFormats, output.Format) {
		return fmt.Errorf("output[%d]: unsupported format %s", index, output.Format)
	}

	if len(output.Config) == 0 {
		return fmt.Errorf("output[%d]: config is required", index)
	}

	if _, ok := output.Config["path"]; !ok {
		return fmt.Errorf("output[%d]: config.path is required", index)
	}

	return v.validateSchema(output.Schema, fmt.Sprintf("output[%d]", index))
}

func (v *Validator) validateSchema(schema SchemaConfig, prefix string) error {
	if len(schema.Fields) == 0 {
		return fmt.Errorf("%s: schema must have at least one field", prefix)
	}

	for i, field := range schema.Fields {
		if field.Name == "" {
			return fmt.Errorf("%s.schema.fields[%d]: name is required", prefix, i)
		}

		if field.Type == "" {
			return fmt.Errorf("%s.schema.fields[%d]: type is required", prefix, i)
		}

		supportedTypes := []string{"string", "number", "boolean", "array", "object"}
		if !contains(supportedTypes, field.Type) {
			return fmt.Errorf("%s.schema.fields[%d]: unsupported type %s", prefix, i, field.Type)
		}
	}

	return nil
}

func (v *Validator) validateEvaluation(eval EvaluationConfig) error {
	if eval.Provider == "" {
		return fmt.Errorf("evaluation.provider is required")
	}

	supportedProviders := []string{"openai", "anthropic", "gemini", "bedrock"}
	if !contains(supportedProviders, eval.Provider) {
		return fmt.Errorf("evaluation: unsupported provider %s", eval.Provider)
	}

	if eval.Model == "" {
		return fmt.Errorf("evaluation.model is required")
	}

	if eval.Auth.APIKeyEnv == "" {
		return fmt.Errorf("evaluation.auth.api_key_env is required")
	}

	if eval.Strategy == "" {
		return fmt.Errorf("evaluation.strategy is required")
	}

	supportedStrategies := []string{"classification", "extraction", "generation"}
	if !contains(supportedStrategies, eval.Strategy) {
		return fmt.Errorf("evaluation: unsupported strategy %s", eval.Strategy)
	}

	if eval.Prompt == "" {
		return fmt.Errorf("evaluation.prompt is required")
	}

	if !strings.Contains(eval.Prompt, "{{") {
		return fmt.Errorf("evaluation.prompt must contain at least one template variable")
	}

	return nil
}

func (v *Validator) validateControls(controls ControlsConfig) error {
	if controls.Concurrency <= 0 {
		return fmt.Errorf("controls.concurrency must be greater than 0")
	}

	if controls.OnError == "" {
		return fmt.Errorf("controls.on_error is required")
	}

	supportedErrorHandling := []string{"retry", "skip", "fail"}
	if !contains(supportedErrorHandling, controls.OnError) {
		return fmt.Errorf("controls: unsupported on_error value %s", controls.OnError)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}