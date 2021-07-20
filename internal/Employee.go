package internal

import "github.com/google/uuid"

type Employee struct {
	ID         uuid.UUID    `json:"id"`
	FirstName  string    `json:"firstName"`
	LasName    string    `json:"lasName"`
	PositionID *Position `json:"position"`
}
