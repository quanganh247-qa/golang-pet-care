package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/services"
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
	geminiService     *services.GeminiService
	openFDAService    *services.OpenFDAService
	drugInfoService   *services.DrugInfoService
	sideEffectService *services.SideEffectService
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

	return &ChatHandler{
		geminiService:     geminiService,
		openFDAService:    openFDAService,
		drugInfoService:   drugInfoService,
		sideEffectService: sideEffectService,
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

	// log.Printf("Received chat request for %s: %s", request.BotType, request.Message)

	// Process the request based on the chatbot type
	switch request.BotType {
	case models.MediBot:
		h.handleMediBotRequest(c, request)
	case models.SideEffectHelper:
		h.handleSideEffectRequest(c, request)
	default: // HealthTrendBot or unknown types
		h.handleHealthTrendRequest(c, request)
	}
}

// handleMediBotRequest processes drug information requests
func (h *ChatHandler) handleMediBotRequest(c *gin.Context, request models.ChatRequest) {
	// Extract drug info directly using the new service
	drugInfo, err := h.drugInfoService.ExtractDrugInfo(request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract drug information", "details": err.Error()})
		return
	}

	// Create a user-friendly message
	message := fmt.Sprintf("<h3>Thông tin về thuốc %s</h3>", drugInfo.DrugName)

	// Add indications if available
	if drugInfo.Indications != "" {
		message += "<h4>Công dụng</h4>"
		message += "<p>" + drugInfo.Indications + "</p>"
	}

	// Add side effects if available
	if len(drugInfo.SideEffects) > 0 {
		message += "<h4>Tác dụng phụ phổ biến</h4>"
		message += "<ul>"
		for _, effect := range drugInfo.SideEffects {
			message += "<li>" + effect + "</li>"
		}
		message += "</ul>"
	}

	// Add contraindications if available
	if drugInfo.Contraindications != "" {
		message += "<h4>Chống chỉ định (Khi nào không nên dùng)</h4>"
		message += "<p>" + drugInfo.Contraindications + "</p>"
	}

	// Add warnings if available
	if drugInfo.Warnings != "" {
		message += "<h4>Cảnh báo quan trọng</h4>"
		message += "<p>" + drugInfo.Warnings + "</p>"
	}

	// Add recalls if available and not just the "no recalls" message
	if len(drugInfo.Recalls) > 0 && !strings.Contains(drugInfo.Recalls[0], "No active recalls") {
		message += "<h4>Thu hồi hoặc cảnh báo an toàn gần đây</h4>"
		message += "<ul>"
		for _, recall := range drugInfo.Recalls {
			message += "<li>" + recall + "</li>"
		}
		message += "</ul>"
	}

	// Add source information
	message += "<p><small>Nguồn dữ liệu: " + drugInfo.SourceDetails + "</small></p>"

	// Construct the response with drug information
	response := models.ChatResponse{
		Message:       message,
		SourceDetails: drugInfo.SourceDetails,
		DrugInfo:      drugInfo,
		BotType:       models.MediBot,
	}

	c.JSON(http.StatusOK, response)
}

// handleSideEffectRequest processes side effect reporting requests
func (h *ChatHandler) handleSideEffectRequest(c *gin.Context, request models.ChatRequest) {
	// Extract side effect information
	sideEffectReport, err := h.sideEffectService.ExtractSideEffectReport(request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract side effect information", "details": err.Error()})
		return
	}

	// Create a user-friendly message
	message := fmt.Sprintf("<h3>Báo cáo tác dụng phụ: %s do thuốc %s</h3>",
		sideEffectReport.SideEffectName, sideEffectReport.DrugName)

	// Add frequency information
	if sideEffectReport.IsCommon {
		message += fmt.Sprintf("<p>Đây là tác dụng phụ <strong>phổ biến</strong> của thuốc này. Tần suất: %s</p>",
			sideEffectReport.Frequency)
	} else {
		message += fmt.Sprintf("<p>Đây là tác dụng phụ <strong>ít gặp</strong> của thuốc này. Tần suất: %s</p>",
			sideEffectReport.Frequency)
	}

	// Add reporting instructions
	message += "<h4>Cách báo cáo tác dụng phụ cho FDA</h4>"
	message += "<p>" + sideEffectReport.ReportingSteps + "</p>"

	// Add reporting link
	message += fmt.Sprintf("<p>Đường dẫn trực tiếp: <a href='%s' target='_blank'>%s</a></p>",
		sideEffectReport.ReportingLink, "Báo cáo tác dụng phụ tới FDA")

	// Add advice
	message += "<h4>Lời khuyên</h4>"
	message += "<p>Khi bạn gặp tác dụng phụ của thuốc:</p>"
	message += "<ul>"
	message += "<li>Liên hệ ngay với bác sĩ nếu tác dụng phụ nghiêm trọng</li>"
	message += "<li>Không tự ý ngừng thuốc mà không có chỉ định của bác sĩ</li>"
	message += "<li>Báo cáo tác dụng phụ để giúp FDA thu thập thông tin về an toàn thuốc</li>"
	message += "</ul>"

	// Construct the response
	response := models.ChatResponse{
		Message:          message,
		SideEffectReport: sideEffectReport,
		BotType:          models.SideEffectHelper,
	}

	c.JSON(http.StatusOK, response)
}

// handleHealthTrendRequest processes health trend analysis requests
func (h *ChatHandler) handleHealthTrendRequest(c *gin.Context, request models.ChatRequest) {
	// Step 1: Process the user query with Gemini to understand intent and extract key parameters
	geminiResponse, err := h.geminiService.ProcessQuery(request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query", "details": err.Error()})
		return
	}

	// Step 2: Use the processed query to fetch data from OpenFDA
	fdaData, fdaQueryType, err := h.openFDAService.FetchDataBasedOnIntent(geminiResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch health data", "details": err.Error()})
		return
	}

	// Step 3: Get Gemini to analyze the FDA data and generate a response
	analysis, err := h.geminiService.AnalyzeHealthData(request.Message, fdaData, fdaQueryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze health data", "details": err.Error()})
		return
	}

	// Construct the response for health trend analysis
	response := models.ChatResponse{
		Message:       analysis.Summary,
		Data:          fdaData,
		ChartData:     analysis.ChartData,
		ChartType:     analysis.ChartType,
		ChartTitle:    analysis.ChartTitle,
		SourceDetails: analysis.SourceDetails,
		BotType:       models.HealthTrendBot,
	}

	c.JSON(http.StatusOK, response)
}
