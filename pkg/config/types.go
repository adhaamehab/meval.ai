package config

// Config represents the meval.yaml configuration
type Config struct {
	Experiment ExperimentConfig `yaml:"experiment"`
	Inputs     []InputConfig    `yaml:"inputs"`
	Outputs    []OutputConfig   `yaml:"outputs"`
	Evaluation EvaluationConfig `yaml:"evaluation"`
	Controls   ControlsConfig   `yaml:"controls"`
}

// ExperimentConfig represents experiment metadata
type ExperimentConfig struct {
	Name     string                 `yaml:"name"`
	Version  string                 `yaml:"version"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty"`
}

// InputConfig represents input source configuration
type InputConfig struct {
	ID     string                 `yaml:"id"`
	Format string                 `yaml:"format"`
	Config map[string]interface{} `yaml:"config"`
	Schema SchemaConfig           `yaml:"schema"`
}

// OutputConfig represents output source configuration
type OutputConfig struct {
	ID     string                 `yaml:"id"`
	Format string                 `yaml:"format"`
	Config map[string]interface{} `yaml:"config"`
	Schema SchemaConfig           `yaml:"schema"`
}

// SchemaConfig represents schema configuration
type SchemaConfig struct {
	Fields []FieldConfig `yaml:"fields"`
}

// FieldConfig represents a field in the schema
type FieldConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// EvaluationConfig represents evaluation configuration
type EvaluationConfig struct {
	Provider string                 `yaml:"provider"`
	Model    string                 `yaml:"model"`
	Params   map[string]interface{} `yaml:"params"`
	Auth     AuthConfig             `yaml:"auth"`
	Strategy string                 `yaml:"strategy"`
	Prompt   string                 `yaml:"prompt"`
	Mappings MappingsConfig         `yaml:"mappings"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	APIKeyEnv string `yaml:"api_key_env"`
}

// MappingsConfig represents field mappings
type MappingsConfig struct {
	Input  map[string]string `yaml:"input"`
	Output map[string]string `yaml:"output"`
}

// ControlsConfig represents execution controls
type ControlsConfig struct {
	Concurrency int    `yaml:"concurrency"`
	OnError     string `yaml:"on_error"`
}