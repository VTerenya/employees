package service

import (
	"errors"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/google/uuid"
	"strconv"
)

type ServiceHandler interface {
	CreatePosition(p internal.Position) error
	CreateEmployee(e internal.Employee) error
	GetPositions(limit, offset string) ([]internal.Position,error)
	GetEmployees(limit, offset string) ([]internal.Employee, error)
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

func (t Service) CreatePosition(p internal.Position) error {
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

func (t Service) CreateEmployee(e internal.Employee) error {
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value.LasName == e.FirstName &&
			value.FirstName == e.FirstName &&
			value.PositionID == e.PositionID{
			return errors.New("create error: employee is exists")
		}
	}
	e.ID = uuid.New()
	t.repo.AddEmployee(e)
	return nil
}

func (t Service) GetPositions(limit, offset string) ([]internal.Position,error) {
	m := t.repo.GetPositions()
	positions := make([]internal.Position, 0)
	for _, value := range m {
		positions = append(positions, value)
	}
	intLimit, err := strconv.ParseInt(limit,10, 32)
	if err!=nil{
		return nil, err
	}
	intOffset, err := strconv.ParseInt(offset,10,32)
	if err != nil{
		return nil, err
	}
	if len(positions)/int(intLimit) >= int(intOffset){
		return nil, errors.New("incorrect data")
	}
	answer := make([]internal.Position, 0)
	for i := len(positions)/int(intLimit)*5; i < len(positions)/int(intLimit)*5 + int(intLimit); i++ {
		answer = append(answer, positions[i])
	}
	return answer, nil
}

func (t Service) GetEmployees(limit, offset string) ([]internal.Employee,error) {
	m := t.repo.GetEmployees()
	employees := make([]internal.Employee, 0)
	for _, value := range m {
		employees = append(employees, value)
	}
	intLimit, err := strconv.ParseInt(limit,10, 32)
	if err!=nil{
		return nil, err
	}
	intOffset, err := strconv.ParseInt(offset,10,32)
	if err != nil{
		return nil, err
	}
	if len(employees)/int(intLimit) >= int(intOffset){
		return nil, errors.New("incorrect data")
	}
	answer := make([]internal.Employee, 0)
	for i := len(employees)/int(intLimit)*5; i < len(employees)/int(intLimit)*5 + int(intLimit); i++ {
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
			value.PositionID = e.PositionID
			return nil
		}
	}
	return errors.New("update error: no this employee")
}
