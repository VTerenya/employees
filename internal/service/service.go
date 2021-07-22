package service

import (
	"errors"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/google/uuid"
)

type ServiceHandler interface {
	CreatePosition(p *internal.Position) error
	CreateEmployee(e *internal.Employee) error
	GetPositions(limit, offset int) ([]internal.Position, error)
	GetEmployees(limit, offset int) ([]internal.Employee, error)
	GetPosition(id string) (internal.Position, error)
	GetEmployee(id string) (internal.Employee, error)
	DeletePosition(id string) error
	DeleteEmployee(id string) error
	UpdatePosition(p internal.Position) error
	UpdateEmployee(e internal.Employee) error
}

type Service struct {
	repo *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (t Service) CreatePosition(p *internal.Position) error {
	m := t.repo.GetPositions()
	for _, value := range m {
		if value.Salary == p.Salary && value.Name == p.Name {
			return errors.New("create error: this position is exists")
		}
	}
	p.ID = uuid.New()
	t.repo.AddPosition(p)
	return nil
}

func (t Service) CreateEmployee(e *internal.Employee) error {
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value.LasName == e.LasName &&
			value.FirstName == e.FirstName {
			return errors.New("create error: employee is exists")
		}
	}
	e.ID = uuid.New()
	t.repo.AddEmployee(e)
	return nil
}

func (t Service) GetPositions(limit, offset int) ([]internal.Position, error) {
	m := t.repo.GetPositions()
	positions := make([]internal.Position, 0)
	for _, value := range m {
		positions = append(positions, value)
	}
	offset-=1
	if float64(len(positions))/float64(limit) <float64(offset) || limit <1 || offset <0{
		return nil, errors.New("incorrect data")
	}
	var answer []internal.Position
	if len(positions) == 0{
		return answer,nil
	}

	for i := limit*offset; i < limit*offset+limit && i < len(positions); i++ {
		answer = append(answer, positions[i])
	}
	return answer, nil
}

func (t Service) GetEmployees(limit, offset int) ([]internal.Employee, error) {
	m := t.repo.GetEmployees()
	employees := make([]internal.Employee, 0)
	for _, value := range m {
		employees = append(employees, value)
	}
	offset-=1
	if float64(len(employees))/float64(limit) <float64(offset) || limit <1 || offset <0{
		return nil, errors.New("incorrect data")
	}
	var answer []internal.Employee
	if len(employees) == 0{
		return answer,nil
	}
	for i := limit*offset; i < limit*offset+limit && i < len(employees); i++ {
		answer = append(answer, employees[i])
	}
	return answer, nil
}

func (t Service) GetPosition(id string) (internal.Position, error) {
	m := t.repo.GetPositions()
	uId, err := uuid.Parse(id)
	if err != nil {
		return internal.Position{}, err
	}
	for _, value := range m {
		if value.ID == uId {
			return value, nil
		}
	}
	return internal.Position{}, errors.New("get error: no this position")
}

func (t Service) GetEmployee(id string) (internal.Employee, error) {
	m := t.repo.GetEmployees()
	uId, err := uuid.Parse(id)
	if err != nil {
		return internal.Employee{}, err
	}
	for _, value := range m {
		if value.ID == uId {
			return value, nil
		}
	}
	return internal.Employee{}, errors.New("get error: no this employee")
}

func (t Service) DeletePosition(id string) error {
	m := t.repo.GetPositions()
	uId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	for key, value := range m {
		if value.ID == uId {
			delete(m, key)
		}
	}
	return errors.New("delete error: no this position")
}

func (t Service) DeleteEmployee(id string) error {
	m := t.repo.GetEmployees()
	uId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	for key, value := range m {
		if value.ID == uId {
			delete(m, key)
		}
	}
	return errors.New("delete error: no this employee")
}

func (t Service) UpdatePosition(p internal.Position) error {
	m := t.repo.GetPositions()
	for _, value := range m {
		if value.ID == p.ID {
			value.Name = p.Name
			value.Salary = p.Salary
			return nil
		}
	}
	return errors.New("update error: no this employee")
}

func (t Service) UpdateEmployee(e internal.Employee) error {
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value.ID == e.ID {
			value.FirstName = e.FirstName
			value.LasName = e.LasName
			value.Position = e.Position
			return nil
		}
	}
	return errors.New("update error: no this employee")
}
