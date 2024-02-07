package ws

import (
	"fmt"
	"log"
)

// Literally just keeps the addresses of the rooms
type Hub struct {
	Rooms map[*Room]bool
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[*Room]bool),
	}
}

func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room] = true
	log.Println("Room " + room.Id.String() + " added to hub")
	log.Println("There are now " + fmt.Sprintf("%v", len(h.Rooms)) + " rooms")
}

func (h *Hub) DestroyRoom(room *Room) {
	delete(h.Rooms, room)
	close(room.Register)
	close(room.Unregister)
	close(room.Recv)
	log.Println("Room " + fmt.Sprintf("%p", room) + " destroyed")
}

/*
func (h *Hub) Run() {
	for {
		select {
		case room := <-h.Register:
			h.Rooms[room] = true
			log.Println("Room " + room.Id.String() + " registered")
		case room := <-h.Unregister:
			if _, ok := h.Rooms[room]; ok {
				delete(h.Rooms, room)
				close(room.Register)
				close(room.Unregister)
				close(room.Recv)
			}
			log.Println("Room " + room.Id.String() + " unregistered")
		}

	}
}
*/
