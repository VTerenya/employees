package internal

import "github.com/google/uuid"

type Position struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Salary string    `json:"salary"`
}
