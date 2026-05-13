package hub

import (
	"sync"

	"fileserver/internal/realtime/client"
)

type Hub struct {
	mu sync.RWMutex

	clients map[string]*client.Client

	register chan *client.Client

	unregister chan *client.Client

	broadcast chan []byte
}

func NewHub() *Hub {

	return &Hub{
		clients: make(map[string]*client.Client),

		register: make(chan *client.Client),

		unregister: make(chan *client.Client),

		broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {

	for {

		select {

		case client := <-h.register:

			h.mu.Lock()

			h.clients[client.ID] = client

			h.mu.Unlock()

		case client := <-h.unregister:

			h.mu.Lock()

			if _, exists := h.clients[client.ID]; exists {

				delete(h.clients, client.ID)

				close(client.Send)
			}

			h.mu.Unlock()

		case message := <-h.broadcast:

			h.mu.Lock()

			for _, client := range h.clients {

				select {

				case client.Send <- message:

				default:

					close(client.Send)

					delete(h.clients, client.ID)
				}
			}

			h.mu.Unlock()
		}
	}
}

func (h *Hub) Broadcast(
	message []byte,
) {

	h.broadcast <- message
}

func (h *Hub) Register(
	client *client.Client,
) {

	h.register <- client
}

func (h *Hub) Unregister(
	client *client.Client,
) {

	h.unregister <- client
}
