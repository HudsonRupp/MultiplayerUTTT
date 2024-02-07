package wsDto

import (
	"github.com/google/uuid"
)

type RoomInfo struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Occupants int       `json:"occupants"`
	Capacity  int       `json:"capacity"`
	Board     [][]int   `json:"board"`
}
