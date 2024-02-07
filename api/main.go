package main

import (
	"log"
	"main/dto/httpDto"
	"main/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {

	//hub := NewHub();
	//go hub.Run();

	r := gin.Default()
	hub := ws.NewHub()

	r.Static("/assets", "../www/assets")
	r.LoadHTMLGlob("../www/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/game", func(c *gin.Context) {
		c.HTML(http.StatusOK, "game.html", gin.H{})
	})

	r.GET("/ws/:roomId", func(c *gin.Context) {
		var roomId uuid.UUID
		var err error
		roomId, err = uuid.Parse(c.Param("roomId"))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid roomId"})
			return
		}
		if roomId.String() == "00000000-0000-0000-0000-000000000000" {
			roomId = uuid.New()
		}

		serveWS(c, roomId, hub)
	})
	r.Use(corsMiddleware())

	r.GET("/rooms", func(c *gin.Context) {
		resp := []httpDto.RoomInfoResponse{}
		for room := range hub.Rooms {
			resp = append(resp, httpDto.RoomInfoResponse{
				Name:      room.Name,
				Id:        room.Id,
				Occupants: len(room.Clients),
				Capacity:  2,
				Status:    room.Status,
			})
		}
		c.JSON(200, resp)
	})

	r.Run("localhost:8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWS(c *gin.Context, roomId uuid.UUID, hub *ws.Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var client *ws.Client

	room := getRoom(hub, roomId)

	client = ws.NewClient(1, room, conn)
	addClient(client, room)
}

func getRoom(hub *ws.Hub, roomId uuid.UUID) *ws.Room {

	log.Println(hub.Rooms)
	for room := range hub.Rooms {
		if room.Id.String() == roomId.String() {
			log.Println("Room already exists, adding to existing room")
			return room
		}
	}
	room := ws.NewRoom(hub, roomId)
	go room.Run()
	hub.AddRoom(room)
	return room
}

func addClient(client *ws.Client, room *ws.Room) {
	room.Register <- client

	go client.Write()
	go client.Read()
}
