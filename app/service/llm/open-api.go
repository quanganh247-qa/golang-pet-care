package llm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/patrickmn/go-cache"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
)

// Constants
const (
	ActionAppointment  = "appointment"
	ActionPetLog       = "pet_log"
	ActionPetSchedule  = "pet_schedule"
	MaxRetries         = 3
	CacheExpiration    = 5 * time.Minute
	CacheCleanup       = 10 * time.Minute
	LanguageEnglish    = "en"
	LanguageVietnamese = "vi"
)

// Global variables
var (
	sem           = semaphore.NewWeighted(15)
	responseCache = cache.New(CacheExpiration, CacheCleanup)
	logger        *zap.Logger
)

// Initialize logger
func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
}

type SuggestionResponse interface {
	GetAction() string
	Validate() error
}

// BaseResponse contains common fields for all response types
type BaseResponse struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (b BaseResponse) GetAction() string {
	return b.Action
}

type GeminiRequest struct {
	Prompt string `json:"prompt"`
}

type AppointmentResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title           string `json:"title"`
		PetName         string `json:"pet_name"`
		AppointmentType string `json:"appointment_type"`
		Date            string `json:"date"`
		Time            string `json:"time"`
	} `json:"parameters"`
}

type LogResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title    string `json:"title"`
		PetName  string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

type ScheduleResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title    string `json:"title"	`
		PetName  string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

type ActionResponse struct {
	Action string `json:"action"`
}

// Enhanced prompt templates
const (
	actionDeterminationPromptEN = `As an AI assistant, your task is to interpret the user's request and determine the appropriate action to take. The possible actions are:

- "appointment": For scheduling a one-time appointment or event.
- "pet_log": For logging a one-time activity or event related to a pet.
- "pet_schedule": For setting up a recurring schedule (e.g., daily, weekly).

Task:
Analyze the following user description and determine the most suitable action type.

User Input:
%s

Response Format:
Return only a valid JSON object, structured as follows:
{
	"action": "<appointment | pet_log | pet_schedule>"
}`

	actionDeterminationPromptVI = `Là một trợ lý AI, nhiệm vụ của bạn là phân tích yêu cầu của người dùng và xác định hành động phù hợp. Các hành động có thể là:

- "appointment": Để lên lịch hẹn hoặc sự kiện một lần.
- "pet_log": Để ghi lại hoạt động hoặc sự kiện một lần liên quan đến thú cưng.
- "pet_schedule": Để thiết lập lịch trình định kỳ (ví dụ: hàng ngày, hàng tuần).

Nhiệm vụ:
Phân tích mô tả của người dùng sau đây và xác định loại hành động phù hợp nhất.

Đầu vào người dùng:
%s

Định dạng phản hồi:
Chỉ trả về một đối tượng JSON hợp lệ, có cấu trúc như sau:
{
	"action": "<appointment | pet_log | pet_schedule>"
}`

	suggestionPromptEN = `As an AI assistant, your task is to interpret the user's request and determine the appropriate action to take. The possible actions are:
- "appointment": For scheduling a one-time appointment or event.
- "pet_log": For logging a one-time activity or event related to a pet.
- "pet_schedule": For setting up a recurring schedule (e.g., daily, weekly).
- "unknown": If the request doesn't match any of the above or is unclear.

Based on the user's description, identify the action type and extract the necessary parameters. Respond with a JSON object containing:

- "action": The determined action type.
- "parameters": A dictionary of extracted parameters relevant to the action.

Rules:
- For "appointment", extract "pet_name", "appointment_type", "date" (in YYYY-MM-DD format, using current year for relative dates), and "time" (in HH:MM format).
  Generate a title in format: "{pet_name}'s {appointment_type} Appointment"
- For "pet_log", extract "pet_name", "activity", "date" (in YYYY-MM-DD format), and "time" (in HH:MM format if mentioned).
  Generate a title in format: "{pet_name}'s {activity} Log"
- For "pet_schedule", extract "pet_name", "activity", "frequency" (e.g., "daily", "weekly"), and "start_date" (in YYYY-MM-DD format).
  Generate a title in format: "{pet_name}'s {frequency} {activity}"
- If the action is "unknown", provide a brief explanation in "parameters" under "reason".
- Always use the current time %s for any dates mentioned without a year.
- Always include a generated title in the parameters based on the rules above.

Example:
- User: "Log that Max had his walk today at 2:30 PM"
- Response:
{
	"action": "pet_log",
	"parameters": {
		"title": "Max's walk Log",
		"pet_name": "Max",
		"activity": "walk",
		"date": "%s",
		"time": "14:30"
	}
}

Analyze the following request and provide the JSON response:

Request: %s with action: %s`

	suggestionPromptVI = `Là một trợ lý AI, nhiệm vụ của bạn là phân tích yêu cầu của người dùng và xác định hành động phù hợp. Các hành động có thể là:
- "appointment": Để lên lịch hẹn hoặc sự kiện một lần.
- "pet_log": Để ghi lại hoạt động hoặc sự kiện một lần liên quan đến thú cưng.
- "pet_schedule": Để thiết lập lịch trình định kỳ (ví dụ: hàng ngày, hàng tuần).
- "unknown": Nếu yêu cầu không khớp với bất kỳ hành động nào ở trên hoặc không rõ ràng.

Dựa trên mô tả của người dùng, xác định loại hành động và trích xuất các tham số cần thiết. Trả về một đối tượng JSON chứa:

- "action": Loại hành động đã xác định.
- "parameters": Một từ điển chứa các tham số liên quan đến hành động.

Quy tắc:
- Đối với "appointment", trích xuất "pet_name", "appointment_type", "date" (định dạng YYYY-MM-DD, sử dụng năm hiện tại cho ngày tương đối), và "time" (định dạng HH:MM).
  Tạo tiêu đề theo định dạng: "Lịch hẹn {appointment_type} của {pet_name}"
- Đối với "pet_log", trích xuất "pet_name", "activity", "date" (định dạng YYYY-MM-DD), và "time" (định dạng HH:MM nếu được đề cập).
  Tạo tiêu đề theo định dạng: "Nhật ký {activity} của {pet_name}"
- Đối với "pet_schedule", trích xuất "pet_name", "activity", "frequency" (ví dụ: "hàng ngày", "hàng tuần"), và "start_date" (định dạng YYYY-MM-DD).
  Tạo tiêu đề theo định dạng: "Lịch trình {frequency} {activity} của {pet_name}"
- Nếu hành động là "unknown", cung cấp giải thích ngắn gọn trong "parameters" dưới "reason".
- Luôn sử dụng thời gian hiện tại %s cho bất kỳ ngày nào được đề cập mà không có năm.
- Luôn bao gồm tiêu đề được tạo trong các tham số dựa trên các quy tắc trên.

Ví dụ:
- Người dùng: "Ghi lại rằng Max đã đi dạo hôm nay lúc 2:30 chiều"
- Phản hồi:
{
	"action": "pet_log",
	"parameters": {
		"title": "Nhật ký đi dạo của Max",
		"pet_name": "Max",
		"activity": "đi dạo",
		"date": "%s",
		"time": "14:30"
	}
}

Phân tích yêu cầu sau và cung cấp phản hồi JSON:

Yêu cầu: %s với hành động: %s`
)

// Enhanced DetermineActionGemini function with language support
func DetermineActionGemini(ctx *gin.Context, config *util.Config, description string) (*ActionResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("action:%s", description)
	if cached, found := responseCache.Get(cacheKey); found {
		if response, ok := cached.(*ActionResponse); ok {
			logger.Info("Cache hit for action determination", zap.String("description", description))
			return response, nil
		}
	}

	// Detect language
	language := detectLanguage(description)
	logger.Info("Detected language", zap.String("language", language))

	// Acquire semaphore
	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("server busy: %v", err)
	}
	defer sem.Release(1)

	var actionResponse *ActionResponse
	var lastErr error

	// Retry logic
	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Second * time.Duration(attempt))
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(config.GoogleAPIKey))
		if err != nil {
			lastErr = err
			continue
		}

		model := client.GenerativeModel("gemini-2.0-flash")
		model.ResponseMIMEType = "application/json"

		// Select prompt based on language
		var prompt string
		if language == LanguageVietnamese {
			prompt = fmt.Sprintf(actionDeterminationPromptVI, description)
		} else {
			prompt = fmt.Sprintf(actionDeterminationPromptEN, description)
		}

		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(prompt)},
		}

		resp, err := model.GenerateContent(ctx, genai.Text(description))
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		responseText, err := extractResponseText(resp)
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		actionResponse, err = parseActionResponse(responseText)
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		// Validate response
		if err := validateActionResponse(actionResponse); err != nil {
			lastErr = err
			client.Close()
			continue
		}

		// Cache successful response
		responseCache.Set(cacheKey, actionResponse, cache.DefaultExpiration)
		client.Close()
		return actionResponse, nil
	}

	logger.Error("Failed to determine action after retries",
		zap.String("description", description),
		zap.Error(lastErr))
	return nil, fmt.Errorf("failed to determine action after %d retries: %v", MaxRetries, lastErr)
}

// Enhanced GenerateSuggestionGemini function with language support
func GenerateSuggestionGemini(ctx *gin.Context, config *util.Config, action, description string) (*BaseResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("suggestion:%s:%s", action, description)
	if cached, found := responseCache.Get(cacheKey); found {
		if response, ok := cached.(*BaseResponse); ok {
			logger.Info("Cache hit for suggestion",
				zap.String("action", action),
				zap.String("description", description))
			return response, nil
		}
	}

	// Detect language
	language := detectLanguage(description)
	logger.Info("Detected language", zap.String("language", language))

	// Acquire semaphore
	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("server busy: %v", err)
	}
	defer sem.Release(1)

	var baseResponse *BaseResponse
	var lastErr error

	// Retry logic
	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Second * time.Duration(attempt))
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(config.GoogleAPIKey))
		if err != nil {
			lastErr = err
			continue
		}

		model := client.GenerativeModel("gemini-2.0-flash")
		model.ResponseMIMEType = "application/json"

		// Select prompt based on language
		var prompt string
		if language == LanguageVietnamese {
			prompt = fmt.Sprintf(suggestionPromptVI,
				time.Now().Format("2006-01-02"),
				time.Now().Format("2006-01-02"),
				description,
				action)
		} else {
			prompt = fmt.Sprintf(suggestionPromptEN,
				time.Now().Format("2006-01-02"),
				time.Now().Format("2006-01-02"),
				description,
				action)
		}

		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(prompt)},
		}

		resp, err := model.GenerateContent(ctx, genai.Text(description))
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		responseText, err := extractResponseText(resp)
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		baseResponse, err = parseBaseResponse(responseText, action)
		if err != nil {
			lastErr = err
			client.Close()
			continue
		}

		// Validate response
		if err := validateBaseResponse(baseResponse); err != nil {
			lastErr = err
			client.Close()
			continue
		}

		// Cache successful response
		responseCache.Set(cacheKey, baseResponse, cache.DefaultExpiration)
		client.Close()
		return baseResponse, nil
	}

	logger.Error("Failed to generate suggestion after retries",
		zap.String("action", action),
		zap.String("description", description),
		zap.Error(lastErr))
	return nil, fmt.Errorf("failed to generate suggestion after %d retries: %v", MaxRetries, lastErr)
}

// Helper function to detect language
func detectLanguage(text string) string {
	// Simple language detection based on common Vietnamese characters
	vietnameseChars := []rune{'á', 'à', 'ả', 'ã', 'ạ', 'ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ', 'â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ',
		'đ', 'é', 'è', 'ẻ', 'ẽ', 'ẹ', 'ê', 'ế', 'ề', 'ể', 'ễ', 'ệ', 'í', 'ì', 'ỉ', 'ĩ', 'ị',
		'ó', 'ò', 'ỏ', 'õ', 'ọ', 'ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ', 'ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ',
		'ú', 'ù', 'ủ', 'ũ', 'ụ', 'ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự', 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ'}

	for _, char := range text {
		for _, vnChar := range vietnameseChars {
			if char == vnChar {
				return LanguageVietnamese
			}
		}
	}
	return LanguageEnglish
}

// Helper functions
func extractResponseText(resp *genai.GenerateContentResponse) (string, error) {
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText = string(txt)
		}
	}
	if responseText == "" {
		return "", fmt.Errorf("empty response from model")
	}
	return responseText, nil
}

func parseActionResponse(responseText string) (*ActionResponse, error) {
	var actionResponse ActionResponse
	if err := json.Unmarshal([]byte(responseText), &actionResponse); err != nil {
		// Try to extract JSON from the response if it's wrapped in other text
		start := strings.Index(responseText, "{")
		end := strings.LastIndex(responseText, "}")
		if start >= 0 && end > start {
			jsonStr := responseText[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &actionResponse); err != nil {
				return nil, fmt.Errorf("error parsing action response: %v", err)
			}
		} else {
			return nil, fmt.Errorf("error parsing action response: %v", err)
		}
	}
	return &actionResponse, nil
}

func parseBaseResponse(responseText string, action string) (*BaseResponse, error) {
	var baseResponse BaseResponse
	switch action {
	case ActionAppointment:
		var v AppointmentResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			return nil, fmt.Errorf("error parsing appointment response: %v", err)
		}
		baseResponse = BaseResponse{
			Action: v.Action,
			Parameters: map[string]interface{}{
				"title":            v.Parameters.Title,
				"pet_name":         v.Parameters.PetName,
				"appointment_type": v.Parameters.AppointmentType,
				"date":             v.Parameters.Date,
				"time":             v.Parameters.Time,
			},
		}
	case ActionPetLog:
		var v LogResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			return nil, fmt.Errorf("error parsing log response: %v", err)
		}
		baseResponse = BaseResponse{
			Action: v.Action,
			Parameters: map[string]interface{}{
				"title":    v.Parameters.Title,
				"pet_name": v.Parameters.PetName,
				"activity": v.Parameters.Activity,
				"date":     v.Parameters.Date,
				"time":     v.Parameters.Time,
				"notes":    v.Parameters.Notes,
			},
		}
	case ActionPetSchedule:
		var v ScheduleResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			return nil, fmt.Errorf("error parsing schedule response: %v", err)
		}
		baseResponse = BaseResponse{
			Action: v.Action,
			Parameters: map[string]interface{}{
				"title":    v.Parameters.Title,
				"pet_name": v.Parameters.PetName,
				"activity": v.Parameters.Activity,
				"date":     v.Parameters.Date,
				"time":     v.Parameters.Time,
				"notes":    v.Parameters.Notes,
			},
		}
	default:
		return nil, fmt.Errorf("unknown action type: %s", action)
	}
	return &baseResponse, nil
}

func validateActionResponse(response *ActionResponse) error {
	if response == nil {
		return fmt.Errorf("nil response")
	}
	if response.Action == "" {
		return fmt.Errorf("empty action in response")
	}
	switch response.Action {
	case ActionAppointment, ActionPetLog, ActionPetSchedule:
		return nil
	default:
		return fmt.Errorf("invalid action type: %s", response.Action)
	}
}

func validateBaseResponse(response *BaseResponse) error {
	if response == nil {
		return fmt.Errorf("nil response")
	}
	if response.Action == "" {
		return fmt.Errorf("empty action in response")
	}
	if response.Parameters == nil {
		return fmt.Errorf("nil parameters in response")
	}
	return nil
}
