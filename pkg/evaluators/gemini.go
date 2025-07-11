package evaluators

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adhaamehab/meval.ai/pkg/config"
	"github.com/adhaamehab/meval.ai/pkg/sources"
)

// GeminiEvaluator implements the Evaluator interface for Google Gemini
type GeminiEvaluator struct {
	apiKey     string
	model      string
	params     map[string]interface{}
	httpClient *http.Client
}

// NewGeminiEvaluator creates a new Gemini evaluator
func NewGeminiEvaluator(cfg config.EvaluationConfig) (*GeminiEvaluator, error) {
	apiKey := os.Getenv(cfg.Auth.APIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("API key environment variable %s is not set", cfg.Auth.APIKeyEnv)
	}

	return &GeminiEvaluator{
		apiKey: apiKey,
		model:  cfg.Model,
		params: cfg.Params,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Evaluate performs evaluation on a single record
func (g *GeminiEvaluator) Evaluate(ctx context.Context, record sources.Record, prompt string) (Result, error) {
	// Apply prompt templating
	processedPrompt := g.applyPromptTemplate(prompt, record)

	// Prepare request
	requestBody := g.buildRequestBody(processedPrompt)

	// Make API call
	response, err := g.makeAPICall(ctx, requestBody)
	if err != nil {
		return Result{
			Input: record,
			Error: err,
		}, err
	}

	// Parse response
	output, metadata, err := g.parseResponse(response)
	if err != nil {
		return Result{
			Input: record,
			Error: err,
		}, err
	}

	return Result{
		Input:    record,
		Output:   output,
		Metadata: metadata,
		Error:    nil,
	}, nil
}

// BatchEvaluate performs evaluation on multiple records
func (g *GeminiEvaluator) BatchEvaluate(ctx context.Context, records []sources.Record, prompt string) ([]Result, error) {
	results := make([]Result, len(records))

	// Process each record individually
	// TODO: Implement concurrent processing based on controls.concurrency
	for i, record := range records {
		result, err := g.Evaluate(ctx, record, prompt)
		if err != nil {
			results[i] = Result{
				Input: record,
				Error: err,
			}
		} else {
			results[i] = result
		}
	}

	return results, nil
}

// applyPromptTemplate replaces template variables with values from the record
func (g *GeminiEvaluator) applyPromptTemplate(prompt string, record sources.Record) string {
	processedPrompt := prompt

	// Replace template variables like {{field_name}} with actual values
	for key, value := range record {
		placeholder := fmt.Sprintf("{{%s}}", key)
		valueStr := fmt.Sprintf("%v", value)
		processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, valueStr)
	}

	return processedPrompt
}

// buildRequestBody builds the API request body
func (g *GeminiEvaluator) buildRequestBody(prompt string) map[string]interface{} {
	// Build request based on Gemini API format
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	// Add generation config from params
	if g.params != nil {
		generationConfig := make(map[string]interface{})
		
		if temp, ok := g.params["temperature"]; ok {
			generationConfig["temperature"] = temp
		}
		
		if maxTokens, ok := g.params["max_tokens"]; ok {
			generationConfig["maxOutputTokens"] = maxTokens
		}
		
		if len(generationConfig) > 0 {
			requestBody["generationConfig"] = generationConfig
		}
	}

	return requestBody
}

// makeAPICall makes the HTTP request to Gemini API
func (g *GeminiEvaluator) makeAPICall(ctx context.Context, requestBody map[string]interface{}) (map[string]interface{}, error) {
	// Construct API URL
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	// Marshal request body
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error: %v", errorResponse)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// parseResponse extracts the output and metadata from Gemini API response
func (g *GeminiEvaluator) parseResponse(response map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
	output := make(map[string]interface{})
	metadata := make(map[string]interface{})

	// Extract candidates
	candidates, ok := response["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return nil, nil, fmt.Errorf("no candidates in response")
	}

	// Get first candidate
	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("invalid candidate format")
	}

	// Extract content
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("no content in candidate")
	}

	// Extract parts
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return nil, nil, fmt.Errorf("no parts in content")
	}

	// Get text from first part
	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("invalid part format")
	}

	text, ok := part["text"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("no text in part")
	}

	// Store the raw text response
	output["response"] = text

	// Try to parse as JSON if it looks like JSON
	trimmedText := strings.TrimSpace(text)
	if strings.HasPrefix(trimmedText, "{") || strings.HasPrefix(trimmedText, "[") {
		var jsonOutput interface{}
		if err := json.Unmarshal([]byte(trimmedText), &jsonOutput); err == nil {
			output["parsed"] = jsonOutput
		}
	}

	// Add metadata
	if safetyRatings, ok := candidate["safetyRatings"]; ok {
		metadata["safetyRatings"] = safetyRatings
	}

	if finishReason, ok := candidate["finishReason"]; ok {
		metadata["finishReason"] = finishReason
	}

	if usageMetadata, ok := response["usageMetadata"]; ok {
		metadata["usage"] = usageMetadata
	}

	return output, metadata, nil
}