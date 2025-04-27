package chat

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
)

type ChatController struct {
	service *ChatService
}

// NewChatController creates a new chat controller
func NewChatController(service *ChatService) *ChatController {
	return &ChatController{
		service: service,
	}
}

// getUserIDFromPayload gets the user ID from the username in the auth payload
func (c *ChatController) getUserIDFromPayload(ctx *gin.Context, authPayload *token.Payload) (int64, error) {
	// Get the user from the database based on username
	user, err := c.service.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	return user.ID, nil
}

// CreateConversation handles the creation of a new conversation
func (c *ChatController) CreateConversation(ctx *gin.Context) {
	var req CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the conversation
	conversation, err := c.service.CreateConversation(ctx, req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, conversation)
}

// GetConversations returns all conversations for the current user
func (c *ChatController) GetConversations(ctx *gin.Context) {
	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse pagination parameters
	var pagination PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	conversations, err := c.service.GetConversations(ctx, userID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, conversations)
}

// GetConversation returns a single conversation by ID
func (c *ChatController) GetConversation(ctx *gin.Context) {
	// Parse conversation ID from URL
	conversationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	conversation, err := c.service.GetConversation(ctx, conversationID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, conversation)
}

// GetMessages returns messages for a conversation
func (c *ChatController) GetMessages(ctx *gin.Context) {
	// Parse conversation ID from URL
	conversationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse pagination parameters
	var pagination PaginationParams
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	messages, err := c.service.GetMessages(ctx, conversationID, userID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

// SendMessage sends a new message to a conversation
func (c *ChatController) SendMessage(ctx *gin.Context) {
	// Parse request body
	var req SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	message, err := c.service.SendMessage(ctx, req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, message)
}

// MarkMessagesAsRead marks messages as read by the current user
func (c *ChatController) MarkMessagesAsRead(ctx *gin.Context) {
	// Parse request body
	var req ReadMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	err = c.service.MarkMessagesAsRead(ctx, req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// AddParticipants adds participants to a conversation
func (c *ChatController) AddParticipants(ctx *gin.Context) {
	// Parse conversation ID from URL
	conversationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Parse request body
	var req AddParticipantsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	err = c.service.AddParticipants(ctx, conversationID, req, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// RemoveParticipant removes a participant from a conversation
func (c *ChatController) RemoveParticipant(ctx *gin.Context) {
	// Parse conversation ID and user ID from URL
	conversationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	participantID, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	err = c.service.RemoveParticipant(ctx, conversationID, participantID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// LeaveConversation allows a user to leave a conversation
func (c *ChatController) LeaveConversation(ctx *gin.Context) {
	// Parse conversation ID from URL
	conversationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Get user ID from JWT token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	err = c.service.LeaveConversation(ctx, conversationID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// HandleWebSocketChat handles WebSocket connections for chat
func (c *ChatController) HandleWebSocketChat(ctx *gin.Context) {
	// Get user information for client ID
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user ID from username
	userID, err := c.getUserIDFromPayload(ctx, authPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set client ID in query parameters
	clientID := fmt.Sprintf("user_%s", authPayload.Username)
	ctx.Request.Header.Set("X-Client-ID", clientID)

	// Set user info in context for WebSocket handler
	ctx.Set("username", authPayload.Username)
	ctx.Set("userID", userID)

	// Pass to the WebSocket handler
	c.service.wsManager.HandleWebSocket(ctx)
}
