package database

import (
	"github.com/VTerenya/employees/internal/employee"
	"github.com/VTerenya/employees/internal/position"
)

type Database struct {
	Positions map[string]position.Position
	Employees map[string]employee.Employee
}

func NewDatabase() *Database {
	return &Database{
		Positions: map[string]position.Position{},
		Employees: map[string]employee.Employee{},
	}
}
