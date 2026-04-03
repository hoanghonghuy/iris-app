package ws

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// writeWait thời gian tối đa cho phép ghi message vào connection
	writeWait = 10 * time.Second

	// pongWait thời gian tối đa chờ pong từ client
	pongWait = 60 * time.Second

	// pingPeriod gửi ping đến client theo chu kỳ (phải nhỏ hơn pongWait)
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize kích thước tối đa của tin nhắn từ client (bytes)
	maxMessageSize = 4096
)

// Client đại diện cho một WebSocket connection của một user cụ thể
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID uuid.UUID
	Send   chan []byte

	// OnMessage callback được gọi khi client gửi tin nhắn lên server
	OnMessage func(userID uuid.UUID, data []byte)
}

// ReadPump đọc messages từ WebSocket connection.
// Chạy trong goroutine riêng cho mỗi client.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WS] read error: %v", err)
			}
			break
		}
		if c.OnMessage != nil {
			c.OnMessage(c.UserID, message)
		}
	}
}

// WritePump ghi messages từ hub đến WebSocket connection.
// Chạy trong goroutine riêng cho mỗi client.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				// Hub đã đóng channel
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
