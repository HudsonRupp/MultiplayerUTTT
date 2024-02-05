package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {

	//hub := NewHub();
	//go hub.Run();

	r := gin.Default()
	hub := NewHub()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		roomId := uuid.MustParse("524e40f9-f5c1-48c2-8d65-0151c54d503c")
		serveWS(c, roomId, hub)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWS(c *gin.Context, roomId uuid.UUID, hub *Hub) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Upgraded connection")

	var client *Client

	log.Println("roomid Inp" + roomId.String())
	room := getRoom(hub, roomId)
	log.Println("roomid out" + room.Id.String())

	client = NewClient(1, room, ws)
	addClient(client, room)
}

func getRoom(hub *Hub, roomId uuid.UUID) *Room {

	log.Println(hub.Rooms)
	for room := range hub.Rooms {
		if room.Id.String() == roomId.String() {
			log.Println("Room already exists, adding to existing room")
			return room
		}
	}
	room := NewRoom(hub, roomId)
	go room.Run()
	hub.AddRoom(room)
	return room
}

func addClient(client *Client, room *Room) {
	room.Register <- client

	go client.Write()
	go client.Read()
}
