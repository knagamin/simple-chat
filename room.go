package main

import (
	"log"
	"net/http"
	"simple-chat/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message    // channel for fowarding message to other clients from a client
	join    chan *client     // channel for managing clients who is getting to join
	leave   chan *client     // channel for managing clients who is getting to leave
	clients map[*client]bool // have exsisting clients
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined.")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("A client left.")
		case msg := <-r.forward:
			for client := range r.clients {
				r.tracer.Trace("Recieved a message: ", msg.Message)
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- A message was sent to clients")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- Sending was failed. The client is being cleaned up.")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Faild to get the cookie:", err)
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
