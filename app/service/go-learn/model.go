package golearn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LocalAIService struct {
	baseURL string
	client  *http.Client
}

type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Response string `json:"response"` // Changed from Message to string
}

type SymptomAnalysis struct {
	PossibleDiseases []string `json:"possible_diseases"`
	Severity         int      `json:"severity"`
	NeedDoctor       bool     `json:"need_doctor"`
	Recommendations  string   `json:"recommendations"`
}

func NewLocalAIService() *LocalAIService {
	return &LocalAIService{
		baseURL: "http://localhost:11434/api/generate",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// performPost is a helper that marshals the given payload, sets headers, checks the status code,
// and decodes the response into the provided result.
func (s *LocalAIService) performPost(ctx context.Context, payload interface{}, result interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling API: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// GetPetAdvice sends a question to the LLM and returns the pet advice.
func (s *LocalAIService) GetPetAdvice(ctx context.Context, question string) (string, error) {
	systemPrompt := fmt.Sprintf(`Bạn là một bác sĩ thú y có kinh nghiệm. 
Hãy đưa ra lời khuyên chuyên môn về chăm sóc thú cưng.
Trả lời ngắn gọn, dễ hiểu và chính xác.
Nếu tình trạng nghiêm trọng, khuyên người dùng đến gặp bác sĩ.
Câu hỏi: %s`, question)

	chatReq := ChatRequest{
		Model:  "mistral",
		Prompt: systemPrompt,
	}

	var chatResp ChatResponse
	if err := s.performPost(ctx, chatReq, &chatResp); err != nil {
		return "", err
	}

	return chatResp.Response, nil
}

// AnalyzeSymptoms sends a list of symptoms to the LLM and decodes the JSON response into a SymptomAnalysis.
func (s *LocalAIService) AnalyzeSymptoms(ctx context.Context, symptoms []string) (*SymptomAnalysis, error) {
	prompt := fmt.Sprintf(`Phân tích các triệu chứng sau của thú cưng:
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
}`, symptoms)

	chatReq := ChatRequest{
		Model:  "mistral",
		Prompt: prompt,
	}

	var analysis SymptomAnalysis
	if err := s.performPost(ctx, chatReq, &analysis); err != nil {
		return nil, err
	}

	return &analysis, nil
}
