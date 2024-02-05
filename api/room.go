package main

import (
	"log"

	"github.com/google/uuid"
)

type Room struct {
	Id uuid.UUID

	Clients map[*Client]bool

	PlayerOne *Client
	PlayerTwo *Client

	//inbound from clients
	Recv chan []byte

	Register chan *Client

	Unregister chan *Client

	Board map[int]int

	Hub *Hub
}

func NewRoom(hub *Hub, roomId uuid.UUID) *Room {
	return &Room{
		Clients:    make(map[*Client]bool),
		Id:         roomId,
		Recv:       make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Hub:        hub,
	}
}

func (room *Room) Close() {
	close(room.Register)
	close(room.Unregister)
	close(room.Recv)
}

func (room *Room) Run() {
	defer func() {
		room.Hub.DestroyRoom(room)
	}()
	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
			log.Println("Client " + client.Conn.RemoteAddr().String() + " registered")
		case client := <-room.Unregister:
			close(client.Send)
			delete(room.Clients, client)
			if len(room.Clients) == 0 {
				return
			}
			log.Println("Client " + client.Conn.RemoteAddr().String() + " unregistered")
		case msgBytes := <-room.Recv:
			log.Println("bytes Received, echoing data: " + string(msgBytes))
			for client := range room.Clients {
				client.Send <- msgBytes
			}
		}
	}
}
