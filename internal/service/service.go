package service

import (
	"errors"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
)

type ServiceHandler interface {
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
	ServiceHandler
	repo *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (t Service) CreatePosition(p internal.Position) error{
	m := t.repo.GetPositions()
	for _, value := range m {
		if value == p {
			return errors.New("create error: this position is exists")
		}
	}
	t.repo.AddPosition(p)
	return nil
}

func (t Service) CreateEmployee(e internal.Employee) error{
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value == e {
			return errors.New("create error: employee is exists")
		}
	}
	t.repo.AddEmployee(e)
	return nil
}

func (t Service) GetPositions() []internal.Position {
	m:=t.repo.GetPositions()
	positions :=make([]internal.Position,0)
	for _, value := range m {
		positions = append(positions, value)
	}
	return positions
}

func (t Service) GetEmployees() []internal.Employee {
	m:=t.repo.GetEmployees()
	employees :=make([]internal.Employee,0)
	for _, value := range m {
		employees = append(employees, value)
	}
	return employees
}

func (t Service) GetPosition(id string) (internal.Position, error) {
	m:=t.repo.GetPositions()
	for _, value := range m{
		if value.(internal.PositionRequest).ID == id{
			return value, nil
		}
	}
	return internal.PositionResponse{}, errors.New("get error: no this position")
}

func (t Service) GetEmployee(id string) (internal.Employee, error) {
	m:=t.repo.GetEmployees()
	for _, value := range m{
		if value.(internal.EmployeeRequest).ID == id{
			return value, nil
		}
	}
	return internal.EmployeeResponse{}, errors.New("get error: no this employee")
}

func (t Service) DeletePosition(id string) error {
	m:=t.repo.GetPositions()
	for key, value := range m{
		if value.(internal.PositionRequest).ID == id{
			delete(m, key)
		}
	}
	return errors.New("delete error: no this position")
}

func (t Service) DeleteEmployee(id string) error {
	m:=t.repo.GetEmployees()
	for key, value := range m{
		if value.(internal.EmployeeRequest).ID == id{
			delete(m, key)
		}
	}
	return errors.New("delete error: no this employee")
}

func (t Service) UpdatePosition(p internal.Position) error {
	m:=t.repo.GetPositions()
	position := p.(internal.PositionRequest)
	for _, value := range m{
		tempPosition := value.(internal.PositionRequest)
		if tempPosition.ID == position.ID{
			tempPosition.Name = position.Name
			tempPosition.Salary = position.Salary
			return nil
		}
	}
	return errors.New("update error: no this employee")
}

func (t Service) UpdateEmployee(e internal.Employee) error {
	m:=t.repo.GetEmployees()
	employee := e.(internal.EmployeeRequest)
	for _, value := range m{
		tempEmployee := value.(internal.EmployeeRequest)
		if tempEmployee.ID == employee.ID{
			tempEmployee.FirstName = employee.FirstName
			tempEmployee.LasName = employee.LasName
			tempEmployee.PositionID = employee.PositionID
			return nil
		}
	}
	return errors.New("update error: no this employee")
}

