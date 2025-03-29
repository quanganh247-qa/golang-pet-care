package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

// WSClient represents a connected WebSocket client
type WSClient struct {
	conn     *websocket.Conn
	send     chan []byte
	clientID string
}

// WSClientManager manages WebSocket connections
type WSClientManager struct {
	clientsMap map[string]*WSClient // Map with clientID as key
	mutex      sync.Mutex
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
}

// Update NewWSClientManager
func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		clientsMap: make(map[string]*WSClient),
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (manager *WSClientManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clientsMap[client.clientID] = client
			manager.mutex.Unlock()
			log.Printf("Client connected: %s. Total clients: %d", client.clientID, len(manager.clientsMap))

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clientsMap[client.clientID]; ok {
				delete(manager.clientsMap, client.clientID)
				close(client.send)
			}
			manager.mutex.Unlock()
			log.Printf("Client disconnected: %s. Total clients: %d", client.clientID, len(manager.clientsMap))

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
}

// BroadcastNotification sends a notification to all connected clients
func (manager *WSClientManager) BroadcastNotification(notification db.Notification) {
	// Format notification as a WebSocket message
	wsMessage := WebSocketMessage{
		Type: "notification",
		Data: notification,
	}

	manager.BroadcastToAll(wsMessage)
}

// HandleWebSocket handles WebSocket connections
func (manager *WSClientManager) HandleWebSocket(c *gin.Context) {
	// Get client ID from query parameters
	clientID := c.Request.URL.Query().Get("clientId")
	if clientID == "" {
		clientID = time.Now().String() // Use timestamp as fallback ID
	}

	// Create WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create and register client
	client := &WSClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		clientID: clientID,
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
}

// Update SendToClient method to use clientsMap
func (manager *WSClientManager) SendToClient(clientID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if client, ok := manager.clientsMap[clientID]; ok {
		select {
		case client.send <- data:
			log.Printf("Sent message to client: %s", clientID)
		default:
			close(client.send)
			delete(manager.clientsMap, clientID)
			log.Printf("Failed to send message to client: %s", clientID)
		}
	} else {
		log.Printf("Client not found: %s", clientID)
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

		// Process the received message
		log.Printf("Received message from client %s: %s", client.clientID, string(message))

		// Here you can process client->server messages if needed
		// For now, we'll just echo it back
		client.send <- message
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
