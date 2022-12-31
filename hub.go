package main

import (
	"encoding/json"

	"github.com/nakesto/chat-api/models"
)

type message struct {
	SenderName  string
	Message     string
	ReceiveName string
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			clientId := client.ID
			for client := range h.clients {
				msg := []byte("some one join room (ID: " + clientId + ")")
				client.send <- msg
			}

			h.clients[client] = true

		case client := <-h.unregister:
			clientId := client.ID
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			for client := range h.clients {
				msg := []byte("some one leave room (ID:" + clientId + ")")
				client.send <- msg
			}
		case userMessage := <-h.broadcast:
			var data message
			err := json.Unmarshal(userMessage, &data)

			if err != nil {
				break
			}

			for client := range h.clients {
				//prevent self and not receiver receive the message
				if client.Username == string(data.SenderName) {
					continue
				} else if client.Username == string(data.ReceiveName) {
					chat := models.Chat{}

					chat.SenderName = data.SenderName
					chat.ReceiveName = data.ReceiveName
					chat.Message = data.Message

					_, err := chat.SaveChat()

					if err != nil {
						continue
					}

					select {
					case client.send <- []byte(data.Message):
					default:
						close(client.send)
						delete(h.clients, client)
					}
				} else {
					continue
				}
			}
		}
	}
}
