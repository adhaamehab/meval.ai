package evaluators

import (
	"fmt"

	"github.com/adhaamehab/meval.ai/pkg/config"
)

// DefaultFactory implements the Factory interface for evaluators
type DefaultFactory struct{}

// NewDefaultFactory creates a new evaluator factory
func NewDefaultFactory() *DefaultFactory {
	return &DefaultFactory{}
}

// CreateEvaluator creates an evaluator based on provider and configuration
func (f *DefaultFactory) CreateEvaluator(provider string, cfg config.EvaluationConfig) (Evaluator, error) {
	switch provider {
	case "gemini":
		return NewGeminiEvaluator(cfg)
	case "openai":
		return nil, fmt.Errorf("OpenAI evaluator not yet implemented")
	case "anthropic":
		return nil, fmt.Errorf("Anthropic evaluator not yet implemented")
	case "bedrock":
		return nil, fmt.Errorf("Bedrock evaluator not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}