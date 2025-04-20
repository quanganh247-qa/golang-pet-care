package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// MessageStatus represents the delivery status of a message
type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusDelivered MessageStatus = "delivered"
	StatusFailed    MessageStatus = "failed"
)

// OfflineMessage represents a message that could not be delivered
type OfflineMessage struct {
	ID          int64           `json:"id"`
	ClientID    string          `json:"client_id"`
	MessageType string          `json:"message_type"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
	DeliveredAt *time.Time      `json:"delivered_at,omitempty"`
	Status      MessageStatus   `json:"status"`
	RetryCount  int             `json:"retry_count"`
	Username    string          `json:"username"`
}

// MessageStore manages the storage and retrieval of offline messages
type MessageStore struct {
	storeDB       db.Store
	mutex         sync.RWMutex
	cleanupTicker *time.Ticker
	retryTicker   *time.Ticker
}

// NewMessageStore creates a new message store
func NewMessageStore(store db.Store) *MessageStore {
	ms := &MessageStore{
		storeDB:       store,
		cleanupTicker: time.NewTicker(24 * time.Hour),  // Clean up old messages daily
		retryTicker:   time.NewTicker(5 * time.Minute), // Retry delivery every 5 minutes
	}

	// Start cleanup routine
	go ms.cleanupRoutine()

	// Start retry routine
	go ms.retryRoutine()

	return ms
}

// StoreMessage stores a message for a client who is offline
func (ms *MessageStore) StoreMessage(ctx context.Context, clientID, username, messageType string, data interface{}) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %w", err)
	}

	// Check if this is a doctor client ID and handle the username correctly
	if strings.HasPrefix(clientID, "doctor_") {
		// Extract doctor ID from client ID (e.g., "doctor_1" -> "1")
		doctorIDStr := strings.TrimPrefix(clientID, "doctor_")
		doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid doctor ID format: %w", err)
		}

		// Get the doctor record to find the user_id
		doctor, err := ms.storeDB.GetDoctor(ctx, doctorID)
		if err != nil {
			return fmt.Errorf("failed to get doctor: %w", err)
		}

		// Get the user record using doctor's user_id
		user, err := ms.storeDB.GetUserByID(ctx, doctor.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user for doctor ID %d: %w", doctorID, err)
		}

		// Use the actual username from the user record
		username = user.Username
	} else if username == "" {
		// If no username is provided, we can't store the message
		return fmt.Errorf("username is required for non-doctor client IDs")
	}

	// Store the message in the database
	_, err = ms.storeDB.CreateOfflineMessage(ctx, db.CreateOfflineMessageParams{
		ClientID:    clientID,
		Username:    username,
		MessageType: messageType,
		Data:        jsonData,
		Status:      string(StatusPending),
		RetryCount:  0,
	})

	if err != nil {
		return fmt.Errorf("failed to store offline message: %w", err)
	}

	log.Printf("Stored offline message for client %s of type %s", clientID, messageType)
	return nil
}

// GetPendingMessages retrieves pending messages for a client
func (ms *MessageStore) GetPendingMessages(ctx context.Context) ([]OfflineMessage, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	// Get pending messages from the database
	dbMessages, err := ms.storeDB.GetPendingMessagesForClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending messages: %w", err)
	}

	messages := make([]OfflineMessage, 0, len(dbMessages))
	for _, dbMsg := range dbMessages {
		msg := OfflineMessage{
			ID:          dbMsg.ID,
			ClientID:    dbMsg.ClientID,
			MessageType: dbMsg.MessageType,
			Data:        dbMsg.Data,
			CreatedAt:   dbMsg.CreatedAt.Time,
			Status:      MessageStatus(dbMsg.Status),
			RetryCount:  int(dbMsg.RetryCount),
			Username:    dbMsg.Username,
		}

		if dbMsg.DeliveredAt.Valid {
			deliveredAt := dbMsg.DeliveredAt.Time
			msg.DeliveredAt = &deliveredAt
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// MarkMessageDelivered marks a message as delivered
func (ms *MessageStore) MarkMessageDelivered(ctx context.Context, messageID int64) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	err := ms.storeDB.MarkMessageDelivered(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}

	return nil
}

// MarkMessageFailed marks a message as failed and increments retry count
func (ms *MessageStore) MarkMessageFailed(ctx context.Context, messageID int64) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	err := ms.storeDB.IncrementMessageRetryCount(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	return nil
}

// cleanupRoutine periodically cleans up old delivered messages
func (ms *MessageStore) cleanupRoutine() {
	for {
		<-ms.cleanupTicker.C

		// Delete messages older than 30 days that have been delivered
		ctx := context.Background()
		cutoff := time.Now().AddDate(0, 0, -30)

		err := ms.storeDB.DeleteOldMessages(ctx, db.DeleteOldMessagesParams{
			Status:    string(StatusDelivered),
			CreatedAt: pgtype.Timestamptz{Time: cutoff, Valid: true},
		})

		if err != nil {
			log.Printf("Error cleaning up old messages: %v", err)
		}
	}
}

// retryRoutine periodically retries sending failed messages
func (ms *MessageStore) retryRoutine() {
	for {
		<-ms.retryTicker.C

		ctx := context.Background()

		// Get messages that need retry (status pending or failed, retry count < max)
		messages, err := ms.storeDB.GetMessagesForRetry(ctx, db.GetMessagesForRetryParams{
			Column1:    []string{string(StatusPending), string(StatusFailed)},
			RetryCount: 5, // Maximum number of retries
		})

		if err != nil {
			log.Printf("Error retrieving messages for retry: %v", err)
			continue
		}

		for _, msg := range messages {
			// Try to find the client and send the message
			// This will be handled by the WSClientManager
			log.Printf("Message %d scheduled for retry (attempt %d)", msg.ID, msg.RetryCount)
		}
	}
}
