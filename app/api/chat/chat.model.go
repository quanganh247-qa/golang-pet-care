package chat

import (
	"encoding/json"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// ChatMessageType defines the types of messages that can be exchanged
type ChatMessageType string

const (
	ChatMessageTypeText  ChatMessageType = "text"
	ChatMessageTypeImage ChatMessageType = "image"
	ChatMessageTypeFile  ChatMessageType = "file"
)

// ConversationType defines the types of conversations
type ConversationType string

const (
	ConversationTypePrivate ConversationType = "private"
	ConversationTypeGroup   ConversationType = "group"
)

// WebSocketChatMessage represents a chat message sent over WebSocket
type WebSocketChatMessage struct {
	Type           string          `json:"type"`           // "chat_message"
	Action         string          `json:"action"`         // "new", "read", etc.
	MessageID      int64           `json:"messageId"`      // Set for existing messages
	ConversationID int64           `json:"conversationId"` // The conversation this message belongs to
	SenderID       int64           `json:"senderId"`       // User ID of sender
	Content        string          `json:"content"`        // Message content
	MessageType    ChatMessageType `json:"messageType"`    // text, image, file
	Metadata       json.RawMessage `json:"metadata"`       // Additional data
	CreatedAt      time.Time       `json:"createdAt"`      // When the message was created
	SenderUsername string          `json:"senderUsername"` // Username of the sender
	SenderName     string          `json:"senderName"`     // Full name of the sender
}

// ConversationResponse represents a conversation in the API
type ConversationResponse struct {
	ID           int64                 `json:"id"`
	Type         ConversationType      `json:"type"`
	Name         string                `json:"name,omitempty"`
	Participants []ParticipantResponse `json:"participants"`
	LastMessage  *MessageResponse      `json:"lastMessage,omitempty"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	UnreadCount  int64                 `json:"unreadCount"`
}

// MessageResponse represents a message in the API
type MessageResponse struct {
	ID             int64           `json:"id"`
	ConversationID int64           `json:"conversationId"`
	SenderID       int64           `json:"senderId"`
	Content        string          `json:"content"`
	MessageType    ChatMessageType `json:"messageType"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	SenderUsername string          `json:"senderUsername"`
	SenderName     string          `json:"senderName"`
	Read           bool            `json:"read"`
}

// ParticipantResponse represents a conversation participant in the API
type ParticipantResponse struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	IsAdmin  bool      `json:"isAdmin"`
	JoinedAt time.Time `json:"joinedAt"`
}

// CreateConversationRequest represents a request to create a new conversation
type CreateConversationRequest struct {
	Type           ConversationType `json:"type" binding:"required,oneof=private group"`
	Name           string           `json:"name"` // Required for group conversations
	ParticipantIDs []int64          `json:"participantIds" binding:"required,min=1"`
}

// SendMessageRequest represents a request to send a new message
type SendMessageRequest struct {
	ConversationID int64           `json:"conversationId" binding:"required"`
	Content        string          `json:"content" binding:"required"`
	MessageType    ChatMessageType `json:"messageType" binding:"required,oneof=text image file"`
	// Metadata       json.RawMessage `json:"metadata,omitempty"`
}

// ReadMessageRequest represents a request to mark messages as read
type ReadMessageRequest struct {
	MessageIDs []int64 `json:"messageIds" binding:"required,min=1"`
}

// AddParticipantsRequest represents a request to add participants to a conversation
type AddParticipantsRequest struct {
	ParticipantIDs []int64 `json:"participantIds" binding:"required,min=1"`
}

// PaginationParams represents pagination parameters for list endpoints
type PaginationParams struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

// GetDefaultPagination returns default pagination parameters
func GetDefaultPagination(p PaginationParams) PaginationParams {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p
}

// ChatUser represents a user in the chat system
type ChatUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// ChatParticipant extends ChatUser with conversation-specific info
type ChatParticipant struct {
	ChatUser
	IsAdmin  bool      `json:"isAdmin"`
	JoinedAt time.Time `json:"joinedAt"`
}

// ToParticipantResponse converts a ChatParticipant to ParticipantResponse
func (p ChatParticipant) ToParticipantResponse() ParticipantResponse {
	return ParticipantResponse{
		ID:       p.ID,
		Username: p.Username,
		Email:    p.Email,
		FullName: p.FullName,
		IsAdmin:  p.IsAdmin,
		JoinedAt: p.JoinedAt,
	}
}

// DbMessageToResponse converts a database message to a response message
func DbMessageToResponse(dbMessage db.GetMessagesByConversationRow, read bool) MessageResponse {
	return MessageResponse{
		ID:             dbMessage.ID,
		ConversationID: dbMessage.ConversationID,
		SenderID:       dbMessage.SenderID,
		Content:        dbMessage.Content,
		MessageType:    ChatMessageType(dbMessage.MessageType),
		Metadata:       dbMessage.Metadata,
		CreatedAt:      dbMessage.CreatedAt.Time,
		SenderUsername: dbMessage.SenderUsername,
		SenderName:     dbMessage.SenderName,
		Read:           read,
	}
}

// DbConversationToResponse converts a database conversation to a response conversation
func DbConversationToResponse(
	dbConv db.Conversation,
	participants []ParticipantResponse,
	lastMessage *MessageResponse,
	unreadCount int64,
) ConversationResponse {
	return ConversationResponse{
		ID:           dbConv.ID,
		Type:         ConversationType(dbConv.Type),
		Name:         dbConv.Name.String,
		Participants: participants,
		LastMessage:  lastMessage,
		CreatedAt:    dbConv.CreatedAt.Time,
		UpdatedAt:    dbConv.UpdatedAt.Time,
		UnreadCount:  unreadCount,
	}
}

// ChatApi defines the API structure for chat functionality
type ChatApi struct {
	controller *ChatController
}
