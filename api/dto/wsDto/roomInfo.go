package wsDto

import (
	"github.com/google/uuid"
)

type RoomInfo struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Board       [][]int  `json:"board"`
	MetaBoard   [][]int  `json:"metaBoard"`
	AllowedMove [][]bool `json:"allowedMove"`
	GameStatus  int      `json:"gameStatus"`

	Occupants int `json:"occupants"`
	Capacity  int `json:"capacity"`
}
