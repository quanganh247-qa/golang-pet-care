package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SuggestionsGenerator provides contextual follow-up questions and suggestions
type SuggestionsGenerator struct {
	gemini *GeminiService
	cache  *CacheService
}

// NewSuggestionsGenerator creates a new suggestions generator
func NewSuggestionsGenerator(gemini *GeminiService, cache *CacheService) *SuggestionsGenerator {
	return &SuggestionsGenerator{
		gemini: gemini,
		cache:  cache,
	}
}

// GenerateDrugInfoFollowUps generates follow-up questions for drug information
func (s *SuggestionsGenerator) GenerateDrugInfoFollowUps(drugName string, language string) ([]string, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("followups:drug:%s:%s", strings.ToLower(drugName), language)
	var cachedQuestions []string
	if s.cache.GetJSON(cacheKey, &cachedQuestions) && len(cachedQuestions) > 0 {
		return cachedQuestions, nil
	}

	// Define language-appropriate templates
	var templates []string
	if language == "vi" {
		templates = []string{
			"Tác dụng phụ phổ biến của %s là gì?",
			"%s có tương tác với thuốc nào khác không?",
			"Khi nào không nên dùng %s?",
			"Liều dùng thông thường của %s là bao nhiêu?",
			"Có cần kê đơn để mua %s không?",
		}
	} else {
		templates = []string{
			"What are the common side effects of %s?",
			"Does %s interact with other medications?",
			"When should I not take %s?",
			"What is the typical dosage for %s?",
			"Is %s available over the counter?",
		}
	}

	// Generate at least 3 follow-up questions using templates
	questions := make([]string, 0, 3)
	for i := 0; i < 3 && i < len(templates); i++ {
		questions = append(questions, fmt.Sprintf(templates[i], drugName))
	}

	// Try to generate more personalized questions with Gemini if available
	if s.gemini != nil {
		promptLang := language
		if promptLang != "vi" {
			promptLang = "en"
		}

		prompt := ""
		if promptLang == "vi" {
			prompt = fmt.Sprintf(`
				Hãy tạo 3 câu hỏi tiếp theo liên quan đến thuốc %s mà người dùng có thể muốn hỏi.
				Trả về dưới dạng mảng JSON chứa các câu hỏi.
				Mỗi câu hỏi phải ngắn gọn, cụ thể và khác với các mẫu sau:
				- "Tác dụng phụ phổ biến của %s là gì?"
				- "%s có tương tác với thuốc nào khác không?"
				- "Khi nào không nên dùng %s?"
				
				Chỉ trả về mảng JSON, không có text khác.
			`, drugName, drugName, drugName, drugName)
		} else {
			prompt = fmt.Sprintf(`
				Generate 3 follow-up questions about the medication %s that a user might want to ask.
				Return just a JSON array containing the questions.
				Each question should be concise, specific, and different from these templates:
				- "What are the common side effects of %s?"
				- "Does %s interact with other medications?"
				- "When should I not take %s?"
				
				Return only the JSON array with no other text.
			`, drugName, drugName, drugName, drugName)
		}

		response, err := s.gemini.sendRequest(prompt)
		if err == nil {
			// Extract JSON array from response using utility functions
			jsonStr := s.extractJSON(response)
			if jsonStr != "" {
				var geminiQuestions []string
				if err := json.Unmarshal([]byte(jsonStr), &geminiQuestions); err == nil && len(geminiQuestions) > 0 {
					// Replace template questions with generated questions
					questions = geminiQuestions
				}
			}
		}
	}

	// Cache the result for future use (1 day)
	if len(questions) > 0 {
		s.cache.Set(cacheKey, questions, 86400)
	}

	return questions, nil
}

// GenerateSideEffectFollowUps generates follow-up questions for side effect reports
func (s *SuggestionsGenerator) GenerateSideEffectFollowUps(drugName, sideEffect string, language string) ([]string, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("followups:sideeffect:%s:%s:%s",
		strings.ToLower(drugName), strings.ToLower(sideEffect), language)
	var cachedQuestions []string
	if s.cache.GetJSON(cacheKey, &cachedQuestions) && len(cachedQuestions) > 0 {
		return cachedQuestions, nil
	}

	// Define language-appropriate templates
	var templates []string
	if language == "vi" {
		templates = []string{
			"Tôi có nên ngừng dùng %s khi gặp %s không?",
			"Có cách nào giảm thiểu %s khi dùng %s?",
			"Khi nào tác dụng phụ %s trở nên nguy hiểm?",
			"Có thuốc thay thế nào cho %s ít gây %s hơn không?",
			"Bao lâu thì tác dụng phụ %s biến mất sau khi ngừng %s?",
		}
	} else {
		templates = []string{
			"Should I stop taking %s if I experience %s?",
			"Are there ways to reduce %s when taking %s?",
			"When does the side effect %s become dangerous?",
			"Are there alternative medications to %s that cause less %s?",
			"How long will %s persist after stopping %s?",
		}
	}

	// Generate follow-up questions using templates
	questions := make([]string, 0, 3)
	questions = append(questions, fmt.Sprintf(templates[0], drugName, sideEffect))
	questions = append(questions, fmt.Sprintf(templates[1], sideEffect, drugName))
	questions = append(questions, fmt.Sprintf(templates[2], sideEffect))

	// Try to generate more personalized questions with Gemini if available
	if s.gemini != nil {
		prompt := ""
		if language == "vi" {
			prompt = fmt.Sprintf(`
				Hãy tạo 3 câu hỏi tiếp theo về tác dụng phụ %s của thuốc %s mà người dùng có thể muốn hỏi.
				Trả về dưới dạng mảng JSON chứa các câu hỏi.
				Mỗi câu hỏi phải ngắn gọn và cụ thể.
				
				Chỉ trả về mảng JSON, không có text khác.
			`, sideEffect, drugName)
		} else {
			prompt = fmt.Sprintf(`
				Generate 3 follow-up questions about the side effect %s of medication %s that a user might want to ask.
				Return just a JSON array containing the questions.
				Each question should be concise and specific.
				
				Return only the JSON array with no other text.
			`, sideEffect, drugName)
		}

		response, err := s.gemini.sendRequest(prompt)
		if err == nil {
			jsonStr := s.extractJSON(response)
			if jsonStr != "" {
				var geminiQuestions []string
				if err := json.Unmarshal([]byte(jsonStr), &geminiQuestions); err == nil && len(geminiQuestions) > 0 {
					questions = geminiQuestions
				}
			}
		}
	}

	// Cache the result (1 day)
	if len(questions) > 0 {
		s.cache.Set(cacheKey, questions, 86400)
	}

	return questions, nil
}

// GenerateHealthTrendFollowUps generates follow-up questions for health trend analysis
func (s *SuggestionsGenerator) GenerateHealthTrendFollowUps(query, queryType string, language string) ([]string, error) {
	// Generate a unique cache key
	queryTrimmed := strings.ToLower(strings.TrimSpace(query))
	if len(queryTrimmed) > 50 {
		queryTrimmed = queryTrimmed[:50]
	}
	cacheKey := fmt.Sprintf("followups:health:%s:%s:%s",
		queryTrimmed, queryType, language)

	var cachedQuestions []string
	if s.cache.GetJSON(cacheKey, &cachedQuestions) && len(cachedQuestions) > 0 {
		return cachedQuestions, nil
	}

	// Default follow-up questions based on query type
	var questions []string

	if language == "vi" {
		if strings.Contains(queryType, "drug") {
			questions = []string{
				"Xu hướng báo cáo tác dụng phụ này đã thay đổi như thế nào trong 5 năm qua?",
				"Những nhóm tuổi nào bị ảnh hưởng nhiều nhất?",
				"So sánh số liệu này với các loại thuốc tương tự?",
			}
		} else if strings.Contains(queryType, "device") {
			questions = []string{
				"Thiết bị y tế này đã gây ra những sự cố nghiêm trọng nào?",
				"Có bao nhiêu báo cáo về thiết bị này trong năm ngoái?",
				"So sánh an toàn của thiết bị này với các thiết bị tương tự?",
			}
		} else if strings.Contains(queryType, "food") {
			questions = []string{
				"Những phản ứng phổ biến nhất với thực phẩm này là gì?",
				"Có bao nhiêu vụ thu hồi thực phẩm này trong 3 năm qua?",
				"Các nhóm dân số nào dễ bị ảnh hưởng nhất bởi sản phẩm này?",
			}
		} else {
			questions = []string{
				"Làm thế nào để xu hướng này thay đổi theo mùa?",
				"Có sự khác biệt đáng kể nào giữa nam và nữ không?",
				"Xu hướng này có thay đổi theo khu vực địa lý không?",
			}
		}
	} else {
		if strings.Contains(queryType, "drug") {
			questions = []string{
				"How has the trend in reporting these side effects changed over the past 5 years?",
				"Which age groups are most affected?",
				"How do these figures compare to similar medications?",
			}
		} else if strings.Contains(queryType, "device") {
			questions = []string{
				"What serious incidents has this medical device caused?",
				"How many reports were filed about this device in the past year?",
				"How does the safety of this device compare to similar devices?",
			}
		} else if strings.Contains(queryType, "food") {
			questions = []string{
				"What are the most common reactions to this food product?",
				"How many recalls of this food product were there in the past 3 years?",
				"Which demographic groups are most affected by this product?",
			}
		} else {
			questions = []string{
				"How does this trend vary seasonally?",
				"Is there any significant difference between males and females?",
				"Does this trend vary by geographic region?",
			}
		}
	}

	// Try to generate more relevant questions with Gemini
	if s.gemini != nil {
		prompt := ""
		if language == "vi" {
			prompt = fmt.Sprintf(`
				Dựa trên truy vấn của người dùng: "%s"
				Và loại dữ liệu: "%s"
				
				Hãy tạo 3 câu hỏi tiếp theo mà người dùng có thể muốn hỏi để phân tích sâu hơn.
				Trả về dưới dạng mảng JSON chứa các câu hỏi.
				Mỗi câu hỏi phải ngắn gọn, cụ thể và liên quan đến dữ liệu y tế.
				
				Chỉ trả về mảng JSON, không có text khác.
			`, query, queryType)
		} else {
			prompt = fmt.Sprintf(`
				Based on the user query: "%s"
				And data type: "%s"
				
				Generate 3 follow-up questions the user might want to ask to further analyze the data.
				Return just a JSON array containing the questions.
				Each question should be concise, specific, and related to health data.
				
				Return only the JSON array with no other text.
			`, query, queryType)
		}

		response, err := s.gemini.sendRequest(prompt)
		if err == nil {
			jsonStr := s.extractJSON(response)
			if jsonStr != "" {
				var geminiQuestions []string
				if err := json.Unmarshal([]byte(jsonStr), &geminiQuestions); err == nil && len(geminiQuestions) > 0 {
					questions = geminiQuestions
				}
			}
		}
	}

	// Cache the result (1 day)
	if len(questions) > 0 {
		s.cache.Set(cacheKey, questions, 86400)
	}

	return questions, nil
}

// extractJSON attempts to extract a JSON array or object from text
func (s *SuggestionsGenerator) extractJSON(text string) string {
	// Look for JSON array start and end
	startIdx := strings.Index(text, "[")
	endIdx := strings.LastIndex(text, "]")

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return text[startIdx : endIdx+1]
	}

	// Look for JSON object start and end
	startIdx = strings.Index(text, "{")
	endIdx = strings.LastIndex(text, "}")

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return text[startIdx : endIdx+1]
	}

	return ""
}
