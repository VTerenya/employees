package internal

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Position struct {
	ID     uuid.UUID       `json:"ID"`
	Name   string          `json:"name"`
	Salary decimal.Decimal `json:"salary"`
}
