# meval.ai

AI-powered evaluation platform for machine learning models.

## Overview

meval.ai uses LLMs to evaluate ML model outputs. It provides a YAML-based configuration system for defining evaluation experiments.

## Current Status

This is an early MVP implementation with:
- ✅ YAML configuration reader and validator
- ✅ Extensible interface design
- ✅ Schema validation
- ✅ JSON source with wildcard support
- 🚧 CLI implementation (coming soon)
- 🚧 CSV and Parquet sources (coming soon)
- 🚧 Evaluator implementations (coming soon)

## Configuration

Create a `meval.yaml` file:

```yaml
experiment:
  name: sentiment-evaluation
  version: 0.1
  metadata:
    author: your-name
    description: Evaluate sentiment analysis model

inputs:
  - id: predictions
    format: json
    config:
      path: ./data/input.json  # Can use wildcards: ./data/*.json
      mode: array              # "array" or "lines" (default: array)
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
        - name: predicted_sentiment
          type: string
        - name: evaluation_sentiment
          type: string
        - name: explanation
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
      predicted_sentiment: prediction
    output:
      evaluation_sentiment: $.label
      explanation: $.explanation

controls:
  concurrency: 4
  on_error: retry
```

## Development

### Project Structure

```
meval.ai/
├── cmd/                # CLI commands (TBD)
├── pkg/
│   ├── config/        # Configuration management
│   ├── sources/       # Data source interfaces and implementations
│   ├── evaluators/    # AI evaluator interfaces and implementations
│   ├── controller/    # Pipeline controller interface and implementation
│   └── results/       # Result handling (TBD)
└── go.mod
```

### Building

```bash
# Clone the repository
git clone https://github.com/adhaamehab/meval.ai.git
cd meval.ai

# Install dependencies
go mod download

# Run tests
go test ./...
```

### Implemented Components

#### Config Package
- `Reader`: Reads and parses YAML configuration files
- `Validator`: Validates configuration structure and values
- Support for experiment metadata with key-value pairs

#### Sources Package
- `JSONSource`: Reads/writes JSON files with support for:
  - JSON array format (standard JSON array of objects)
  - JSON lines format (one JSON object per line)
  - Wildcard path patterns (e.g., `data/*.json`)
  - Schema validation for all records
- `Factory`: Creates sources based on format configuration

#### Package Organization
Each package owns its interfaces and implementations:
- `sources`: Source interface and implementations (JSON implemented, CSV/Parquet coming)
- `evaluators`: Evaluator interface and future provider implementations
- `controller`: Controller interface for pipeline orchestration
- `config`: Configuration types, reader, and validator with their interfaces

### Supported Configuration

- **Experiment**: name, version, metadata (key-value pairs)
- **Inputs/Outputs**: JSON, CSV, Parquet formats
- **Providers**: OpenAI, Anthropic, Gemini, Bedrock
- **Strategies**: classification, extraction, generation
- **Error Handling**: retry, skip, fail


## License
GNU Affero General Public License v3.0 (AGPL-3.0)