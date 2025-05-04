package services

import (
	"fmt"
	"sort"
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

	if found := s.cache.GetJSON(cacheKey, &state); !found {
		return nil, fmt.Errorf("conversation not found")
	}

	return &state, nil
}

// SaveConversation saves a conversation state
func (s *ConversationService) SaveConversation(state *models.ConversationState) error {
	cacheKey := fmt.Sprintf("conversation:%s", state.ConversationID)

	// Update last activity timestamp
	state.LastActivity = time.Now().Unix()

	// Save to cache with expiration (e.g., 30 days)
	return s.cache.SetJSON(cacheKey, state, 30*24*time.Hour)
}

// DeleteConversation deletes a conversation from the cache
func (s *ConversationService) DeleteConversation(conversationID string) error {
	cacheKey := fmt.Sprintf("conversation:%s", conversationID)
	s.cache.Delete(cacheKey)
	return nil
}

// ListActiveConversations returns a list of active conversations for a user
func (s *ConversationService) ListActiveConversations(userID string) ([]models.ConversationState, error) {
	// This is a simple implementation that would be better with a proper database
	// In a real system, you would query a database for active conversations

	// For now, we'll scan all conversation keys in the cache
	// This is inefficient and just for demonstration
	keys := s.cache.CacheKeys()

	conversations := make([]models.ConversationState, 0)
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

	// Sort by last activity, most recent first
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].LastActivity > conversations[j].LastActivity
	})

	// Limit the number of results
	if len(conversations) > limit {
		conversations = conversations[:limit]
	}

	return conversations, nil
}
