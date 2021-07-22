package repository

import (
	"github.com/VTerenya/employees/internal"
)

type RepositoryHandler interface {
	GetPositions() map[string]internal.Position
	GetEmployees() map[string]internal.Employee
	AddPosition(p internal.Position)
	AddEmployee(e internal.Employee)
}

type Repository struct {
	data *database
}

func NewRepository(data *database) *Repository {
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

type database struct {
	employees map[string]internal.Employee
	positions map[string]internal.Position
}

func NewDataBase() *database {
	return &database{
		employees: map[string]internal.Employee{},
		positions: map[string]internal.Position{},
	}
}

func (d database) GetEmployees() map[string]internal.Employee {
	return d.employees
}

func (d database) GetPosition() map[string]internal.Position {
	return d.positions
}
