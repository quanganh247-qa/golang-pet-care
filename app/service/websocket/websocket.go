package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	clientsMap map[string]*WSClient // Map with clientID as key
	mutex      sync.Mutex
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	storeDB    db.Store
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

	return manager
}

func (manager *WSClientManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clientsMap[client.clientID] = client
			manager.mutex.Unlock()

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clientsMap[client.clientID]; ok {
				delete(manager.clientsMap, client.clientID)
				close(client.send)
			}
			manager.mutex.Unlock()

		case message := <-manager.broadcast:
			manager.mutex.Lock()
			for clientID, client := range manager.clientsMap {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clientsMap, clientID)
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
	// Chat-specific fields
	ConversationID int64  `json:"conversationId,omitempty"`
	Action         string `json:"action,omitempty"`    // For chat: "send", "read", "typing"
	MessageID      int64  `json:"messageId,omitempty"` // For referencing specific messages
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development - should be restricted in production
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// MessageTypeHandler is a function that handles a specific type of WebSocket message
type MessageTypeHandler func(message WebSocketMessage, clientID string, userID int64)

// messageHandlers stores registered handlers for message types
var messageHandlers = make(map[string]MessageTypeHandler)

// RegisterMessageTypeHandler registers a handler for a specific message type
func (manager *WSClientManager) RegisterMessageTypeHandler(messageType string, handler MessageTypeHandler) {
	messageHandlers[messageType] = handler
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

// HandleWebSocket handles WebSocket connections
func (manager *WSClientManager) HandleWebSocket(c *gin.Context) {
	// Get client ID from header (set by middleware or controller)
	clientID := c.Request.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = fmt.Sprintf("temp_%d", time.Now().UnixNano()) // Use timestamp as fallback ID
	}

	// Get username if available (might be empty for unauthenticated connections)
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
	manager.register <- client

	// Start goroutines for reading and writing
	go client.writePump(manager)
	go client.readPump(manager)

	// Send a welcome message
	welcomeMsg := WebSocketMessage{
		Type:    "connected",
		Message: fmt.Sprintf("Successfully connected to WebSocket server with ID: %s", clientID),
	}
	client.sendJSON(welcomeMsg)

	// For unauthenticated connections, send auth required message
	if !client.isAuthSuccess {
		authRequiredMsg := WebSocketMessage{
			Type:    "auth_required",
			Message: "Please authenticate to continue",
		}
		client.sendJSON(authRequiredMsg)
	}
}

// SendToClient sends a message to a specific client and returns true if successful
func (manager *WSClientManager) SendToClient(clientID string, message interface{}) bool {
	data, err := json.Marshal(message)
	if err != nil {
		return false
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	client, ok := manager.clientsMap[clientID]
	if !ok {
		return false
	}

	select {
	case client.send <- data:
		return true
	default:
		close(client.send)
		delete(manager.clientsMap, clientID)
		return false
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
				continue
			}

			// Allow authenticate messages even for unauthenticated clients
			if wsMessage.Type == "authenticate" {
				// Check if we have a handler for this message type
				if handler, ok := messageHandlers["authenticate"]; ok {
					// Call the handler with the message - pass 0 for userID since not authenticated yet
					handler(wsMessage, client.clientID, 0)
				}
				continue
			}

			// For all other message types, require authentication
			if !client.isAuthSuccess && wsMessage.Type != "authenticate" {
				client.sendJSON(WebSocketMessage{
					Type:    "error",
					Message: "Authentication required",
				})
				continue
			}

			// Extract user ID from client ID if it exists
			var userID int64 = 0
			if strings.HasPrefix(client.clientID, "user_") {
				// Find user details from username
				if client.username != "" {
					// In a real implementation, you'd look up the user's ID
					// This is simplified for now
					// You could add a new method to store to get user ID by username
					// userID = getUserIDFromUsername(client.username)
				}
			} else if strings.HasPrefix(client.clientID, "doctor_") {
				doctorIDStr := strings.TrimPrefix(client.clientID, "doctor_")
				if doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64); err == nil {
					userID = doctorID
				}
			}

			// Check if we have a handler for this message type
			if handler, ok := messageHandlers[wsMessage.Type]; ok {
				// Call the handler with the message
				handler(wsMessage, client.clientID, userID)
			} else {
				log.Printf("No handler registered for message type: %s", wsMessage.Type)
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

// IsClientConnected checks if a client is connected
func (manager *WSClientManager) IsClientConnected(clientID string) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	_, exists := manager.clientsMap[clientID]
	return exists
}

// UpdateClientUser updates client information after authentication
// This allows a client to authenticate after connecting with a temporary ID
func (manager *WSClientManager) UpdateClientUser(oldClientID string, newClientID string, username string, userID int64) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	// Find the client with the old ID
	client, exists := manager.clientsMap[oldClientID]
	if !exists {
		return false
	}

	// Update client information
	client.clientID = newClientID
	client.username = username
	client.isAuthSuccess = true

	// Remove the old client entry and add with the new ID
	delete(manager.clientsMap, oldClientID)
	manager.clientsMap[newClientID] = client

	// Log the change
	log.Printf("Client ID updated: %s -> %s (user: %s, ID: %d)", oldClientID, newClientID, username, userID)

	return true
}
