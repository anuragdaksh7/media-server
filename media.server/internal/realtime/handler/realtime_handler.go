package handler

import (
	"encoding/json"
	"fileserver/internal/auth/utils"
	"fileserver/internal/realtime/client"
	"fileserver/internal/realtime/hub"
	"fileserver/internal/realtime/ws"
	"fileserver/internal/torrent/manager"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RealtimeHandler struct {
	Hub *hub.Hub

	TorrentManager *manager.TorrentManager
}

func NewRealtimeHandler(
	hub *hub.Hub,
	torrentManager *manager.TorrentManager,
) *RealtimeHandler {

	return &RealtimeHandler{
		Hub: hub,

		TorrentManager: torrentManager,
	}
}

type incomingMessage struct {
	Type string `json:"type"`
}

type frame struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

var pingFrame = frame{
	Type: "ping",
}

func (h *RealtimeHandler) WebSocket(
	c *gin.Context,
) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "missing token",
			},
		)
		return
	}

	if _, err := utils.ParseJWT(token); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid token",
			},
		)
		return
	}

	conn, err := ws.Upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	client := &client.Client{
		ID: uuid.NewString(),

		Connection: conn,

		Send: make(chan []byte, 32),
	}

	h.Hub.Register(client)

	h.sendSnapshot(client)

	go h.writePump(client)

	go h.readPump(client)
}

func (h *RealtimeHandler) sendSnapshot(
	client *client.Client,
) {
	bytes, err := json.Marshal(frame{
		Type: "torrent.snapshot",
		Data: h.TorrentManager.List(),
	})
	if err != nil {
		return
	}

	_ = client.Connection.WriteMessage(
		websocket.TextMessage,
		bytes,
	)
}

func (h *RealtimeHandler) writePump(
	client *client.Client,
) {

	defer func() {

		h.Hub.Unregister(client)

		client.Connection.Close()
	}()

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}

			if err := client.Connection.WriteMessage(
				websocket.TextMessage,
				message,
			); err != nil {
				return
			}
		case <-ticker.C:
			bytes, err := json.Marshal(pingFrame)
			if err != nil {
				continue
			}
			if err := client.Connection.WriteMessage(
				websocket.TextMessage,
				bytes,
			); err != nil {
				return
			}
		}
	}
}

func (h *RealtimeHandler) readPump(
	client *client.Client,
) {

	defer func() {

		h.Hub.Unregister(client)

		client.Connection.Close()
	}()

	for {

		_, payload, err := client.Connection.ReadMessage()

		if err != nil {
			return
		}

		var message incomingMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			continue
		}

		if message.Type == "ping" {
			bytes, err := json.Marshal(pingFrame)
			if err != nil {
				continue
			}

			select {
			case client.Send <- bytes:
			default:
			}
		}
	}
}

func (h *RealtimeHandler) BroadcastJSON(
	eventType string,
	payload interface{},
) {
	bytes, err := json.Marshal(frame{
		Type: eventType,
		Data: payload,
	})

	if err != nil {
		return
	}

	h.Hub.Broadcast(bytes)
}
