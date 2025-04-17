package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// WSClient represents a connected WebSocket client
type WSClient struct {
	conn          *websocket.Conn
	send          chan []byte
	clientID      string
	username      string
	clientType    string // Can be "user", "doctor", etc.
	isAuthSuccess bool
	lastActivity  time.Time
}

// WSClientManager manages WebSocket connections
type WSClientManager struct {
	clientsMap   map[string]*WSClient // Map with clientID as key
	mutex        sync.Mutex
	broadcast    chan []byte
	register     chan *WSClient
	unregister   chan *WSClient
	MessageStore *MessageStore // Changed from messageStore to MessageStore (public)
	storeDB      db.Store
}

// Update NewWSClientManager to accept store
func NewWSClientManager(store db.Store) *WSClientManager {
	manager := &WSClientManager{
		clientsMap: make(map[string]*WSClient),
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		storeDB:    store,
	}

	// Initialize the message store
	manager.MessageStore = NewMessageStore(store)

	return manager
}

func (manager *WSClientManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clientsMap[client.clientID] = client
			manager.mutex.Unlock()

			// Check for pending messages on client registration
			go manager.deliverPendingMessages(client)

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clientsMap[client.clientID]; ok {
				delete(manager.clientsMap, client.clientID)
				close(client.send)
			}
			manager.mutex.Unlock()

		case message := <-manager.broadcast:
			manager.mutex.Lock()
			for _, client := range manager.clientsMap {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clientsMap, client.clientID)
				}
			}
			manager.mutex.Unlock()
		}
	}
}

// WebSocketMessage represents a message to be sent over WebSocket
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	ID        string      `json:"id,omitempty"`
	PatientID string      `json:"patientId,omitempty"`
	Message   string      `json:"message,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development - should be restricted in production
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// BroadcastToAll sends a message to all connected clients
func (manager *WSClientManager) BroadcastToAll(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}
	manager.broadcast <- data
	log.Printf("Sent message to all clients: %s", string(data))
}

// BroadcastNotification sends a notification to all connected clients
func (manager *WSClientManager) BroadcastNotification(notification db.Notification) {
	// Format notification as a WebSocket message
	wsMessage := WebSocketMessage{
		Type: "notification",
		Data: notification,
	}

	// Store the notification in the database first
	err := manager.storeDB.ExecWithTransaction(context.Background(), func(q *db.Queries) error {
		_, err := q.CreatetNotification(context.Background(), db.CreatetNotificationParams{
			Username:    notification.Username,
			Title:       notification.Title,
			Content:     notification.Content,
			NotifyType:  notification.NotifyType,
			RelatedID:   notification.RelatedID,
			RelatedType: notification.RelatedType,
		})
		return err
	})

	if err != nil {
		log.Printf("Error storing notification in database: %v", err)
	}

	// Try to send directly to the user if online
	clientID := fmt.Sprintf("user_%s", notification.Username)
	if !manager.SendToClient(clientID, wsMessage) {
		// User is offline, store the message for later delivery
		ctx := context.Background()
		err := manager.MessageStore.StoreMessage(ctx, clientID, notification.Username, "notification", notification)
		if err != nil {
			log.Printf("Error storing offline notification: %v", err)
		}
	}
}

// HandleWebSocket handles WebSocket connections
func (manager *WSClientManager) HandleWebSocket(c *gin.Context) {
	// Get client ID from query parameters
	clientID := c.Request.URL.Query().Get("clientId")
	if clientID == "" {
		// Try from header (set by middleware)
		clientID = c.Request.Header.Get("X-Client-ID")
		if clientID == "" {
			clientID = time.Now().String() // Use timestamp as fallback ID
		}
	}

	// Get username if available
	var username string
	if userVal, exists := c.Get("username"); exists {
		if u, ok := userVal.(string); ok {
			username = u
		}
	}

	// Extract client type (user, doctor, etc.)
	clientType := "user" // Default
	if strings.HasPrefix(clientID, "doctor_") {
		clientType = "doctor"
	}

	// Create WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create and register client
	client := &WSClient{
		conn:          conn,
		send:          make(chan []byte, 256),
		clientID:      clientID,
		username:      username,
		clientType:    clientType,
		isAuthSuccess: username != "", // If username is set, consider auth successful
		lastActivity:  time.Now(),
	}

	// Register client before starting pumps
	manager.mutex.Lock()
	manager.clientsMap[clientID] = client
	manager.mutex.Unlock()

	// Start goroutines for reading and writing
	go client.writePump(manager)
	go client.readPump(manager)

	// Send a welcome message
	welcomeMsg := WebSocketMessage{
		Type:    "connected",
		Message: fmt.Sprintf("Successfully connected to WebSocket server with ID: %s", clientID),
	}
	client.sendJSON(welcomeMsg)

	// Mark the client as successfully authenticated if we have username
	if username != "" {
		log.Printf("WebSocket client %s authenticated as %s", clientID, username)
		client.isAuthSuccess = true

		// Deliver any pending messages
		go manager.deliverPendingMessages(client)
	}
}

// SendToClient sends a message to a specific client and returns true if successful
func (manager *WSClientManager) SendToClient(clientID string, message interface{}) bool {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return false
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	client, ok := manager.clientsMap[clientID]
	if !ok {
		log.Printf("Client not found: %s", clientID)
		return false
	}

	select {
	case client.send <- data:
		return true
	default:
		log.Printf("Client channel full, closing connection: %s", clientID)
		close(client.send)
		delete(manager.clientsMap, clientID)
		return false
	}
}

// deliverPendingMessages delivers any pending messages to a client that just connected
func (manager *WSClientManager) deliverPendingMessages(client *WSClient) {
	// Wait a short time for connection to stabilize
	time.Sleep(500 * time.Millisecond)

	// Only deliver if the client is authenticated
	if !client.isAuthSuccess || client.username == "" {
		return
	}

	ctx := context.Background()
	messages, err := manager.MessageStore.GetPendingMessages(ctx)
	if err != nil {
		log.Printf("Error retrieving pending messages for client %s: %v", client.clientID, err)
		return
	}

	if len(messages) == 0 {
		return
	}

	log.Printf("Delivering %d pending messages to client %s", len(messages), client.clientID)

	for _, msg := range messages {
		// Parse the message data
		var msgData interface{}

		switch msg.MessageType {
		case "notification":
			var notification db.Notification
			if err := json.Unmarshal(msg.Data, &notification); err != nil {
				log.Printf("Error unmarshaling notification: %v", err)
				continue
			}
			msgData = WebSocketMessage{
				Type: "notification",
				Data: notification,
			}
		case "appointment_alert":
			// This is just an example, adjust based on your actual message types
			msgData = WebSocketMessage{
				Type: "appointment_alert",
				Data: json.RawMessage(msg.Data),
			}
		default:
			// Generic message
			msgData = WebSocketMessage{
				Type: msg.MessageType,
				Data: json.RawMessage(msg.Data),
			}
		}

		// Send the message
		if manager.SendToClient(client.clientID, msgData) {
			// Mark as delivered
			if err := manager.MessageStore.MarkMessageDelivered(ctx, msg.ID); err != nil {
				log.Printf("Error marking message %d as delivered: %v", msg.ID, err)
			}
		} else {
			// Failed to deliver
			if err := manager.MessageStore.MarkMessageFailed(ctx, msg.ID); err != nil {
				log.Printf("Error marking message %d as failed: %v", msg.ID, err)
			}
		}
	}
}

// readPump reads messages from the WebSocket connection
func (client *WSClient) readPump(manager *WSClientManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512 * 1024) // 512KB
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		client.lastActivity = time.Now()
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Update last activity time
		client.lastActivity = time.Now()

		// Process the received message
		log.Printf("Received message from client %s: %s", client.clientID, string(message))

		// Handle message - could be a client command or heartbeat
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err == nil {
			// Handle based on message type
			if wsMessage.Type == "heartbeat" {
				// Send pong back
				client.sendJSON(WebSocketMessage{
					Type: "pong",
				})
			}
		}
	}
}

// writePump writes messages to the WebSocket connection
func (client *WSClient) writePump(manager *WSClientManager) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The manager closed the channel
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sendJSON sends a JSON message to the client
func (client *WSClient) sendJSON(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling JSON message: %v", err)
		return
	}
	client.send <- data
}

// Add this method to check if a client is connected
func (manager *WSClientManager) IsClientConnected(clientID string) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	_, exists := manager.clientsMap[clientID]
	return exists
}
