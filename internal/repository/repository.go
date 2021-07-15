package repository

import (
	"github.com/VTerenya/employees/internal/database"
	"github.com/VTerenya/employees/internal/employee"
	"github.com/VTerenya/employees/internal/position"
	"github.com/VTerenya/employees/internal/todoDatabase"
)

type Todo interface {
	CreatePosition(p position.Position)
	CreateEmployee(e employee.Employee)
	GetPositions() []position.Position
	GetEmployees() []employee.Employee
	GetPosition(id string) (position.Position, error)
	GetEmployee(id string) (employee.Employee, error)
	DeletePosition(id string) error
	DeleteEmployee(id string) error
	UpdatePosition(p position.Position) error
	UpdateEmployee(e employee.Employee) error
}

type Repository struct {
	Todo
}

func NewRepository(data *database.Database) *Repository {
	return &Repository{
		todoDatabase.NewTodoDatabase(data),
	}
}
