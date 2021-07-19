package repository

import (
	"github.com/VTerenya/employees/internal"
	"strconv"
)

type RepositoryHandler interface {
	GetPositions() map[string]internal.Position
	GetEmployees() map[string]internal.Employee
	AddPosition(p internal.Position)
	AddEmployee(e internal.Employee)
}

type Repository struct {
	RepositoryHandler
}

func NewRepository(data *database) *Repository {
	return &Repository{
		NewTodoRepository(data),
	}
}

type todoRepository struct {
	data *database
}

func NewTodoRepository(data *database) *todoRepository {
	return &todoRepository{data: data}
}

func (t todoRepository) GetPositions() map[string]internal.Position {
	return t.data.GetPosition()

}

func (t todoRepository) GetEmployees() map[string]internal.Employee {
	return t.data.GetEmployees()
}

func (t todoRepository) AddPosition(p internal.Position) {
	t.data.GetPosition()[strconv.Itoa(len(t.data.GetPosition())+1)] = p
}

func (t todoRepository) AddEmployee(e internal.Employee) {
	t.data.GetEmployees()[strconv.Itoa(len(t.data.GetEmployees())+1)] = e
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
