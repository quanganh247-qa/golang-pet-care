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

const systemPrompt = `Bạn là một trợ lý ảo tên PetCareHealth, chuyên hỗ trợ tư vấn sức khỏe con người và chăm sóc thú cưng.
- Luôn trả lời bằng tiếng Việt.
- Phân biệt rõ ràng giữa câu hỏi về con người và thú cưng.
- Giải thích ngắn gọn, dễ hiểu, kèm ví dụ cụ thể.
- Nếu không biết câu trả lời, hãy nói "Tôi chưa có thông tin về điều này."
- Không đưa ra lời khuyên y khoa hoặc thú y thay thế chuyên gia.`

type GeminiService struct {
	apiKey     string
	apiBaseURL string
	client     *http.Client
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey:     apiKey,
		apiBaseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
		client:     &http.Client{},
	}
}

// GeminiRequest represents a request to the Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents a content part in Gemini API request
type Content struct {
	Role  string  `json:"role,omitempty"`
	Parts []Parts `json:"parts"`
}

// Parts represents a part within a Content
type Parts struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from the Gemini API
type GeminiResponse struct {
	Candidates     []Candidate `json:"candidates"`
	PromptFeedback interface{} `json:"promptFeedback"`
}

// Candidate represents a response candidate
type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
		Role string `json:"role"`
	} `json:"content"`
	FinishReason string `json:"finishReason"`
	Index        int    `json:"index"`
}

// ProcessChatRequest processes a chat request and returns a response
func (s *GeminiService) ProcessChatRequest(request models.ChatRequest) (models.ChatResponse, error) {
	var response models.ChatResponse

	// Prepare conversation history
	var conversation []Content

	// Add system prompt
	conversation = append(conversation, Content{
		Role: "model",
		Parts: []Parts{
			{Text: systemPrompt},
		},
	})

	// Add previous messages if available
	if len(request.PreviousMessages) > 0 {
		for i, msg := range request.PreviousMessages {
			role := "user"
			if i%2 == 1 {
				role = "model"
			}
			conversation = append(conversation, Content{
				Role: role,
				Parts: []Parts{
					{Text: msg},
				},
			})
		}
	}

	// Add current message
	conversation = append(conversation, Content{
		Role: "user",
		Parts: []Parts{
			{Text: request.Message},
		},
	})

	// Create request payload
	geminiReq := GeminiRequest{
		Contents: conversation,
	}

	// Call Gemini API
	responseText, err := s.callGeminiAPI(geminiReq)
	if err != nil {
		return response, err
	}

	// Remove Markdown formatting like ** from the response text
	responseText = removeMarkdownFormatting(responseText)

	// Prepare response
	response = models.ChatResponse{
		Message:        responseText,
		ConversationID: request.ConversationID,
		Language:       "vi", // Default to Vietnamese as per system prompt
		PriorityLevel:  models.PriorityLow,
	}

	// Generate follow-up questions based on the context
	followUps, _ := s.generateFollowUpQuestions(request.Message, responseText)
	response.FollowUpQuestions = followUps

	return response, nil
}

// callGeminiAPI calls the Gemini API with the given request
func (s *GeminiService) callGeminiAPI(request GeminiRequest) (string, error) {
	// Serialize request body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s?key=%s", s.apiBaseURL, s.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract text from response
	if len(geminiResp.Candidates) > 0 &&
		len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "Tôi chưa có thông tin về điều này.", nil
}

// generateFollowUpQuestions generates follow-up questions based on the conversation context
func (s *GeminiService) generateFollowUpQuestions(userMessage, botResponse string) ([]string, error) {
	// Create a simple prompt to generate follow-up questions
	prompt := fmt.Sprintf(`Dựa vào cuộc trò chuyện sau, hãy đề xuất 3 câu hỏi tiếp theo mà người dùng có thể hỏi:

Người dùng: %s

PetCareHealth: %s

Đề xuất 3 câu hỏi tiếp theo (chỉ trả về danh sách câu hỏi, mỗi câu một dòng):`, userMessage, botResponse)

	// Create content for the follow-up question generator
	geminiReq := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Parts{
					{Text: prompt},
				},
			},
		},
	}

	// Call Gemini API
	responseText, err := s.callGeminiAPI(geminiReq)
	if err != nil {
		return []string{}, err
	}

	// Parse the questions from the response
	questions := parseQuestions(responseText)

	// Return up to 3 questions
	if len(questions) > 3 {
		questions = questions[:3]
	}

	return questions, nil
}

// parseQuestions extracts questions from the response text
func parseQuestions(text string) []string {
	// Split by line and clean up
	lines := strings.Split(text, "\n")
	var questions []string

	for _, line := range lines {
		// Clean the line
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Remove numbered bullets (1., 2., etc.)
		line = strings.TrimLeft(line, "0123456789.- ")

		// Remove other bullet points
		line = strings.TrimPrefix(line, "• ")
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		// Format the question properly for UI display
		if line != "" {
			// Ensure the question ends with a question mark
			if !strings.HasSuffix(line, "?") {
				line += "?"
			}

			// Capitalize first letter if not already
			if len(line) > 0 {
				firstChar := string(line[0])
				if firstChar == strings.ToLower(firstChar) {
					line = strings.ToUpper(firstChar) + line[1:]
				}
			}

			questions = append(questions, line)
		}
	}

	// Filter out non-question text or directives that might have been included
	var filteredQuestions []string
	for _, q := range questions {
		// Skip any lines that appear to be instructions rather than questions
		if strings.Contains(strings.ToLower(q), "câu hỏi") &&
			(strings.Contains(strings.ToLower(q), "đề xuất") ||
				strings.Contains(strings.ToLower(q), "gợi ý")) {
			continue
		}
		filteredQuestions = append(filteredQuestions, q)
	}

	return filteredQuestions
}

// removeMarkdownFormatting removes markdown formatting characters from text
func removeMarkdownFormatting(text string) string {
	// Remove bold markdown (**text**)
	text = strings.ReplaceAll(text, "**", "")

	// Có thể thêm các trường hợp loại bỏ định dạng Markdown khác ở đây nếu cần
	// Ví dụ: loại bỏ dấu gạch nghiêng (*text*), dấu gạch ngang (~text~), v.v.

	return text
}
