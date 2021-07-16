package service

import (
	"errors"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
)

type Handler interface {
	CreatePosition(p internal.Position) error
	CreateEmployee(e internal.Employee) error
	GetPositions() []internal.Position
	GetEmployees() []internal.Employee
	GetPosition(id string) (internal.Position, error)
	GetEmployee(id string) (internal.Employee, error)
	DeletePosition(id string) error
	DeleteEmployee(id string) error
	UpdatePosition(p internal.Position) error
	UpdateEmployee(e internal.Employee) error
}

type Service struct {
	Handler
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		newTodoService(repository),
	}
}

type todoService struct {
	repo *repository.Repository
}

func (t todoService) CreatePosition(p internal.Position) error{
	m := t.repo.GetPositions()
	for _, value := range m {
		if value == p {
			return errors.New("create error: this position is exists")
		}
	}
	t.repo.AddPosition(p)
	return nil
}

func (t todoService) CreateEmployee(e internal.Employee) error{
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value == e {
			return errors.New("create error: no this employee")
		}
	}
	t.repo.AddEmployee(e)
	return nil
}

func (t todoService) GetPositions() []internal.Position {
	m:=t.repo.GetPositions()
	positions :=make([]internal.Position,0)
	for _, value := range m {
		positions = append(positions, value)
	}
	return positions
}

func (t todoService) GetEmployees() []internal.Employee {
	m:=t.repo.GetEmployees()
	employees :=make([]internal.Employee,0)
	for _, value := range m {
		employees = append(employees, value)
	}
	return employees
}

func (t todoService) GetPosition(id string) (internal.Position, error) {
	m:=t.repo.GetPositions()
	for _, value := range m{
		if value.ID == id{
			return value, nil
		}
	}
	return internal.Position{}, errors.New("get error: no this position")
}

func (t todoService) GetEmployee(id string) (internal.Employee, error) {
	m:=t.repo.GetEmployees()
	for _, value := range m{
		if value.ID == id{
			return value, nil
		}
	}
	return internal.Employee{}, errors.New("get error: no this employee")
}

func (t todoService) DeletePosition(id string) error {
	m:=t.repo.GetPositions()
	for key, value := range m{
		if value.ID == id{
			delete(m, key)
		}
	}
	return errors.New("delete error: no this position")
}

func (t todoService) DeleteEmployee(id string) error {
	m:=t.repo.GetEmployees()
	for key, value := range m{
		if value.ID == id{
			delete(m, key)
		}
	}
	return errors.New("delete error: no this employee")
}

func (t todoService) UpdatePosition(p internal.Position) error {
	m:=t.repo.GetPositions()
	for _, value := range m{
		if value.ID == p.ID{
			value.Name = p.Name
			value.Salary = p.Salary
			return nil
		}
	}
	return errors.New("update error: no this employee")
}

func (t todoService) UpdateEmployee(e internal.Employee) error {
	m:=t.repo.GetEmployees()
	for _, value := range m{
		if value.ID == e.ID{
			value.FirstName = e.FirstName
			value.LasName = e.LasName
			value.PositionID = e.PositionID
			return nil
		}
	}
	return errors.New("update error: no this employee")
}

func newTodoService(repository *repository.Repository) *todoService{
	return &todoService{
		repo: repository,
	}
}

