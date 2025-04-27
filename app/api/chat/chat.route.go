package chat

import (
	"context"
	"log"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type ChatRouteInterface interface {
	SetupRoutes(routerGroup middleware.RouterGroup)
}

// ChatApi implements the ChatRouteInterface
func NewChatApi(store db.Store, wsManager *websocket.WSClientManager) *ChatApi {
	service := NewChatService(store, wsManager)
	controller := NewChatController(service)
	return &ChatApi{controller: controller}
}

// SetupRoutes sets up all chat related routes
func (api *ChatApi) SetupRoutes(routerGroup middleware.RouterGroup) {
	chatRoutes := routerGroup.RouterDefault.Group("/")

	{
		// Public route for WebSocket connection
		chatRoutes.GET("/ws/messages", middleware.AuthMiddleware(token.TokenMaker), api.controller.HandleWebSocketChat)

		// Protected routes
		authRoutes := routerGroup.RouterAuth(chatRoutes)
		{
			// Conversations
			authRoutes.POST("/messages/conversations", api.controller.CreateConversation)
			authRoutes.GET("/messages/conversations", api.controller.GetConversations)
			authRoutes.GET("/messages/conversations/:id", api.controller.GetConversation)

			// Messages
			authRoutes.GET("/messages/conversations/:id/messages", api.controller.GetMessages)
			authRoutes.POST("/messages/messages", api.controller.SendMessage)
			authRoutes.POST("/messages/messages/read", api.controller.MarkMessagesAsRead)

			// Participants
			authRoutes.POST("/messages/conversations/:id/participants", api.controller.AddParticipants)
			authRoutes.DELETE("/messages/conversations/:id/participants/:userId", api.controller.RemoveParticipant)
			authRoutes.DELETE("/messages/conversations/:id/leave", api.controller.LeaveConversation)
		}
	}
}

// RegisterChatWebSocketHandlers registers handlers for chat-related WebSocket messages
func RegisterChatWebSocketHandlers(chatService *ChatService, wsManager *websocket.WSClientManager) {
	// Register a handler for chat messages
	wsManager.RegisterMessageTypeHandler("chat_message", func(message websocket.WebSocketMessage, clientID string, userID int64) {
		if userID <= 0 {
			log.Println("Cannot process chat message: Invalid user ID")
			return
		}

		// Process the chat message
		go chatService.HandleWebSocketChatMessage(message, clientID, userID)
	})

	// Register typing status handler
	wsManager.RegisterMessageTypeHandler("typing", func(message websocket.WebSocketMessage, clientID string, userID int64) {
		if userID <= 0 {
			return
		}

		// Process typing notification - just broadcast it directly
		if message.ConversationID > 0 {
			go chatService.broadcastTypingStatus(message.ConversationID, userID)
		}
	})

	// Register read receipts handler
	wsManager.RegisterMessageTypeHandler("read_receipt", func(message websocket.WebSocketMessage, clientID string, userID int64) {
		if userID <= 0 || message.MessageID <= 0 {
			return
		}

		// Mark message as read
		req := ReadMessageRequest{
			MessageIDs: []int64{message.MessageID},
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := chatService.MarkMessagesAsRead(ctx, req, userID)
			if err != nil {
				log.Printf("Error marking message as read via WebSocket: %v", err)
			}
		}()
	})
}

// RegisterRoutes registers all chat-related routes
func RegisterRoutes(routerGroup middleware.RouterGroup, store db.Store, wsManager *websocket.WSClientManager) {
	chatAPI := NewChatApi(store, wsManager)
	chatAPI.SetupRoutes(routerGroup)

	// Register WebSocket handlers
	service := NewChatService(store, wsManager)
	RegisterChatWebSocketHandlers(service, wsManager)
}
