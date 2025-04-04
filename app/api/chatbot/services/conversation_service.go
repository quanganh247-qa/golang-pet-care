package services

import (
	"fmt"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
)

// ConversationService manages conversation state and history
type ConversationService struct {
	cache *CacheService
}

// NewConversationService creates a new conversation service
func NewConversationService(cache *CacheService) *ConversationService {
	return &ConversationService{
		cache: cache,
	}
}

// GetConversation retrieves a conversation state by ID
func (s *ConversationService) GetConversation(conversationID string) (*models.ConversationState, error) {
	cacheKey := fmt.Sprintf("conversation:%s", conversationID)
	var state models.ConversationState

	if !s.cache.GetJSON(cacheKey, &state) {
		return nil, fmt.Errorf("conversation not found or expired")
	}

	return &state, nil
}

// SaveConversation saves a conversation state
func (s *ConversationService) SaveConversation(state *models.ConversationState) error {
	cacheKey := fmt.Sprintf("conversation:%s", state.ConversationID)

	// Trim conversation history if too long (keep last 20 messages)
	if len(state.Messages) > 20 {
		state.Messages = state.Messages[len(state.Messages)-20:]
	}

	// Cache conversation for 24 hours (or use default from cache service)
	s.cache.SetWithPriority(cacheKey, state, 86400, PriorityMedium)

	return nil
}

// ListActiveConversations returns a list of active conversations for a user
func (s *ConversationService) ListActiveConversations(userID string) ([]models.ConversationState, error) {
	// This is a simple implementation that would be better with a proper database
	// In a real system, you would query a database for active conversations

	// For now, we'll scan all conversation keys in the cache
	// This is inefficient and just for demonstration
	keys := s.cache.CacheKeys()

	var conversations []models.ConversationState
	for _, key := range keys {
		if len(key) > 13 && key[:13] == "conversation:" {
			var state models.ConversationState
			if s.cache.GetJSON(key, &state) {
				if state.UserID == userID {
					conversations = append(conversations, state)
				}
			}
		}
	}

	return conversations, nil
}

// GetRecentConversations returns recent conversations, limited by count
func (s *ConversationService) GetRecentConversations(userID string, limit int) ([]models.ConversationState, error) {
	conversations, err := s.ListActiveConversations(userID)
	if err != nil {
		return nil, err
	}

	// Sort by last activity (most recent first)
	// Simple bubble sort for demonstration
	for i := 0; i < len(conversations); i++ {
		for j := i + 1; j < len(conversations); j++ {
			if conversations[i].LastActivity < conversations[j].LastActivity {
				conversations[i], conversations[j] = conversations[j], conversations[i]
			}
		}
	}

	// Limit results
	if len(conversations) > limit {
		conversations = conversations[:limit]
	}

	return conversations, nil
}

// AddMessageToConversation adds a new message to a conversation
func (s *ConversationService) AddMessageToConversation(conversationID string, message string, isBot bool) error {
	state, err := s.GetConversation(conversationID)
	if err != nil {
		return err
	}

	state.Messages = append(state.Messages, models.MessageEntry{
		Message:   message,
		Timestamp: time.Now().Unix(),
		IsBot:     isBot,
	})

	state.LastActivity = time.Now().Unix()

	return s.SaveConversation(state)
}

// DeleteConversation removes a conversation from the cache
func (s *ConversationService) DeleteConversation(conversationID string) error {
	cacheKey := fmt.Sprintf("conversation:%s", conversationID)
	s.cache.Delete(cacheKey)
	return nil
}

// SearchConversations searches for conversations by keyword
func (s *ConversationService) SearchConversations(userID string, keyword string) ([]models.ConversationState, error) {
	conversations, err := s.ListActiveConversations(userID)
	if err != nil {
		return nil, err
	}

	var results []models.ConversationState
	for _, conv := range conversations {
		// Check if any message contains the keyword
		for _, msg := range conv.Messages {
			if contains(msg.Message, keyword) {
				results = append(results, conv)
				break
			}
		}
	}

	return results, nil
}

// contains is a helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	// This is a simple implementation
	// In a real system, you would use a more sophisticated search
	for i := 0; i <= len(s)-len(substr); i++ {
		if equal(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

// equal is a helper function to check if two strings are equal (case-insensitive)
func equal(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if lower(s1[i]) != lower(s2[i]) {
			return false
		}
	}
	return true
}

// lower is a helper function to convert a byte to lowercase
func lower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}
