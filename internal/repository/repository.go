package repository

import (
	"github.com/VTerenya/employees/internal"
)

type Repo interface {
	GetPositions() map[string]internal.Position
	GetEmployees() map[string]internal.Employee
	AddPosition(p internal.Position)
	AddEmployee(e internal.Employee)
}

type Repository struct {
	data *Database
}

func NewRepository(data *Database) *Repository {
	return &Repository{
		data,
	}
}

func (t Repository) GetPositions() map[string]internal.Position {
	return t.data.GetPosition()
}

func (t Repository) GetEmployees() map[string]internal.Employee {
	return t.data.GetEmployees()
}

func (t Repository) AddPosition(p *internal.Position) {
	t.data.GetPosition()[p.ID.String()] = *p
}

func (t Repository) AddEmployee(e *internal.Employee) {
	t.data.GetEmployees()[e.ID.String()] = *e
}

type Database struct {
	employees map[string]internal.Employee
	positions map[string]internal.Position
}

func NewDataBase() *Database {
	return &Database{
		employees: map[string]internal.Employee{},
		positions: map[string]internal.Position{},
	}
}

func (d Database) GetEmployees() map[string]internal.Employee {
	return d.employees
}

func (d Database) GetPosition() map[string]internal.Position {
	return d.positions
}
