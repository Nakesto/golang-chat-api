package main

import (
	"encoding/json"
	"fmt"

	"github.com/nakesto/chat-api/models"
)

type message struct {
	SenderName  string
	Message     string
	ReceiveName string
}

type Response struct {
	Tipe string      `json:"type"`
	Data interface{} `json:"data"`
}

type Welcome struct {
	Pesan string `json:"message"`
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
				result := Response{Tipe: "welcome", Data: Welcome{Pesan: "Someone join socket dengan Id" + clientId}}
				sasa := []interface{}{result}
				end, err := json.Marshal(sasa)

				if err != nil {
					continue
				}

				msg := end
				client.send <- msg
			}
			h.clients[client] = true

		case client := <-h.unregister:
			clientId := client.ID
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			fmt.Println("some one leave room (ID:" + clientId + ")")
		case userMessage := <-h.broadcast:
			var data message
			err := json.Unmarshal(userMessage, &data)

			if err != nil {
				break
			}

			chat := models.Chat{}

			chat.SenderName = data.SenderName
			chat.ReceiveName = data.ReceiveName
			chat.Message = data.Message

			c, err := chat.SaveChat()

			if err != nil {
				continue
			}

			init := Response{Tipe: "message", Data: c}

			rooms := models.ChatRoom{SenderName: c.SenderName, ReceiveName: c.ReceiveName, LastMessage: c.Message, UpdatedAt: c.UpdatedAt, Sender: c.Sender, Receiver: c.Receiver}

			init2 := Response{Tipe: "room", Data: rooms}

			for client := range h.clients {
				//prevent self and not receiver receive the message
				if client.Username == string(data.SenderName) {
					resp := []interface{}{init, init2}
					result, err := json.Marshal(resp)

					if err != nil {
						continue
					}

					select {
					case client.send <- result:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				} else if client.Username == string(data.ReceiveName) {

					resp := []interface{}{init, init2}

					result, err := json.Marshal(resp)

					if err != nil {
						continue
					}

					select {
					case client.send <- result:
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
