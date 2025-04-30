package hub

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Id string
	P  []byte
}

type Hub struct {
	clients      map[*websocket.Conn]bool
	broadcast    chan *Message
	broadcastall chan *Message
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
	mutex        sync.RWMutex
}

func New() *Hub {
	hub := new(Hub)

	hub.clients = make(map[*websocket.Conn]bool)

	hub.broadcast = make(chan *Message)

	hub.broadcastall = make(chan *Message)

	hub.register = make(chan *websocket.Conn)

	hub.unregister = make(chan *websocket.Conn)

	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mutex.Lock()

			h.clients[conn] = true

			h.mutex.Unlock()

		case conn := <-h.unregister:
			h.mutex.Lock()

			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)

				conn.Close()
			}

			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()

			for client := range h.clients {
				if message.Id != "" && message.Id == client.Params("id") {
					if err := client.WriteMessage(websocket.TextMessage, message.P); err != nil {
						delete(h.clients, client)

						client.Close()
					}
				}
			}

			h.mutex.RUnlock()

		case message := <-h.broadcastall:
			h.mutex.RLock()

			for client := range h.clients {
				if client.Params("id") == "" {
					if err := client.WriteMessage(websocket.TextMessage, message.P); err != nil {
						delete(h.clients, client)

						client.Close()
					}
				}
			}

			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.register <- conn
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.unregister <- conn
}

func (h *Hub) Broadcast(id string, p []byte) {
	h.broadcast <- &Message{Id: id, P: p}
}

func (h *Hub) BroadcastAll(p []byte) {
	h.broadcastall <- &Message{P: p}
}
