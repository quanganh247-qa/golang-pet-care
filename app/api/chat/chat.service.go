package chat

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type ChatService struct {
	store     db.Store
	wsManager *websocket.WSClientManager
}

// NewChatService creates a new chat service
func NewChatService(store db.Store, wsManager *websocket.WSClientManager) *ChatService {
	return &ChatService{
		store:     store,
		wsManager: wsManager,
	}
}

// CreateConversation creates a new conversation
func (s *ChatService) CreateConversation(ctx context.Context, req CreateConversationRequest, creatorID int64) (*ConversationResponse, error) {
	// For private conversations, validate that there are exactly 2 participants
	if req.Type == ConversationTypePrivate && len(req.ParticipantIDs) != 1 {
		return nil, fmt.Errorf("private conversations must have exactly 1 participant besides the creator")
	}

	// For group conversations, validate that there is a name
	if req.Type == ConversationTypeGroup && req.Name == "" {
		return nil, fmt.Errorf("group conversations must have a name")
	}

	var conversation db.Conversation
	var participants []ParticipantResponse

	err := s.store.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Create the conversation
		var dbName pgtype.Text
		if req.Name != "" {
			dbName = pgtype.Text{
				String: req.Name,
				Valid:  true,
			}
		}

		var err error
		conversation, err = q.CreateConversation(ctx, db.CreateConversationParams{
			Type: string(req.Type),
			Name: dbName,
		})
		if err != nil {
			return fmt.Errorf("failed to create conversation: %w", err)
		}

		// Add creator as participant and admin
		_, err = q.AddParticipantToConversation(ctx, db.AddParticipantToConversationParams{
			ConversationID: conversation.ID,
			UserID:         creatorID,
			IsAdmin:        true,
		})
		if err != nil {
			return fmt.Errorf("failed to add creator as participant: %w", err)
		}

		// Add other participants
		for _, participantID := range req.ParticipantIDs {
			user, err := q.GetDoctor(ctx, participantID)
			if err != nil {
				return fmt.Errorf("failed to get user %d: %w", participantID, err)
			}
			_, err = q.AddParticipantToConversation(ctx, db.AddParticipantToConversationParams{
				ConversationID: conversation.ID,
				UserID:         user.UserID,
				IsAdmin:        false,
			})
			if err != nil {
				return fmt.Errorf("failed to add participant %d: %w", participantID, err)
			}
		}

		// Get participants to include in response
		dbParticipants, err := q.GetConversationParticipants(ctx, conversation.ID)
		if err != nil {
			return fmt.Errorf("failed to get participants: %w", err)
		}

		// Convert DB participants to API response format
		participants = make([]ParticipantResponse, 0, len(dbParticipants))
		for _, p := range dbParticipants {
			participants = append(participants, ParticipantResponse{
				ID:       p.ID,
				Username: p.Username,
				Email:    p.Email,
				FullName: p.FullName,
				IsAdmin:  p.IsAdmin,
				JoinedAt: p.JoinedAt.Time,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	_, err = s.store.CreateMessage(ctx, db.CreateMessageParams{
		ConversationID: conversation.ID,
		SenderID:       creatorID,
		Content:        "New conversation created",
		MessageType:    string(ChatMessageTypeText),
		Metadata:       nil,
	})
	if err != nil {
		return nil, err
	}

	// Update conversation timestamp
	err = s.store.UpdateConversationTimestamp(ctx, conversation.ID)
	if err != nil {
		log.Printf("Error updating conversation timestamp: %v", err)
		// Non-critical error, continue
	}

	// Create response with empty last message
	response := DbConversationToResponse(conversation, participants, nil, 0)

	// Notify other participants about the new conversation via WebSocket
	s.notifyParticipantsNewConversation(response)

	return &response, nil
}

// GetConversations gets all conversations for a user
func (s *ChatService) GetConversations(ctx context.Context, userID int64, pagination PaginationParams) ([]ConversationResponse, error) {

	dbConversations, err := s.store.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Build response
	conversations := make([]ConversationResponse, 0, len(dbConversations))
	for _, dbConv := range dbConversations {
		// Get participants
		dbParticipants, err := s.store.GetConversationParticipants(ctx, dbConv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get participants for conversation %d: %w", dbConv.ID, err)
		}

		participants := make([]ParticipantResponse, 0, len(dbParticipants))
		for _, p := range dbParticipants {
			participants = append(participants, ParticipantResponse{
				ID:       p.ID,
				Username: p.Username,
				Email:    p.Email,
				FullName: p.FullName,
				IsAdmin:  p.IsAdmin,
				JoinedAt: p.JoinedAt.Time,
			})
		}

		// Get last message
		var lastMessage *MessageResponse
		lastDbMessage, err := s.store.GetConversationLastMessage(ctx, dbConv.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to get last message for conversation %d: %w", dbConv.ID, err)
			}
		} else {
			// Check if this message has been read by the current user
			messageRead := false
			// This would require a new query to check if the message has been read
			// For now, assume unread for simplicity

			dbMessage := db.GetMessagesByConversationRow{
				ID:             lastDbMessage.ID,
				ConversationID: lastDbMessage.ConversationID,
				SenderID:       lastDbMessage.SenderID,
				Content:        lastDbMessage.Content,
				MessageType:    lastDbMessage.MessageType,
				Metadata:       lastDbMessage.Metadata,
				CreatedAt:      lastDbMessage.CreatedAt,
				SenderUsername: lastDbMessage.SenderUsername,
				SenderName:     lastDbMessage.SenderName,
			}

			msg := DbMessageToResponse(dbMessage, messageRead)
			lastMessage = &msg
		}

		// Get unread count
		unreadCount, err := s.store.GetUnreadMessageCount(ctx, db.GetUnreadMessageCountParams{
			ConversationID: dbConv.ID,
			UserID:         userID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get unread count for conversation %d: %w", dbConv.ID, err)
		}

		// Add to result
		conversations = append(conversations, DbConversationToResponse(
			dbConv,
			participants,
			lastMessage,
			unreadCount,
		))
	}

	return conversations, nil
}

// GetConversation gets a single conversation by ID
func (s *ChatService) GetConversation(ctx context.Context, conversationID int64, userID int64) (*ConversationResponse, error) {
	// Check if user is a participant
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return nil, fmt.Errorf("user is not a member of this conversation")
	}

	// Get conversation
	dbConv, err := s.store.GetConversation(ctx, conversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get participants
	dbParticipants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	participants := make([]ParticipantResponse, 0, len(dbParticipants))
	for _, p := range dbParticipants {
		participants = append(participants, ParticipantResponse{
			ID:       p.ID,
			Username: p.Username,
			Email:    p.Email,
			FullName: p.FullName,
			IsAdmin:  p.IsAdmin,
			JoinedAt: p.JoinedAt.Time,
		})
	}

	// Get last message
	var lastMessage *MessageResponse
	lastDbMessage, err := s.store.GetConversationLastMessage(ctx, conversationID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get last message: %w", err)
		}
	} else {
		// Check if this message has been read by the current user
		messageRead := false // Simplified
		dbMessage := db.GetMessagesByConversationRow{
			ID:             lastDbMessage.ID,
			ConversationID: lastDbMessage.ConversationID,
			SenderID:       lastDbMessage.SenderID,
			Content:        lastDbMessage.Content,
			MessageType:    lastDbMessage.MessageType,
			Metadata:       lastDbMessage.Metadata,
			CreatedAt:      lastDbMessage.CreatedAt,
			SenderUsername: lastDbMessage.SenderUsername,
			SenderName:     lastDbMessage.SenderName,
		}
		msg := DbMessageToResponse(dbMessage, messageRead)
		lastMessage = &msg
	}

	// Get unread count
	unreadCount, err := s.store.GetUnreadMessageCount(ctx, db.GetUnreadMessageCountParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}

	// Create response
	response := DbConversationToResponse(dbConv, participants, lastMessage, unreadCount)
	return &response, nil
}

// GetMessages gets messages for a conversation
func (s *ChatService) GetMessages(ctx context.Context, conversationID int64, userID int64, pagination PaginationParams) ([]MessageResponse, error) {
	// Check if user is a participant
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return nil, fmt.Errorf("user is not a member of this conversation")
	}

	// Get messages
	pagination = GetDefaultPagination(pagination)
	dbMessages, err := s.store.GetMessagesByConversation(ctx, db.GetMessagesByConversationParams{
		ConversationID: conversationID,
		Limit:          int32(pagination.PageSize),
		Offset:         int32((pagination.Page - 1) * pagination.PageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Convert to response format
	messages := make([]MessageResponse, 0, len(dbMessages))
	for _, msg := range dbMessages {
		// In a real implementation, you'd check if this message has been read by the user
		// For simplicity, we'll assume all messages are unread
		messages = append(messages, DbMessageToResponse(msg, false))
	}

	return messages, nil
}

// SendMessage sends a new message to a conversation
func (s *ChatService) SendMessage(ctx context.Context, req SendMessageRequest, senderID int64) (*MessageResponse, error) {
	// Check if user is a participant
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: req.ConversationID,
		UserID:         senderID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return nil, fmt.Errorf("user is not a member of this conversation")
	}

	// // Create metadata if none provided
	// var metadata sql.NullString
	// if req.Metadata != nil {
	// 	metadata = sql.NullString{
	// 		String: string(req.Metadata),
	// 		Valid:  true,
	// 	}
	// }

	// Create the message
	dbMessage, err := s.store.CreateMessage(ctx, db.CreateMessageParams{
		ConversationID: req.ConversationID,
		SenderID:       senderID,
		Content:        req.Content,
		MessageType:    string(req.MessageType),
		// Metadata:       []byte(metadata.String),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update conversation timestamp
	err = s.store.UpdateConversationTimestamp(ctx, req.ConversationID)
	if err != nil {
		log.Printf("Error updating conversation timestamp: %v", err)
		// Non-critical error, continue
	}

	// Get sender info for the response
	user, err := s.store.GetUserByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender info: %w", err)
	}

	addedUSer := db.User{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	// Broadcast message to all participants
	s.broadcastMessageToParticipants(ctx, req.ConversationID, dbMessage, addedUSer)

	// Get the full message details for response
	messageResponse := MessageResponse{
		ID:             dbMessage.ID,
		ConversationID: dbMessage.ConversationID,
		SenderID:       dbMessage.SenderID,
		Content:        dbMessage.Content,
		MessageType:    ChatMessageType(dbMessage.MessageType),
		// Metadata:       req.Metadata,
		CreatedAt:      dbMessage.CreatedAt.Time,
		SenderUsername: user.Username,
		SenderName:     user.FullName,
		Read:           false,
	}

	return &messageResponse, nil
}

// MarkMessagesAsRead marks messages as read by a user
func (s *ChatService) MarkMessagesAsRead(ctx context.Context, req ReadMessageRequest, userID int64) error {
	for _, messageID := range req.MessageIDs {
		_, err := s.store.MarkMessageAsRead(ctx, db.MarkMessageAsReadParams{
			MessageID: messageID,
			UserID:    userID,
		})
		if err != nil {
			return fmt.Errorf("failed to mark message %d as read: %w", messageID, err)
		}
	}

	// Notify other participants about read status
	// s.notifyMessagesRead(ctx, req.MessageIDs, userID)

	return nil
}

// AddParticipants adds participants to a conversation
func (s *ChatService) AddParticipants(ctx context.Context, conversationID int64, req AddParticipantsRequest, userID int64) error {
	// Check if user is an admin in this conversation
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return fmt.Errorf("user is not a member of this conversation")
	}

	// For simplicity, we're not checking if the user is an admin,
	// but in a real implementation you would add that check

	// Add each participant
	for _, participantID := range req.ParticipantIDs {
		_, err := s.store.AddParticipantToConversation(ctx, db.AddParticipantToConversationParams{
			ConversationID: conversationID,
			UserID:         participantID,
			IsAdmin:        false,
		})
		if err != nil {
			return fmt.Errorf("failed to add participant %d: %w", participantID, err)
		}
	}

	// Notify current participants about new members
	s.notifyParticipantsAdded(ctx, conversationID, req.ParticipantIDs)

	return nil
}

// RemoveParticipant removes a participant from a conversation
func (s *ChatService) RemoveParticipant(ctx context.Context, conversationID int64, participantID int64, userID int64) error {
	// Check if user is an admin in this conversation
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return fmt.Errorf("user is not a member of this conversation")
	}

	// For simplicity, we're not checking if the user is an admin,
	// but in a real implementation you would add that check

	// Remove the participant
	err = s.store.RemoveParticipantFromConversation(ctx, db.RemoveParticipantFromConversationParams{
		ConversationID: conversationID,
		UserID:         participantID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	// Notify participants about member removal
	s.notifyParticipantRemoved(ctx, conversationID, participantID)

	return nil
}

// LeaveConversation allows a user to leave a conversation
func (s *ChatService) LeaveConversation(ctx context.Context, conversationID int64, userID int64) error {
	// Check if user is in this conversation
	isMember, err := s.store.IsUserInConversation(ctx, db.IsUserInConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return fmt.Errorf("failed to check if user is in conversation: %w", err)
	}

	if !isMember {
		return fmt.Errorf("user is not a member of this conversation")
	}

	// Remove the user
	err = s.store.RemoveParticipantFromConversation(ctx, db.RemoveParticipantFromConversationParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return fmt.Errorf("failed to leave conversation: %w", err)
	}

	// Notify other participants
	s.notifyParticipantRemoved(ctx, conversationID, userID)

	return nil
}

// HandleWebSocketChatMessage handles incoming WebSocket chat messages
func (s *ChatService) HandleWebSocketChatMessage(message websocket.WebSocketMessage, clientID string, userID int64) {
	// Parse WebSocket message into chat message
	var chatMsg WebSocketChatMessage
	if err := json.Unmarshal([]byte(message.Message), &chatMsg); err != nil {
		log.Printf("Error parsing chat message: %v", err)
		return
	}

	// Handle different action types
	switch chatMsg.Action {
	case "send":
		// Send a new message
		req := SendMessageRequest{
			ConversationID: chatMsg.ConversationID,
			Content:        chatMsg.Content,
			MessageType:    chatMsg.MessageType,
			// Metadata:       chatMsg.Metadata,
		}

		_, err := s.SendMessage(context.Background(), req, userID)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			// Send error response to the client
			s.sendErrorToClient(clientID, "send_message_failed", err.Error())
		}

	case "read":
		// Mark message(s) as read
		if chatMsg.MessageID > 0 {
			req := ReadMessageRequest{
				MessageIDs: []int64{chatMsg.MessageID},
			}

			err := s.MarkMessagesAsRead(context.Background(), req, userID)
			if err != nil {
				log.Printf("Error marking message as read: %v", err)
				// Send error response to the client
				s.sendErrorToClient(clientID, "mark_read_failed", err.Error())
			}
		}

	case "typing":
		// Broadcast typing status to other participants
		s.broadcastTypingStatus(chatMsg.ConversationID, userID)

	default:
		log.Printf("Unknown chat action: %s", chatMsg.Action)
	}
}

// sendErrorToClient sends an error message to a specific client
func (s *ChatService) sendErrorToClient(clientID string, errorType string, errorMessage string) {
	errorMsg := map[string]interface{}{
		"type":    "error",
		"code":    errorType,
		"message": errorMessage,
	}

	s.wsManager.SendToClient(clientID, errorMsg)
}

// broadcastMessageToParticipants sends a message to all participants of a conversation
func (s *ChatService) broadcastMessageToParticipants(ctx context.Context, conversationID int64, dbMessage db.Message, sender db.User) {
	// Get all participants
	participants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get participants for message broadcast: %v", err)
		return
	}

	// Format the message for WebSocket
	wsMessage := WebSocketChatMessage{
		Type:           "chat_message",
		Action:         "new",
		MessageID:      dbMessage.ID,
		ConversationID: dbMessage.ConversationID,
		SenderID:       dbMessage.SenderID,
		Content:        dbMessage.Content,
		MessageType:    ChatMessageType(dbMessage.MessageType),
		CreatedAt:      dbMessage.CreatedAt.Time,
		SenderUsername: sender.Username,
		SenderName:     sender.FullName,
	}

	if len(dbMessage.Metadata) > 0 {
		wsMessage.Metadata = json.RawMessage(dbMessage.Metadata)
	}

	// Broadcast to all participants
	message := map[string]interface{}{
		"type": "chat_message",
		"data": wsMessage,
	}

	for _, participant := range participants {
		// Don't send to the sender
		if participant.ID == sender.ID {
			continue
		}

		// Check if participant is online
		clientID := fmt.Sprintf("user_%s", participant.Username)
		if !s.wsManager.SendToClient(clientID, message) {
			// User is offline, store the message for later delivery
			log.Printf("User %s offline, storing message for later delivery", participant.Username)
			err := s.wsManager.MessageStore.StoreMessage(
				ctx,
				clientID,
				participant.Username,
				"chat_message",
				message,
			)
			if err != nil {
				log.Printf("Error storing offline chat message: %v", err)
			}
		}
	}
}

// // notifyMessagesRead notifies participants when messages are read
// func (s *ChatService) notifyMessagesRead(ctx context.Context, messageIDs []int64, userID int64) {
// 	// In a real implementation, you would:
// 	// 1. Get the conversation ID from the first message
// 	// 2. Get all participants for that conversation
// 	// 3. Send notifications to all participants except the reader

// 	// This is simplified for now
// 	user, err := s.store.GetUserByID(ctx, userID)
// 	if err != nil {
// 		log.Printf("Failed to get user for read notification: %v", err)
// 		return
// 	}

// 	// Send read status only if we have the message details
// 	if len(messageIDs) > 0 {
// 		// Get message details to get conversation ID
// 		// Note: This assumes all messages are from the same conversation
// 		firstMessageID := messageIDs[0]
// 		// This would need a new query to get message details
// 		// For now, we'll skip this part
// 	}
// }

// notifyParticipantsAdded notifies all conversation participants about new members
func (s *ChatService) notifyParticipantsAdded(ctx context.Context, conversationID int64, newParticipantIDs []int64) {
	// Get all current participants
	participants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get participants for member-added notification: %v", err)
		return
	}

	// Get new participant details
	newParticipants := make([]ChatUser, 0, len(newParticipantIDs))
	for _, id := range newParticipantIDs {
		user, err := s.store.GetUserByID(ctx, id)
		if err != nil {
			log.Printf("Failed to get user %d for member-added notification: %v", id, err)
			continue
		}

		newParticipants = append(newParticipants, ChatUser{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		})
	}

	// Format the notification
	notification := map[string]interface{}{
		"type":           "chat_event",
		"event":          "participants_added",
		"conversationId": conversationID,
		"participants":   newParticipants,
		"timestamp":      time.Now(),
	}

	// Send to all participants
	for _, participant := range participants {
		clientID := fmt.Sprintf("user_%s", participant.Username)
		if !s.wsManager.SendToClient(clientID, notification) {
			// Store offline notification
			err := s.wsManager.MessageStore.StoreMessage(
				ctx,
				clientID,
				participant.Username,
				"chat_event",
				notification,
			)
			if err != nil {
				log.Printf("Error storing offline member-added notification: %v", err)
			}
		}
	}
}

// notifyParticipantRemoved notifies conversation members when someone leaves
func (s *ChatService) notifyParticipantRemoved(ctx context.Context, conversationID int64, removedUserID int64) {
	// Get all current participants
	participants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get participants for member-removed notification: %v", err)
		return
	}

	// Get removed user details
	removedUser, err := s.store.GetUserByID(ctx, removedUserID)
	if err != nil {
		log.Printf("Failed to get removed user for notification: %v", err)
		return
	}

	// Format the notification
	notification := map[string]interface{}{
		"type":            "chat_event",
		"event":           "participant_removed",
		"conversationId":  conversationID,
		"participantId":   removedUserID,
		"participantName": removedUser.FullName,
		"username":        removedUser.Username,
		"timestamp":       time.Now(),
	}

	// Send to all current participants
	for _, participant := range participants {
		clientID := fmt.Sprintf("user_%s", participant.Username)
		if !s.wsManager.SendToClient(clientID, notification) {
			// Store offline notification
			err := s.wsManager.MessageStore.StoreMessage(
				ctx,
				clientID,
				participant.Username,
				"chat_event",
				notification,
			)
			if err != nil {
				log.Printf("Error storing offline member-removed notification: %v", err)
			}
		}
	}
}

// notifyParticipantsNewConversation notifies users when they're added to a new conversation
func (s *ChatService) notifyParticipantsNewConversation(conversation ConversationResponse) {
	ctx := context.Background()

	// Format the notification
	notification := map[string]interface{}{
		"type":         "chat_event",
		"event":        "new_conversation",
		"conversation": conversation,
		"timestamp":    time.Now(),
	}

	// Send to all participants except the creator (assumed to be the first admin)
	for _, participant := range conversation.Participants {
		// Skip admins (creators)
		if participant.IsAdmin {
			continue
		}

		clientID := fmt.Sprintf("user_%s", participant.Username)
		if !s.wsManager.SendToClient(clientID, notification) {
			// Store offline notification
			err := s.wsManager.MessageStore.StoreMessage(
				ctx,
				clientID,
				participant.Username,
				"chat_event",
				notification,
			)
			if err != nil {
				log.Printf("Error storing offline new-conversation notification: %v", err)
			}
		}
	}
}

// broadcastTypingStatus notifies participants when someone is typing
func (s *ChatService) broadcastTypingStatus(conversationID int64, userID int64) {
	ctx := context.Background()

	// Get user info
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Failed to get user for typing notification: %v", err)
		return
	}

	// Get all participants
	participants, err := s.store.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		log.Printf("Failed to get participants for typing notification: %v", err)
		return
	}

	// Format the typing notification
	notification := map[string]interface{}{
		"type":           "chat_event",
		"event":          "typing",
		"conversationId": conversationID,
		"userId":         userID,
		"username":       user.Username,
		"fullName":       user.FullName,
		"timestamp":      time.Now(),
	}

	// Send to all participants except the one typing
	for _, participant := range participants {
		if participant.ID == userID {
			continue
		}

		clientID := fmt.Sprintf("user_%s", participant.Username)
		// For typing notifications, we don't store them for offline users
		s.wsManager.SendToClient(clientID, notification)
	}
}
