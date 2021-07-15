package employee

import (
	"github.com/VTerenya/employees/internal/position"
)

type Employee struct {
	ID         string             `json:"id"`
	FirstName  string             `json:"firstName"`
	LasName    string             `json:"lasName"`
	PositionID *position.Position `json:"position"`
}
