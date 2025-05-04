package inference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/quanganh247-qa/go-blog-be/app/service"
)

const (
	defaultRoboflowAPIURL = "https://serverless.roboflow.com"
)

// RoboflowClient is a client for interacting with Roboflow inference API
type RoboflowClient struct {
	client *service.ClientBuilder
	apiKey string
}

// NewRoboflowClient creates a new Roboflow inference client
func NewRoboflowClient(apiKey string, apiURL string) *RoboflowClient {
	if apiURL == "" {
		apiURL = defaultRoboflowAPIURL
	}

	httpClient := &http.Client{}
	clientBuilder := &service.ClientBuilder{
		HttpClient:  httpClient,
		EndpointUrl: apiURL,
		UserAgent:   "golang-pet-care/roboflow-client",
	}

	return &RoboflowClient{
		client: clientBuilder,
		apiKey: apiKey,
	}
}

// Generic InferenceResult for all response types
type InferenceResult struct {
	Predictions interface{} `json:"predictions"`
	Time        float64     `json:"time"`
}

// CatPrediction represents a cat prediction from Roboflow model
type CatPrediction struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Class  string  `json:"class"`
	Score  float64 `json:"confidence"`
}

// DogPrediction represents a dog prediction from Roboflow model
type DogPrediction struct {
	Class string  `json:"class"`
	Score float64 `json:"confidence"`
}

// Infer sends an image to Roboflow for inference
func (r *RoboflowClient) Infer(imagePath string, modelID string, breed string) (*InferenceResult, error) {
	// Read image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	return r.InferFromBytes(imageData, modelID, breed)
}

// InferFromBytes sends image bytes to Roboflow for inference
func (r *RoboflowClient) InferFromBytes(imageData []byte, modelID, breed string) (*InferenceResult, error) {
	// Build the URL with model ID and API key
	url := fmt.Sprintf("%s/%s?api_key=%s", r.client.EndpointUrl, modelID, r.apiKey)

	// Create request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := r.client.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("inference request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read the raw response first
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create a generic result to return
	result := &InferenceResult{}

	// Parse based on breed type
	if breed == "dog" {
		// Parse into a temporary struct with the correct prediction type
		var tempResult struct {
			Predictions []DogPrediction `json:"predictions"`
			Time        float64         `json:"time"`
		}

		if err := json.Unmarshal(body, &tempResult); err != nil {
			return nil, fmt.Errorf("failed to parse dog inference result: %w", err)
		}

		// Transfer the data to our generic result
		result.Predictions = tempResult.Predictions
		result.Time = tempResult.Time
	} else if breed == "cat" {
		// Parse into a temporary struct with the correct prediction type
		var tempResult struct {
			Predictions []CatPrediction `json:"predictions"`
			Time        float64         `json:"time"`
		}

		if err := json.Unmarshal(body, &tempResult); err != nil {
			return nil, fmt.Errorf("failed to parse cat inference result: %w", err)
		}

		// Transfer the data to our generic result
		result.Predictions = tempResult.Predictions

		result.Time = tempResult.Time
	} else {
		// If no specific breed is provided, just decode the raw JSON
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse generic inference result: %w", err)
		}
	}

	return result, nil
}
