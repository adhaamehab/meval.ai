package evaluators

import (
	"context"

	"github.com/adhaamehab/meval.ai/pkg/config"
	"github.com/adhaamehab/meval.ai/pkg/sources"
)

// Evaluator defines the interface for evaluation providers
type Evaluator interface {
	// Evaluate performs evaluation on a single record
	Evaluate(ctx context.Context, record sources.Record, prompt string) (Result, error)
	// BatchEvaluate performs evaluation on multiple records
	BatchEvaluate(ctx context.Context, records []sources.Record, prompt string) ([]Result, error)
}

// Result contains the result of an evaluation
type Result struct {
	Input    sources.Record
	Output   map[string]interface{}
	Metadata map[string]interface{}
	Error    error
}

// Factory creates evaluators based on provider
type Factory interface {
	CreateEvaluator(provider string, config config.EvaluationConfig) (Evaluator, error)
}