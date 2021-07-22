package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
)

type Service interface {
	CreatePosition(p *internal.Position) error
	CreateEmployee(e *internal.Employee) error
	GetPositions(limit, offset int) ([]internal.Position, error)
	GetEmployees(limit, offset int) ([]internal.Employee, error)
	GetPosition(id string) (internal.Position, error)
	GetEmployee(id string) (internal.Employee, error)
	DeletePosition(id string) error
	DeleteEmployee(id string) error
	UpdatePosition(p *internal.Position) error
	UpdateEmployee(e *internal.Employee) error
}

type Serv struct {
	repo *repository.Repo
}

func NewServ(repository *repository.Repo) *Serv {
	return &Serv{
		repo: repository,
	}
}

func (t Serv) CreatePosition(p *internal.Position) error {
	m := t.repo.GetPositions()
	for _, value := range m {
		if value.Salary == p.Salary && value.Name == p.Name {
			return errors.New("position is exists")
		}
	}
	p.ID = uuid.New()
	t.repo.AddPosition(p)
	return nil
}

func (t Serv) CreateEmployee(e *internal.Employee) error {
	m := t.repo.GetEmployees()
	for _, value := range m {
		if value.LasName == e.LasName &&
			value.FirstName == e.FirstName {
			return errors.New("employee is exists")
		}
	}
	e.ID = uuid.New()
	t.repo.AddEmployee(e)
	return nil
}

func (t Serv) GetPositions(limit, offset int) ([]internal.Position, error) {
	m := t.repo.GetPositions()
	positions := make([]internal.Position, 0)
	for _, value := range m {
		positions = append(positions, value)
	}
	offset--
	if float64(len(positions))/float64(limit) < float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.New("incorrect data")
	}
	var answer []internal.Position
	if len(positions) == 0 {
		return answer, nil
	}
	for i := limit * offset; i < limit*offset+limit && i < len(positions); i++ {
		answer = append(answer, positions[i])
	}
	return answer, nil
}

func (t Serv) GetEmployees(limit, offset int) ([]internal.Employee, error) {
	m := t.repo.GetEmployees()
	employees := make([]internal.Employee, 0)
	for _, value := range m {
		employees = append(employees, value)
	}
	offset--
	if float64(len(employees))/float64(limit) < float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.New("incorrect data")
	}
	var answer []internal.Employee
	if len(employees) == 0 {
		return answer, nil
	}
	for i := limit * offset; i < limit*offset+limit && i < len(employees); i++ {
		answer = append(answer, employees[i])
	}
	return answer, nil
}

func (t Serv) GetPosition(id string) (internal.Position, error) {
	m := t.repo.GetPositions()
	uID, err := uuid.Parse(id)
	if err != nil {
		return internal.Position{}, err
	}
	for _, value := range m {
		if value.ID == uID {
			return value, nil
		}
	}
	return internal.Position{}, errors.New("not found")
}

func (t Serv) GetEmployee(id string) (internal.Employee, error) {
	m := t.repo.GetEmployees()
	uID, err := uuid.Parse(id)
	if err != nil {
		return internal.Employee{}, err
	}
	for _, value := range m {
		if value.ID == uID {
			return value, nil
		}
	}
	return internal.Employee{}, errors.New("not found")
}

func (t Serv) DeletePosition(id string) error {
	m := t.repo.GetPositions()
	uID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	for key, value := range m {
		if value.ID == uID {
			delete(m, key)
		}
	}
	return errors.New("not found")
}

func (t Serv) DeleteEmployee(id string) error {
	m := t.repo.GetEmployees()
	uID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	for key, value := range m {
		if value.ID == uID {
			delete(m, key)
		}
	}
	return errors.New("not found")
}

func (t Serv) UpdatePosition(p *internal.Position) error {
	m := t.repo.GetPositions()
	if p.ID.String() == uuid.Nil.String() {
		return errors.New("incorrect data")
	}
	if _, ok := m[p.ID.String()]; ok {
		delete(m, p.ID.String())
		t.repo.AddPosition(p)
	}
	return errors.New("not found")
}

func (t Serv) UpdateEmployee(e *internal.Employee) error {
	m := t.repo.GetEmployees()
	fmt.Println(e.ID.String())
	if e.ID.String() == uuid.Nil.String() {
		return errors.New("incorrect data")
	}
	if _, ok := m[e.ID.String()]; ok {
		delete(m, e.ID.String())
		t.repo.AddEmployee(e)
	}
	return errors.New("not found")
}
