package httpDto

import (
	"github.com/google/uuid"
)

type RoomInfoResponse struct {
	Name      string    `json:"name"`
	Id        uuid.UUID `json:"id"`
	Occupants int       `json:"occupants"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
}
