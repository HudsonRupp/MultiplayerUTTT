package ws

import (
	"encoding/json"
	"log"
	"main/dto/wsDto"

	"github.com/google/uuid"
)

type Room struct {
	Id      uuid.UUID
	Name    string
	Status  string
	Board   [][]int
	Clients map[*Client]bool
	//inbound from clients
	Recv       chan IncomingMessage
	Register   chan *Client
	Unregister chan *Client
	Hub        *Hub
}

func NewRoom(hub *Hub, roomId uuid.UUID) *Room {

	return &Room{
		Id:         roomId,
		Name:       "New Room",
		Status:     "Waiting",
		Board:      [][]int{},
		Clients:    make(map[*Client]bool),
		Recv:       make(chan IncomingMessage),
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
			updateClients(room)
		case client := <-room.Unregister:
			close(client.Send)
			delete(room.Clients, client)
			if len(room.Clients) == 0 {
				return
			}
			log.Println("Client " + client.Conn.RemoteAddr().String() + " unregistered")
			updateClients(room)
		case msg := <-room.Recv:
			log.Println("Message received")
			log.Println(msg)
			room.handleMessage(msg)
		}
	}
}

func (room *Room) handleMessage(msg IncomingMessage) error {

	switch msgType := msg.MessageType; msgType {
	case "editRoom":
		var req wsDto.EditRoomRequest
		err := json.Unmarshal([]byte(msg.Content), &req)
		if err != nil {
			return err
		}
		room.Name = req.Name
		updateClients(room)
	case "move":
		var req wsDto.MoveRequest
		err := json.Unmarshal([]byte(msg.Content), &req)
		if err != nil {
			return err
		}
		//msg.Sender.Send <- []byte("Move request received")

	}
	return nil
}

func updateClients(room *Room) {
	for client := range room.Clients {
		client.Send <- getRoomInfo(room)
	}
}

func getRoomInfo(room *Room) OutgoingMessage {
	info, _ := json.Marshal(wsDto.RoomInfo{
		Id:        room.Id,
		Name:      room.Name,
		Status:    room.Status,
		Occupants: len(room.Clients),
		Capacity:  2,
		Board:     room.Board,
	})
	return OutgoingMessage{
		MessageType: "RoomInfo",
		Content:     string(info),
	}
}
