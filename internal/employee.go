package internal

import "github.com/google/uuid"

type Employee struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LasName    string    `json:"las_name"`
	PositionID uuid.UUID `json:"positionID"`
}
