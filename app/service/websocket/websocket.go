package websocket

import (
<<<<<<< HEAD
	"encoding/json"
=======
>>>>>>> e859654 (Elastic search)
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
<<<<<<< HEAD
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
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	mutex      sync.Mutex
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

// NewWSClientManager creates a new WebSocket client manager
func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

// Run starts the WebSocket manager
func (manager *WSClientManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()
			log.Printf("Client connected: %s. Total clients: %d", client.clientID, len(manager.clients))

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
			}
			manager.mutex.Unlock()
			log.Printf("Client disconnected: %s. Total clients: %d", client.clientID, len(manager.clients))

		case message := <-manager.broadcast:
			manager.mutex.Lock()
			for client := range manager.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clients, client)
				}
			}
			manager.mutex.Unlock()
=======
)

// upgrader chuyển đổi HTTP connection thành WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Trong production nên kiểm tra origin
	},
}

// Hub quản lý tất cả các kết nối WebSocket
type Hub struct {
	// Map từ username đến danh sách các kết nối
	clients    map[string][]*Client
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	username string
	send     chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run chạy vòng lặp chính của Hub
func (h *Hub) Run() {
	for {
		select {
		// Xử lý các kết nối mới
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.username] = append(h.clients[client.username], client)
			h.mutex.Unlock()
		// Xử lý các kết nối bị đóng
		case client := <-h.unregister:
			h.mutex.Lock()
			if connections, ok := h.clients[client.username]; ok {
				for i, conn := range connections {
					if conn == client {
						h.clients[client.username] = append(connections[:i], connections[i+1:]...)
						break
					}
				}
				if len(h.clients[client.username]) == 0 {
					delete(h.clients, client.username)
				}
			}
			h.mutex.Unlock()
			close(client.send)
>>>>>>> e859654 (Elastic search)
		}
	}
}

<<<<<<< HEAD
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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	clientID := c.Query("clientId")
	if clientID == "" {
		clientID = time.Now().String() // Use timestamp as fallback ID
	}

	client := &WSClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		clientID: clientID,
	}

	manager.register <- client

	// Start goroutines for reading and writing
	go client.writePump(manager)
	go client.readPump(manager)

	// Send a welcome message
	welcomeMsg := WebSocketMessage{
		Type:    "connected",
		Message: "Successfully connected to WebSocket server",
	}
	client.sendJSON(welcomeMsg)
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
=======
// SendToUser gửi thông báo đến một user cụ thể
func (h *Hub) SendToUser(username string, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if clients, ok := h.clients[username]; ok {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				h.unregister <- client
			}
		}
	}
}

// HandleWebSocket xử lý kết nối WebSocket mới
func (h *Hub) HandleWebSocket(c *gin.Context) {
	username := c.GetString("username") // Lấy từ JWT token
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		username: username,
		send:     make(chan []byte, 256),
	}
	client.hub.register <- client

	// Chạy goroutines để xử lý đọc/ghi
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
>>>>>>> e859654 (Elastic search)
			if err != nil {
				return
			}
			w.Write(message)

<<<<<<< HEAD
			// Add queued messages to the current WebSocket message
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
=======
			// Add queued chat messages to the current WebSocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
>>>>>>> e859654 (Elastic search)
			}

			if err := w.Close(); err != nil {
				return
			}
<<<<<<< HEAD
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
=======
>>>>>>> e859654 (Elastic search)
		}
	}
}

<<<<<<< HEAD
// sendJSON sends a JSON message to the client
func (client *WSClient) sendJSON(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling JSON message: %v", err)
		return
	}
	client.send <- data
=======
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Process the message (e.g., broadcast to other clients, save to database, etc.)
		// For simplicity, we'll just log the message here.
		log.Printf("Received message from %s: %s", c.username, message)
	}
>>>>>>> e859654 (Elastic search)
}
