package client

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID string

	Connection *websocket.Conn

	Send chan []byte
}
