package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Hub quản lý tất cả WebSocket clients đang kết nối.
// Chịu trách nhiệm đăng ký, hủy đăng ký, và broadcast tin nhắn đến đúng user.
type Hub struct {
	// mu bảo vệ clients map
	mu sync.RWMutex

	// clients lưu trữ danh sách clients theo userID
	// Một user có thể có nhiều connections (nhiều tab)
	clients map[uuid.UUID]map[*Client]struct{}

	// register channel nhận client mới cần đăng ký
	Register chan *Client

	// unregister channel nhận client cần hủy đăng ký
	Unregister chan *Client
}

// NewHub tạo mới Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run khởi chạy vòng lặp xử lý register/unregister.
// Phải gọi trong goroutine riêng: go hub.Run()
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]struct{})
			}
			h.clients[client.UserID][client] = struct{}{}
			h.mu.Unlock()
			log.Printf("[WS] Client connected: userID=%s", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.UserID]; ok {
				delete(conns, client)
				if len(conns) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			log.Printf("[WS] Client disconnected: userID=%s", client.UserID)
		}
	}
}

// BroadcastEvent gửi một event JSON đến danh sách userIDs.
// Event được gửi đến tất cả connections (tabs) của mỗi user.
type WSEvent struct {
	Type string      `json:"type"` // "new_message" | "conversation_created" | ...
	Data interface{} `json:"data"`
}

// BroadcastToUsers gửi event đến danh sách userIDs
func (h *Hub) BroadcastToUsers(userIDs []uuid.UUID, event WSEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[WS] failed to marshal event: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, uid := range userIDs {
		if conns, ok := h.clients[uid]; ok {
			for client := range conns {
				select {
				case client.Send <- data:
				default:
					// buffer đầy → đóng connection
					close(client.Send)
					delete(conns, client)
				}
			}
		}
	}
}
