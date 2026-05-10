package handler

import (
	"encoding/json"
	"fileserver/internal/realtime/client"
	"fileserver/internal/realtime/hub"
	"fileserver/internal/realtime/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RealtimeHandler struct {
	Hub *hub.Hub
}

func NewRealtimeHandler(
	hub *hub.Hub,
) *RealtimeHandler {

	return &RealtimeHandler{
		Hub: hub,
	}
}

func (h *RealtimeHandler) WebSocket(
	c *gin.Context,
) {

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

		Send: make(chan []byte),
	}

	h.Hub.Register(client)

	go h.writePump(client)

	go h.readPump(client)
}

func (h *RealtimeHandler) writePump(
	client *client.Client,
) {

	defer func() {

		h.Hub.Unregister(client)

		client.Connection.Close()
	}()

	for {

		message, ok := <-client.Send

		if !ok {
			return
		}

		err := client.Connection.WriteMessage(
			1,
			message,
		)

		if err != nil {
			return
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

		_, _, err := client.Connection.ReadMessage()

		if err != nil {
			return
		}
	}
}

func (h *RealtimeHandler) BroadcastJSON(
	eventType string,
	payload interface{},
) {

	data := map[string]interface{}{
		"type": eventType,
		"data": payload,
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		return
	}

	h.Hub.Broadcast(bytes)
}
