package controller

import (
	"context"

	"github.com/adhaamehab/meval.ai/pkg/config"
)

// Controller manages the execution flow and concurrency
type Controller interface {
	// Execute runs the evaluation pipeline
	Execute(ctx context.Context, cfg *config.Config) error
	// Stop gracefully stops the execution
	Stop() error
}