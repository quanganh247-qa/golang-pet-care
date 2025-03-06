package websocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		}
	}
}

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
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current WebSocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

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
}
