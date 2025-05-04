package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/services"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
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
	geminiService *services.GeminiService

	conversationService *services.ConversationService
}

// NewChatHandler creates a new instance of ChatHandler
func NewChatHandler(geminiAPIKey string) *ChatHandler {
	// Create cache service that will be shared across all services
	cache := services.NewCacheService()

	// Create services with cache
	geminiService := services.NewGeminiService(geminiAPIKey)

	// Create new services for enhanced functionality
	conversationService := services.NewConversationService(cache)

	return &ChatHandler{
		geminiService: geminiService,

		conversationService: conversationService,
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

	// If no language specified, default to Vietnamese for Vietnamese system prompt
	if request.Language == "" {
		request.Language = "vi"
	}

	authPayload, err := middleware.GetAuthorizationPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	user, err := db.StoreDB.GetUser(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user", "details": err.Error()})
		return
	}

	userID := strconv.Itoa(int(user.ID))

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

		// convert into string

		// Create new conversation state
		conversationState = &models.ConversationState{
			ConversationID: conversationID,
			UserID:         userID,
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
				UserID:         userID,
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

	// Get previous messages for context
	if len(conversationState.Messages) > 1 {
		var previousMessages []string
		// Get the last 5 messages or all if less than 5
		startIdx := 0
		if len(conversationState.Messages) > 10 {
			startIdx = len(conversationState.Messages) - 10
		}

		for i := startIdx; i < len(conversationState.Messages); i++ {
			previousMessages = append(previousMessages, conversationState.Messages[i].Message)
		}

		// Remove the current message which was just added
		if len(previousMessages) > 0 {
			request.PreviousMessages = previousMessages[:len(previousMessages)-1]
		}
	}

	// Process the request with Gemini API
	response, err := h.geminiService.ProcessChatRequest(request)
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
