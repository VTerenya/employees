package repository

import (
	"github.com/VTerenya/employees/internal"
)

type Repository interface {
	GetPositions() map[string]internal.Position
	GetEmployees() map[string]internal.Employee
	AddPosition(p internal.Position)
	AddEmployee(e internal.Employee)
}

type Repo struct {
	data *Database
}

func NewRepo(data *Database) *Repo {
	return &Repo{
		data,
	}
}

func (t Repo) GetPositions() map[string]internal.Position {
	return t.data.GetPosition()
}

func (t Repo) GetEmployees() map[string]internal.Employee {
	return t.data.GetEmployees()
}

func (t Repo) AddPosition(p *internal.Position) {
	t.data.GetPosition()[p.ID.String()] = *p
}

func (t Repo) AddEmployee(e *internal.Employee) {
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
