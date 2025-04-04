package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/services"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request parameters"`
	Details string `json:"details,omitempty" example:"Field validation failed"`
}

// SuccessResponse represents a success response
// @Description Success response structure
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// ChatHandler manages chat-related HTTP requests
type ChatHandler struct {
	geminiService        *services.GeminiService
	openFDAService       *services.OpenFDAService
	drugInfoService      *services.DrugInfoService
	sideEffectService    *services.SideEffectService
	conversationService  *services.ConversationService
	suggestionsGenerator *services.SuggestionsGenerator
}

// NewChatHandler creates a new instance of ChatHandler
func NewChatHandler(geminiAPIKey string, openFDAAPIKey string) *ChatHandler {
	// Create cache service that will be shared across all services
	cache := services.NewCacheService()

	// Create services with cache
	geminiService := services.NewGeminiService(geminiAPIKey)
	openFDAService := services.NewOpenFDAService(openFDAAPIKey, cache)
	drugInfoService := services.NewDrugInfoService(openFDAService, geminiService, cache)
	sideEffectService := services.NewSideEffectService(openFDAService, geminiService)

	// Create new services for enhanced functionality
	conversationService := services.NewConversationService(cache)
	suggestionsGenerator := services.NewSuggestionsGenerator(geminiService, cache)

	return &ChatHandler{
		geminiService:        geminiService,
		openFDAService:       openFDAService,
		drugInfoService:      drugInfoService,
		sideEffectService:    sideEffectService,
		conversationService:  conversationService,
		suggestionsGenerator: suggestionsGenerator,
	}
}

// HandleChatRequest processes incoming chat requests
func (h *ChatHandler) HandleChatRequest(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate the request
	if request.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
		return
	}

	// If no BotType is specified, use HealthTrendBot as default
	if request.BotType == "" {
		request.BotType = models.HealthTrendBot
	}

	// If no language specified, default to English
	if request.Language == "" {
		request.Language = "en"
	}

	// Create or retrieve conversation state
	var conversationState *models.ConversationState
	if request.ConversationID == "" {
		// Generate a new conversation ID if none provided
		conversationID, err := generateConversationID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate conversation ID"})
			return
		}
		request.ConversationID = conversationID

		// Create new conversation state
		conversationState = &models.ConversationState{
			ConversationID: conversationID,
			UserID:         request.UserID,
			Language:       request.Language,
			LastActivity:   time.Now().Unix(),
			CurrentContext: models.ContextGeneral,
		}
	} else {
		// Load existing conversation state
		var err error
		conversationState, err = h.conversationService.GetConversation(request.ConversationID)
		if err != nil {
			// If conversation not found, create a new one
			conversationState = &models.ConversationState{
				ConversationID: request.ConversationID,
				UserID:         request.UserID,
				Language:       request.Language,
				LastActivity:   time.Now().Unix(),
				CurrentContext: models.ContextGeneral,
			}
		} else {
			// Update last activity timestamp
			conversationState.LastActivity = time.Now().Unix()

			// Update language if it changed
			if request.Language != "" && request.Language != conversationState.Language {
				conversationState.Language = request.Language
			}
		}
	}

	// Add the message to conversation history
	conversationState.Messages = append(conversationState.Messages, models.MessageEntry{
		Message:   request.Message,
		Timestamp: time.Now().Unix(),
		IsBot:     false,
	})

	// Determine the most appropriate bot type based on message content and conversation context
	botType := h.determineIntentAndBotType(request.Message, conversationState)
	if botType != "" {
		request.BotType = botType
	}

	// Process the request based on the chatbot type
	var response models.ChatResponse
	var err error

	switch request.BotType {
	case models.MediBot:
		response, err = h.handleMediBotRequest(c, request, conversationState)
	case models.SideEffectHelper:
		response, err = h.handleSideEffectRequest(c, request, conversationState)
	default: // HealthTrendBot or unknown types
		response, err = h.handleHealthTrendRequest(c, request, conversationState)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request", "details": err.Error()})
		return
	}

	// Add response to conversation history
	conversationState.Messages = append(conversationState.Messages, models.MessageEntry{
		Message:   response.Message,
		Timestamp: time.Now().Unix(),
		IsBot:     true,
	})

	// Update conversation state
	conversationState.SuggestedFollowUps = response.FollowUpQuestions
	h.conversationService.SaveConversation(conversationState)

	// Set conversation ID in response
	response.ConversationID = conversationState.ConversationID

	c.JSON(http.StatusOK, response)
}

// generateConversationID creates a unique conversation ID
func generateConversationID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// determineIntentAndBotType analyzes the message to determine the most appropriate bot type
func (h *ChatHandler) determineIntentAndBotType(message string, state *models.ConversationState) models.BotType {
	// If current context is already established, maintain it unless clear change in intent
	if state.CurrentContext == models.ContextDrugInfo {
		// Check if message still appears to be about drug info
		if strings.Contains(strings.ToLower(message), "side effect") ||
			strings.Contains(strings.ToLower(message), "adverse") ||
			strings.Contains(strings.ToLower(message), "report") {
			state.CurrentContext = models.ContextSideEffect
			return models.SideEffectHelper
		}
		return models.MediBot
	} else if state.CurrentContext == models.ContextSideEffect {
		// Check if user is asking about drug info
		if !strings.Contains(strings.ToLower(message), "side effect") &&
			!strings.Contains(strings.ToLower(message), "report") {
			state.CurrentContext = models.ContextDrugInfo
			return models.MediBot
		}
		return models.SideEffectHelper
	}

	// Check for clear drug information intent
	drugInformationIndicators := []string{"what is", "information", "tell me about", "drug", "medication", "medicine", "use", "dosage", "dose"}
	sideEffectIndicators := []string{"side effect", "adverse", "reaction", "report", "experienced", "suffering from"}

	drugInfoScore := 0
	sideEffectScore := 0

	messageLower := strings.ToLower(message)

	for _, indicator := range drugInformationIndicators {
		if strings.Contains(messageLower, indicator) {
			drugInfoScore++
		}
	}

	for _, indicator := range sideEffectIndicators {
		if strings.Contains(messageLower, indicator) {
			sideEffectScore++
		}
	}

	// Determine appropriate context and bot type
	if sideEffectScore > drugInfoScore {
		state.CurrentContext = models.ContextSideEffect
		return models.SideEffectHelper
	} else if drugInfoScore > 0 {
		state.CurrentContext = models.ContextDrugInfo
		return models.MediBot
	} else {
		state.CurrentContext = models.ContextHealthTrend
		return models.HealthTrendBot
	}
}

// handleMediBotRequest processes drug information requests
func (h *ChatHandler) handleMediBotRequest(c *gin.Context, request models.ChatRequest, state *models.ConversationState) (models.ChatResponse, error) {
	// Extract drug info directly using the new service
	drugInfo, err := h.drugInfoService.ExtractDrugInfo(request.Message)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to extract drug information: %w", err)
	}

	// Update conversation state with current drug
	state.CurrentDrug = drugInfo.DrugName

	// Determine priority level based on content
	priorityLevel := models.PriorityLow
	if drugInfo.Warnings != "" || len(drugInfo.Recalls) > 0 {
		priorityLevel = models.PriorityHigh
	} else if len(drugInfo.SideEffects) > 0 {
		priorityLevel = models.PriorityMedium
	}

	// Generate follow-up questions for this drug
	followUps, _ := h.suggestionsGenerator.GenerateDrugInfoFollowUps(drugInfo.DrugName, state.Language)

	// Create a user-friendly message with collapsible sections for mobile
	var message string

	// Add language-appropriate header
	if state.Language == "vi" {
		message = fmt.Sprintf("<h3>Thông tin về thuốc %s</h3>", drugInfo.DrugName)
	} else {
		message = fmt.Sprintf("<h3>Information about %s</h3>", drugInfo.DrugName)
	}

	// Prioritize critical warnings if present (show these first)
	if priorityLevel == models.PriorityHigh {
		message += generateWarningSection(drugInfo, state.Language)
	}

	// Add indications if available
	if drugInfo.Indications != "" {
		if state.Language == "vi" {
			message += "<details open><summary class='section-header'>Công dụng</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Indications + "</p></div></details>"
		} else {
			message += "<details open><summary class='section-header'>Indications</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Indications + "</p></div></details>"
		}
	}

	// Add side effects if available
	if len(drugInfo.SideEffects) > 0 {
		if state.Language == "vi" {
			message += "<details><summary class='section-header'>Tác dụng phụ phổ biến</summary><div class='section-content'><ul>"
		} else {
			message += "<details><summary class='section-header'>Common Side Effects</summary><div class='section-content'><ul>"
		}

		for _, effect := range drugInfo.SideEffects {
			message += "<li>" + effect + "</li>"
		}
		message += "</ul></div></details>"
	}

	// Add contraindications if available
	if drugInfo.Contraindications != "" {
		if state.Language == "vi" {
			message += "<details><summary class='section-header'>Chống chỉ định (Khi nào không nên dùng)</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Contraindications + "</p></div></details>"
		} else {
			message += "<details><summary class='section-header'>Contraindications</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Contraindications + "</p></div></details>"
		}
	}

	// Add non-critical warnings if available and not already shown
	if drugInfo.Warnings != "" && priorityLevel != models.PriorityHigh {
		if state.Language == "vi" {
			message += "<details><summary class='section-header'>Cảnh báo quan trọng</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Warnings + "</p></div></details>"
		} else {
			message += "<details><summary class='section-header'>Important Warnings</summary>"
			message += "<div class='section-content'><p>" + drugInfo.Warnings + "</p></div></details>"
		}
	}

	// Add recalls if available and not just the "no recalls" message
	if len(drugInfo.Recalls) > 0 && !strings.Contains(drugInfo.Recalls[0], "No active recalls") && priorityLevel != models.PriorityHigh {
		if state.Language == "vi" {
			message += "<details><summary class='section-header'>Thu hồi hoặc cảnh báo an toàn gần đây</summary><div class='section-content'><ul>"
		} else {
			message += "<details><summary class='section-header'>Recent Recalls or Safety Alerts</summary><div class='section-content'><ul>"
		}

		for _, recall := range drugInfo.Recalls {
			message += "<li>" + recall + "</li>"
		}
		message += "</ul></div></details>"
	}

	// Add source information
	if state.Language == "vi" {
		message += "<p><small>Nguồn dữ liệu: " + drugInfo.SourceDetails + "</small></p>"
	} else {
		message += "<p><small>Data source: " + drugInfo.SourceDetails + "</small></p>"
	}

	// Add follow-up questions suggestions
	if len(followUps) > 0 {
		if state.Language == "vi" {
			message += "<div class='follow-up-section'><p>Bạn có thể muốn hỏi:</p><ul class='follow-up-list'>"
		} else {
			message += "<div class='follow-up-section'><p>You might want to ask:</p><ul class='follow-up-list'>"
		}

		for _, question := range followUps {
			message += "<li class='follow-up-item'>" + question + "</li>"
		}
		message += "</ul></div>"
	}

	// Add medication images if available
	var medicationImages []models.MedicationImage
	if drugInfo.ImageURL != "" {
		medicationImages = append(medicationImages, models.MedicationImage{
			URL:         drugInfo.ImageURL,
			Description: drugInfo.DrugName,
			Source:      drugInfo.SourceDetails,
		})
	}

	// Construct the response with drug information
	response := models.ChatResponse{
		Message:           message,
		SourceDetails:     drugInfo.SourceDetails,
		DrugInfo:          drugInfo,
		BotType:           models.MediBot,
		FollowUpQuestions: followUps,
		MedicationImages:  medicationImages,
		PriorityLevel:     priorityLevel,
		Language:          state.Language,
	}

	return response, nil
}

// generateWarningSection creates a highlighted warning section for critical alerts
func generateWarningSection(drugInfo *models.DrugInfo, language string) string {
	var warning string

	if language == "vi" {
		warning = "<div class='critical-warning'><h4>⚠️ CẢNH BÁO QUAN TRỌNG</h4>"
	} else {
		warning = "<div class='critical-warning'><h4>⚠️ IMPORTANT WARNING</h4>"
	}

	if drugInfo.Warnings != "" {
		warning += "<p>" + drugInfo.Warnings + "</p>"
	}

	if len(drugInfo.Recalls) > 0 && !strings.Contains(drugInfo.Recalls[0], "No active recalls") {
		if language == "vi" {
			warning += "<h5>Thu hồi hoặc cảnh báo an toàn:</h5><ul>"
		} else {
			warning += "<h5>Recalls or safety alerts:</h5><ul>"
		}

		for _, recall := range drugInfo.Recalls {
			warning += "<li>" + recall + "</li>"
		}
		warning += "</ul>"
	}

	warning += "</div>"
	return warning
}

// handleSideEffectRequest processes side effect reporting requests
func (h *ChatHandler) handleSideEffectRequest(c *gin.Context, request models.ChatRequest, state *models.ConversationState) (models.ChatResponse, error) {
	// Extract side effect information
	sideEffectReport, err := h.sideEffectService.ExtractSideEffectReport(request.Message)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to extract side effect information: %w", err)
	}

	// Update conversation state with current drug
	state.CurrentDrug = sideEffectReport.DrugName

	// Generate follow-up questions
	followUps, _ := h.suggestionsGenerator.GenerateSideEffectFollowUps(sideEffectReport.DrugName, sideEffectReport.SideEffectName, state.Language)

	// Determine priority level based on severity
	priorityLevel := models.PriorityMedium
	if sideEffectReport.IsCommon {
		priorityLevel = models.PriorityLow
	}

	// Create a user-friendly message
	var message string

	if state.Language == "vi" {
		message = fmt.Sprintf("<h3>Báo cáo tác dụng phụ: %s do thuốc %s</h3>",
			sideEffectReport.SideEffectName, sideEffectReport.DrugName)
	} else {
		message = fmt.Sprintf("<h3>Side Effect Report: %s from %s</h3>",
			sideEffectReport.SideEffectName, sideEffectReport.DrugName)
	}

	// Add frequency information
	if state.Language == "vi" {
		if sideEffectReport.IsCommon {
			message += fmt.Sprintf("<p>Đây là tác dụng phụ <strong>phổ biến</strong> của thuốc này. Tần suất: %s</p>",
				sideEffectReport.Frequency)
		} else {
			message += fmt.Sprintf("<p>Đây là tác dụng phụ <strong>ít gặp</strong> của thuốc này. Tần suất: %s</p>",
				sideEffectReport.Frequency)
		}
	} else {
		if sideEffectReport.IsCommon {
			message += fmt.Sprintf("<p>This is a <strong>common</strong> side effect of this medication. Frequency: %s</p>",
				sideEffectReport.Frequency)
		} else {
			message += fmt.Sprintf("<p>This is a <strong>rare</strong> side effect of this medication. Frequency: %s</p>",
				sideEffectReport.Frequency)
		}
	}

	// Add reporting instructions
	if state.Language == "vi" {
		message += "<details open><summary class='section-header'>Cách báo cáo tác dụng phụ cho FDA</summary>"
		message += "<div class='section-content'><p>" + sideEffectReport.ReportingSteps + "</p>"
		message += fmt.Sprintf("<p>Đường dẫn trực tiếp: <a href='%s' target='_blank'>%s</a></p></div></details>",
			sideEffectReport.ReportingLink, "Báo cáo tác dụng phụ tới FDA")
	} else {
		message += "<details open><summary class='section-header'>How to Report Side Effects to FDA</summary>"
		message += "<div class='section-content'><p>" + sideEffectReport.ReportingSteps + "</p>"
		message += fmt.Sprintf("<p>Direct link: <a href='%s' target='_blank'>%s</a></p></div></details>",
			sideEffectReport.ReportingLink, "Report Side Effects to FDA")
	}

	// Add advice
	if state.Language == "vi" {
		message += "<details><summary class='section-header'>Lời khuyên</summary><div class='section-content'>"
		message += "<p>Khi bạn gặp tác dụng phụ của thuốc:</p><ul>"
		message += "<li>Liên hệ ngay với bác sĩ nếu tác dụng phụ nghiêm trọng</li>"
		message += "<li>Không tự ý ngừng thuốc mà không có chỉ định của bác sĩ</li>"
		message += "<li>Báo cáo tác dụng phụ để giúp FDA thu thập thông tin về an toàn thuốc</li>"
		message += "</ul></div></details>"
	} else {
		message += "<details><summary class='section-header'>Advice</summary><div class='section-content'>"
		message += "<p>When experiencing medication side effects:</p><ul>"
		message += "<li>Contact your doctor immediately if side effects are severe</li>"
		message += "<li>Do not stop taking medication without medical guidance</li>"
		message += "<li>Report side effects to help FDA collect drug safety information</li>"
		message += "</ul></div></details>"
	}

	// Add follow-up questions suggestions
	if len(followUps) > 0 {
		if state.Language == "vi" {
			message += "<div class='follow-up-section'><p>Bạn có thể muốn hỏi:</p><ul class='follow-up-list'>"
		} else {
			message += "<div class='follow-up-section'><p>You might want to ask:</p><ul class='follow-up-list'>"
		}

		for _, question := range followUps {
			message += "<li class='follow-up-item'>" + question + "</li>"
		}
		message += "</ul></div>"
	}

	// Construct the response
	response := models.ChatResponse{
		Message:           message,
		SideEffectReport:  sideEffectReport,
		BotType:           models.SideEffectHelper,
		FollowUpQuestions: followUps,
		PriorityLevel:     priorityLevel,
		Language:          state.Language,
	}

	return response, nil
}

// handleHealthTrendRequest processes health trend analysis requests
func (h *ChatHandler) handleHealthTrendRequest(c *gin.Context, request models.ChatRequest, state *models.ConversationState) (models.ChatResponse, error) {
	// Step 1: Process the user query with Gemini to understand intent and extract key parameters
	geminiResponse, err := h.geminiService.ProcessQuery(request.Message)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to process query: %w", err)
	}

	// Step 2: Use the processed query to fetch data from OpenFDA
	fdaData, fdaQueryType, err := h.openFDAService.FetchDataBasedOnIntent(geminiResponse)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to fetch health data: %w", err)
	}

	// Step 3: Get Gemini to analyze the FDA data and generate a response
	analysis, err := h.geminiService.AnalyzeHealthData(request.Message, fdaData, fdaQueryType)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to analyze health data: %w", err)
	}

	// Generate follow-up questions
	followUps, _ := h.suggestionsGenerator.GenerateHealthTrendFollowUps(request.Message, fdaQueryType, state.Language)

	// Add appropriate language to the response
	if state.Language == "vi" && !strings.Contains(analysis.Summary, "Nguồn dữ liệu:") {
		// Add Vietnamese source attribution if not already present
		analysis.Summary += "<p><small>Nguồn dữ liệu: OpenFDA API</small></p>"

		// Add follow-up suggestions in Vietnamese
		if len(followUps) > 0 {
			var followUpSection string
			followUpSection = "<div class='follow-up-section'><p>Bạn có thể muốn hỏi:</p><ul class='follow-up-list'>"
			for _, question := range followUps {
				followUpSection += "<li class='follow-up-item'>" + question + "</li>"
			}
			followUpSection += "</ul></div>"
			analysis.Summary += followUpSection
		}
	} else if state.Language == "en" && !strings.Contains(analysis.Summary, "Data source:") {
		// Add English source attribution if not already present
		analysis.Summary += "<p><small>Data source: OpenFDA API</small></p>"

		// Add follow-up suggestions in English
		if len(followUps) > 0 {
			var followUpSection string
			followUpSection = "<div class='follow-up-section'><p>You might want to ask:</p><ul class='follow-up-list'>"
			for _, question := range followUps {
				followUpSection += "<li class='follow-up-item'>" + question + "</li>"
			}
			followUpSection += "</ul></div>"
			analysis.Summary += followUpSection
		}
	}

	// Construct the response for health trend analysis
	response := models.ChatResponse{
		Message:           analysis.Summary,
		Data:              fdaData,
		ChartData:         analysis.ChartData,
		ChartType:         analysis.ChartType,
		ChartTitle:        analysis.ChartTitle,
		SourceDetails:     analysis.SourceDetails,
		BotType:           models.HealthTrendBot,
		FollowUpQuestions: followUps,
		Language:          state.Language,
	}

	return response, nil
}

// HandleOptionsRequest handles preflight CORS requests
func (h *ChatHandler) HandleOptionsRequest(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Status(http.StatusOK)
}

// ListConversations returns a list of conversations for the authenticated user
func (h *ChatHandler) ListConversations(c *gin.Context) {
	// Get user from auth middleware
	authPayload, err := middleware.GetAuthorizationPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	// Get limit parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get conversations for this user
	conversations, err := h.conversationService.GetRecentConversations(authPayload.Username, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// GetConversation returns a single conversation by ID
func (h *ChatHandler) GetConversation(c *gin.Context) {
	// Get user from auth middleware
	authPayload, err := middleware.GetAuthorizationPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	// Get conversation
	conversation, err := h.conversationService.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found", "details": err.Error()})
		return
	}

	// Verify user owns this conversation
	if conversation.UserID != "" && conversation.UserID != authPayload.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this conversation"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

// DeleteConversation deletes a conversation by ID
func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	// Get user from auth middleware
	authPayload, err := middleware.GetAuthorizationPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation ID is required"})
		return
	}

	// Get conversation first to verify ownership
	conversation, err := h.conversationService.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found", "details": err.Error()})
		return
	}

	// Verify user owns this conversation
	if conversation.UserID != "" && conversation.UserID != authPayload.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this conversation"})
		return
	}

	// Delete the conversation
	err = h.conversationService.DeleteConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}

// GetDrugFollowUpQuestions generates follow-up questions for a drug
func (h *ChatHandler) GetDrugFollowUpQuestions(c *gin.Context) {
	// Get drug name from path
	drugName := c.Param("drug_name")
	if drugName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Drug name is required"})
		return
	}

	// Get language preference
	language := c.DefaultQuery("language", "en")

	// Generate follow-up questions
	questions, err := h.suggestionsGenerator.GenerateDrugInfoFollowUps(drugName, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate follow-up questions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"drugName":  drugName,
		"language":  language,
		"questions": questions,
	})
}
