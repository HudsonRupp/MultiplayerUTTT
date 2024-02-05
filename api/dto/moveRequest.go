package dto

type MoveRequest struct {
	X      int `json:"Xcord"`
	Y      int `json:"Ycord"`
	Player int
}
