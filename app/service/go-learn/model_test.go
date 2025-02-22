package golearn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestLocalAIService_GetPetAdvice(t *testing.T) {
	service := NewLocalAIService()

	question := "Chó của tôi không ăn và nằm một chỗ 2 ngày nay, tôi nên làm gì?"

	response, err := service.GetPetAdvice(context.Background(), question)
	if err != nil {
		t.Errorf("GetPetAdvice failed: %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response, got empty string")
	}

	t.Logf("AI Response: %s", response)
}

func TestAnalyzeSymptoms(t *testing.T) {
	// Setup
	service := NewLocalAIService()

	// Test case
	symptoms := []string{"ho", "sốt", "chảy nước mũi"}

	// Make the request
	jsonData, err := json.Marshal(ChatRequest{
		Model: "mistral",
		Prompt: fmt.Sprintf(`Phân tích các triệu chứng sau của thú cưng:
		%v
		
		Hãy đưa ra:
		1. Các bệnh có thể gặp phải
		2. Mức độ nghiêm trọng (1-5)
		3. Khuyến nghị có nên đến bác sĩ ngay không
		
		Trả về dưới dạng JSON với format:
		{
			"possible_diseases": ["disease1", "disease2"],
			"severity": number,
			"need_doctor": boolean,
			"recommendations": "text"
		}`, symptoms),
	})
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Print request for debugging
	t.Logf("Request: %s", string(jsonData))

	// Make API call
	resp, err := http.Post(
		service.baseURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print raw response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	t.Logf("Raw Response: %s", string(body))

	// Try to parse the response
	var aiResp struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &aiResp); err != nil {
		t.Fatalf("Failed to parse AI response: %v", err)
	}
	t.Logf("AI Response Text: %s", aiResp.Response)

	// Try to parse the JSON from AI's response
	var analysis SymptomAnalysis
	if err := json.Unmarshal([]byte(aiResp.Response), &analysis); err != nil {
		t.Fatalf("Failed to parse analysis JSON: %v", err)
	}

	t.Logf("Final Parsed Analysis: %+v", analysis)

	// Validation
	if len(analysis.PossibleDiseases) == 0 {
		t.Error("Expected possible diseases, got empty slice")
	}
	if analysis.Severity < 1 || analysis.Severity > 5 {
		t.Errorf("Expected severity between 1-5, got %d", analysis.Severity)
	}
	if analysis.Recommendations == "" {
		t.Error("Expected non-empty recommendations")
	}
}
func TestAnalyzeSymptoms_Simple(t *testing.T) {
	fmt.Println("\n=== Starting Test ===")

	// Setup
	service := NewLocalAIService() // Changed to /api/generate

	// Test case with simple symptoms
	symptoms := []string{"ho", "sốt"}
	fmt.Printf("Testing symptoms: %v\n", symptoms)

	// Create request body - simplified for Ollama
	requestBody := map[string]interface{}{
		"model": "mistral",
		"prompt": fmt.Sprintf(`Return a JSON object analyzing these pet symptoms: %v

        The response must be ONLY a valid JSON object with this exact structure:
        {
            "possible_diseases": ["disease1", "disease2"],
            "severity": 3,
            "need_doctor": true,
            "recommendations": "text here"
        }

        DO NOT include any other text, ONLY the JSON object.`, symptoms),
		"format": "json",
		// Option: if your API supports disabling streaming, you might add:
		// "stream": false,
	}

	// Marshal request
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	fmt.Println("\n=== Request ===")
	fmt.Printf("Request Body: %s\n", string(jsonData))

	// Make API call
	resp, err := http.Post(
		service.baseURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("\n=== Response Status: %s ===\n", resp.Status)

	// Instead of reading the entire response at once, use a JSON decoder to read each fragment.
	decoder := json.NewDecoder(resp.Body)
	var fullResponse string
	for {
		var part struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}
		if err := decoder.Decode(&part); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("Failed to decode response fragment: %v", err)
		}

		fullResponse += part.Response

		// If the stream signals it's done, exit the loop.
		if part.Done {
			break
		}
	}

	fmt.Println("\n=== Aggregated AI Response ===")
	fmt.Printf("Combined Response: %s\n", fullResponse)

	// Now try to parse the aggregated JSON from the AI's response
	var analysis SymptomAnalysis
	if err := json.Unmarshal([]byte(fullResponse), &analysis); err != nil {
		t.Fatalf("Failed to parse analysis: %v", err)
	}

	fmt.Println("\n=== Final Parsed Analysis ===")
	fmt.Printf("Analysis: %+v\n", analysis)

	// Validation
	if len(analysis.PossibleDiseases) == 0 {
		t.Error("Expected possible diseases, got empty slice")
	}
	if analysis.Severity < 1 || analysis.Severity > 5 {
		t.Error("Expected severity between 1-5")
	}
	if analysis.Recommendations == "" {
		t.Error("Expected non-empty recommendations")
	}

	fmt.Println("\n=== Test Complete ===")
}
