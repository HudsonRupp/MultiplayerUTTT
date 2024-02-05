package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	// 0 or 1 (for managing turns)
	Player int

	Room *Room

	Conn *websocket.Conn

	//outbound messages
	Send chan []byte
}

func NewClient(player int, room *Room, ws *websocket.Conn) *Client {
	return &Client{Player: player,
		Room: room,
		Conn: ws,
		Send: make(chan []byte),
	}
}

func (c *Client) Read() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		//var msg dto.MoveRequest
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			log.Println("Closing socket")
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		//msg.Player = c.Player
		log.Println(string(message))
		c.Room.Recv <- message
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					log.Println(err.Error())
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
			log.Println("Ping sent")
		}
	}
}

func (c *Client) Close() {
	close(c.Send)
}
