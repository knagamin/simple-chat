package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client describes one of chat user
type client struct {
	socket   *websocket.Conn        // Websocket for the client
	send     chan *message          // channel for seding a message
	room     *room                  // the room which is joined by the client
	userData map[string]interface{} // user info
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("2006-01-02 03:04:05")
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
