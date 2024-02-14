package ws

import (
	"encoding/json"
	"log"
	"main/dto/wsDto"

	"github.com/google/uuid"
)

type Room struct {
	//room data
	Id   uuid.UUID
	Name string

	//game data
	Board           [][]int
	MetaBoard       [][]int
	AllowedNextMove [][]bool
	PlayerTurn      int
	GameStatus      int

	//ws data
	Clients    map[*Client]bool
	Recv       chan IncomingMessage
	Register   chan *Client
	Unregister chan *Client
	Hub        *Hub
}

func NewRoom(hub *Hub, roomId uuid.UUID) *Room {
	board := make([][]int, 9)
	for i := range board {
		board[i] = make([]int, 9)
	}
	metaBoard := make([][]int, 3)
	for i := range metaBoard {
		metaBoard[i] = make([]int, 3)
	}
	allowed := [][]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	return &Room{
		Id:              roomId,
		Name:            "New Room",
		GameStatus:      0,
		Board:           board,
		MetaBoard:       metaBoard,
		AllowedNextMove: allowed,
		Clients:         make(map[*Client]bool),
		Recv:            make(chan IncomingMessage),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Hub:             hub,
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
			if len(room.Clients) > 0 {
				client.Player = 2
			}
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
			sendError(msg.Sender, "Invalid moveRequest format")
			log.Println(err.Error())
			return err
		}
		room.makeMove(req.X, req.Y, msg.Sender)
		updateClients(room)
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
		Id:          room.Id,
		Name:        room.Name,
		Board:       room.Board,
		MetaBoard:   room.MetaBoard,
		AllowedMove: room.AllowedNextMove,
		GameStatus:  room.GameStatus,
		Occupants:   len(room.Clients),
		Capacity:    2,
	})
	return OutgoingMessage{
		MessageType: "RoomInfo",
		Content:     string(info),
	}
}

func sendError(client *Client, message string) {
	err, _ := json.Marshal(wsDto.ErrorResponse{Message: message})
	client.Send <- OutgoingMessage{
		MessageType: "error",
		Content:     string(err),
	}
}
func (room *Room) makeMove(x int, y int, client *Client) {
	if room.GameStatus != 0 {
		sendError(client, "Game is over")
		return
	}
	if !room.AllowedNextMove[x/3][y/3] {
		sendError(client, "Square is taken")
		return
	}
	if room.MetaBoard[x/3][y/3] != 0 {
		sendError(client, "Square game has been won")
	}

	newBoard := room.Board
	if newBoard[x][y] != 0 {
		sendError(client, "Space already taken")
		return
	}
	newBoard[x][y] = client.Player

	//make move, check if game over,
	metaBoard := getMetaBoard(newBoard)
	gameResult := check3by3(metaBoard)

	var allowedMove [][]bool

	if metaBoard[x%3][y%3] != 0 {
		allowedMove = [][]bool{
			{true, true, true},
			{true, true, true},
			{true, true, true},
		}
	} else {
		allowedMove = [][]bool{
			{false, false, false},
			{false, false, false},
			{false, false, false},
		}
		allowedMove[x%3][y%3] = true
	}
	for i := 0; i < len(metaBoard); i++ {
		for j := 0; j < len(metaBoard); j++ {
			if metaBoard[i][j] != 0 {
				allowedMove[i][j] = false
			}
		}
	}

	room.Board = newBoard
	room.AllowedNextMove = allowedMove
	room.MetaBoard = metaBoard
	if gameResult != 0 {
		room.GameStatus = gameResult
	}

	updateClients(room)
}
func getMetaBoard(board [][]int) [][]int {
	metaBoard := [][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	for i := 0; i < len(board)/3; i++ {
		for j := 0; j < len(board[i])/3; j++ {
			subBoard := [][]int{
				board[(i * 3)][(j * 3) : (j*3)+3],
				board[(i*3)+1][(j * 3) : (j*3)+3],
				board[(i*3)+2][(j * 3) : (j*3)+3],
			}
			metaBoard[i][j] = check3by3(subBoard)
		}
	}
	return metaBoard
}

func check3by3(board [][]int) int {
	//cool magic square trick
	weights := [][]int{
		{4, 9, 2},
		{3, 5, 7},
		{8, 1, 6},
	}

	var xScore int
	var oScore int

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == 1 {
				xScore += weights[i][j]
				if xScore == 15 {
					return 1
				}
			} else if board[i][j] == 2 {
				oScore += weights[i][j]
				if oScore == 15 {
					return 2
				}
			}
		}
	}
	return 0
}
