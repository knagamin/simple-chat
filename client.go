package main

import (
	"github.com/gorilla/websocket"
)

// client describes one of chat user
type client struct {
	socket *websocket.Conn // Websocket for the client
	send   chan []byte     // channel for seding a message
	room   *room           // the room which is joined by the client
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
