package wsDto

type MoveRequest struct {
	X      int `json:"Xcord"`
	Y      int `json:"Ycord"`
	Player int
}
