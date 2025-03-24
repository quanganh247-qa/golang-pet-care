package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
)

// GeminiService handles interactions with the Gemini API
type GeminiService struct {
	apiKey     string
	apiBaseURL string
	client     *http.Client
	cache      *CacheService
}

// NewGeminiService creates a new instance of GeminiService
func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey:     apiKey,
		apiBaseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
		client:     &http.Client{},
		cache:      nil, // We'll set this via a setter if needed
	}
}

// SetCache sets the cache service for GeminiService
func (s *GeminiService) SetCache(cache *CacheService) {
	s.cache = cache
}

// ProcessQuery sends the user's query to Gemini to understand intent and extract parameters
func (s *GeminiService) ProcessQuery(query string) (*models.GeminiQueryResponse, error) {
	// First check if this is a drug info request
	isDrugInfo, drugName := s.checkForDrugInfoRequest(query)
	if isDrugInfo && drugName != "" {
		// Create a response focused on drug info
		return &models.GeminiQueryResponse{
			QueryIntent: "drug_info",
			SearchTerms: drugName,
			IsDrugInfo:  true,
		}, nil
	}

	// If not a drug info request, process as a health trend query
	// Construct a prompt that guides Gemini to extract relevant health trend parameters
	// Using a more concise prompt to reduce token usage while maintaining quality
	prompt := fmt.Sprintf(`
Extract FDA search parameters from this query: "%s"
Return only JSON with these fields:
- queryIntent: Main intent (drug/device/food events, labels)
- searchTerms: Specific terms to search for
- timeFrame: Time period mentioned
- demographicFilters: Any demographic filters
- dataAggregation: How to aggregate data

Format: Valid JSON only, no explanation.
`, query)

	response, err := s.sendRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("error calling Gemini API: %w", err)
	}

	// Extract the JSON from the response
	jsonStr := extractJSONFromText(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in Gemini response")
	}

	// Parse the JSON
	var result models.GeminiQueryResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("error parsing Gemini response: %w", err)
	}

	return &result, nil
}

// checkForDrugInfoRequest detects if a query is specifically asking for drug information
func (s *GeminiService) checkForDrugInfoRequest(query string) (bool, string) {
	// Create a smaller, focused prompt to detect drug info requests
	// and extract the drug name if present
	prompt := fmt.Sprintf(`
Is this query asking about a specific medication/drug? "%s"
If yes, extract the exact drug name only.
If no, respond with "not_drug_info".
Return ONLY the drug name or "not_drug_info" with no additional text.
`, query)

	response, err := s.sendRequest(prompt)
	if err != nil {
		return false, ""
	}

	// Clean and process the response
	cleanResponse := strings.TrimSpace(response)

	// If response is not_drug_info or empty, it's not a drug info request
	if cleanResponse == "not_drug_info" || cleanResponse == "" {
		return false, ""
	}

	// Check for JSON format in response (shouldn't happen, but just in case)
	if strings.HasPrefix(cleanResponse, "{") {
		// Extract from possible JSON
		var resultMap map[string]string
		if err := json.Unmarshal([]byte(cleanResponse), &resultMap); err == nil {
			if drugName, ok := resultMap["drug_name"]; ok && drugName != "" {
				return true, drugName
			}
		}
	}

	// If we get here, the model likely returned just the drug name as instructed
	// Clean up common formatting issues
	cleanResponse = strings.Trim(cleanResponse, `"'`)
	cleanResponse = strings.TrimPrefix(cleanResponse, "Drug name: ")

	// Likely a drug name was returned
	return true, cleanResponse
}

// AnalyzeHealthData uses Gemini to analyze the FDA data and provide insights
func (s *GeminiService) AnalyzeHealthData(originalQuery string, fdaData interface{}, fdaQueryType string) (*models.AnalysisResult, error) {
	// Convert FDA data to a string, but reduce its size by limiting the items if it's too large
	// This helps reduce token usage when sending to Gemini API
	processedData, err := s.prepareDataForAnalysis(fdaData)
	if err != nil {
		return nil, fmt.Errorf("error processing FDA data: %w", err)
	}

	// Enhanced prompt to generate better chart data
	prompt := fmt.Sprintf(`
Analyze this OpenFDA %s data for query: "%s"

DATA:
%s

Return JSON with:
- summary: Key findings answering the user's question (concise HTML formatting allowed)
- sourceDetails: Brief note about data source and limitations

ALWAYS include these fields for visualization:
- chartType: Appropriate chart type (line/bar/pie/scatter)
- chartTitle: Clear descriptive chart title
- chartData: Structured data for visualization containing:
  * For line/bar charts: Include "labels" array and "datasets" array with at least one dataset containing "label" and "data" arrays
  * For pie charts: Include "labels" array and "datasets" array with at least one dataset containing "data" array
  * Limit to 10 most relevant data points
  * Include appropriate colors

Example chartData format:
{
  "labels": ["Item1", "Item2", "Item3"],
  "datasets": [{
    "label": "Dataset Label",
    "data": [10, 20, 30],
    "backgroundColor": ["rgba(54, 162, 235, 0.5)", "rgba(255, 99, 132, 0.5)"]
  }]
}

Format: Valid JSON only.
`, fdaQueryType, originalQuery, processedData)

	// Check cache first if available
	var result models.AnalysisResult
	cacheKey := fmt.Sprintf("analysis:%x", MD5Hash(originalQuery+fdaQueryType))

	// Try to get from cache
	if s.cache != nil && s.cache.GetJSON(cacheKey, &result) {
		return &result, nil
	}

	response, err := s.sendRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("error calling Gemini API for analysis: %w", err)
	}

	// Extract the JSON from the response
	jsonStr := extractJSONFromText(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in Gemini analysis response")
	}

	// Parse the JSON
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("error parsing Gemini analysis response: %w", err)
	}

	// Ensure chart data has a valid format if specified
	if result.ChartType != "" && result.ChartData == nil {
		// Create default chart data structure if missing
		result.ChartData = map[string]interface{}{
			"labels": []string{"No Data Available"},
			"datasets": []map[string]interface{}{
				{
					"label":           "No Data",
					"data":            []int{0},
					"backgroundColor": []string{"rgba(54, 162, 235, 0.5)"},
				},
			},
		}
	}

	// Cache the analysis result for 1 day (86400 seconds)
	if s.cache != nil {
		s.cache.Set(cacheKey, result, 86400)
	}

	return &result, nil
}

// prepareDataForAnalysis processes the FDA data to reduce size for token optimization
// while preserving the most important information
func (s *GeminiService) prepareDataForAnalysis(fdaData interface{}) (string, error) {
	// Convert data to map for processing
	var dataMap map[string]interface{}

	// First convert to JSON
	jsonData, err := json.Marshal(fdaData)
	if err != nil {
		return "", err
	}

	// Then parse back to map
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return string(jsonData), nil // If not a map, return as-is
	}

	// Process the data to reduce size
	// For example, if there are results arrays, limit to most relevant items
	if results, ok := dataMap["results"].([]interface{}); ok && len(results) > 20 {
		// Keep only first 20 results to reduce token count
		dataMap["results"] = results[:20]
		dataMap["note"] = "Data truncated to 20 results for analysis"
	}

	// If there's a count section with many terms, limit it
	if counts, ok := dataMap["results"].([]interface{}); ok {
		for i, count := range counts {
			if countMap, ok := count.(map[string]interface{}); ok {
				if terms, ok := countMap["term"].([]interface{}); ok && len(terms) > 10 {
					countMap["term"] = terms[:10]
					countMap["note"] = "Terms truncated to 10 items"
				}
			}

			// Only process first 10 count entries
			if i >= 10 {
				dataMap["results"] = counts[:10]
				break
			}
		}
	}

	// Convert back to JSON string
	processedJSON, err := json.Marshal(dataMap)
	if err != nil {
		return string(jsonData), nil // Return original if processing fails
	}

	return string(processedJSON), nil
}

// sendRequest sends a prompt to the Gemini API and returns the response text
func (s *GeminiService) sendRequest(prompt string) (string, error) {
	// Check cache first if available
	if s.cache != nil {
		// Create a hash of the prompt to use as cache key
		cacheKey := fmt.Sprintf("gemini:%x", MD5Hash(prompt))

		var cachedResponse string
		if s.cache.GetJSON(cacheKey, &cachedResponse) {
			return cachedResponse, nil
		}
	}

	// Construct the request payload with optimized parameters to reduce token usage
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
		// Add parameters to control response length and specificity
		"generationConfig": map[string]interface{}{
			"temperature":     0.2,  // Lower temperature for more focused responses
			"maxOutputTokens": 1024, // Limit response size
			"topP":            0.95, // More focused token selection
			"topK":            40,   // More focused token selection
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s?key=%s", s.apiBaseURL, s.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response to extract the text
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Navigate through the response structure to get the text
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("content not found in Gemini API response")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("parts not found in Gemini API response")
	}

	textPart, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("text not found in Gemini API response")
	}

	// Cache the result for 1 day (86400 seconds)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("gemini:%x", MD5Hash(prompt))
		s.cache.Set(cacheKey, textPart, 86400)
	}

	return textPart, nil
}

// MD5Hash creates an MD5 hash of a string
func MD5Hash(text string) string {
	// Simple hashing function for cache keys
	result := 0
	for i := 0; i < len(text); i++ {
		result = ((result << 5) - result) + int(text[i])
	}
	return fmt.Sprintf("%d", result)
}

// extractJSONFromText attempts to extract a JSON object from a text string
func extractJSONFromText(text string) string {
	// Look for text between { and }
	startIdx := strings.Index(text, "{")
	if startIdx == -1 {
		return ""
	}

	// Find the matching closing brace
	depth := 1
	for i := startIdx + 1; i < len(text); i++ {
		if text[i] == '{' {
			depth++
		} else if text[i] == '}' {
			depth--
			if depth == 0 {
				return text[startIdx : i+1]
			}
		}
	}

	return ""
}
