package service

import (
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/errors"
	"github.com/google/uuid"
)

type Serv struct {
	repo Repository
}

func NewServ(repository Repository) *Serv {
	return &Serv{
		repo: repository,
	}
}

func (t Serv) CreatePosition(p *internal.Position) error {
	m := t.repo.GetPositions()
	for _, value := range m {
		if value.Salary == p.Salary && value.Name == p.Name {
			return errors.PositionIsExists()
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
			return errors.EmployeeIsExists()
		}
	}
	e.ID = uuid.New()
	t.repo.AddEmployee(e)
	return nil
}

func (t Serv) GetPositions(limit, offset int) ([]internal.Position, error) {
	m := t.repo.GetPositions()
	answer := make([]internal.Position, 0)
	if len(m) == 0 && offset == 1 && limit == 1 {
		return answer, nil
	}
	positions := make([]internal.Position, 0)
	for _, value := range m {
		positions = append(positions, value)
	}
	offset--
	if float64(len(positions))/float64(limit) <= float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.NotFound()
	}
	for i := limit * offset; i < limit*offset+limit && i < len(positions); i++ {
		answer = append(answer, positions[i])
	}
	return answer, nil
}

func (t Serv) GetEmployees(limit, offset int) ([]internal.Employee, error) {
	m := t.repo.GetEmployees()
	answer := make([]internal.Employee, 0)
	if len(m) == 0 && offset == 1 && limit == 1 {
		return answer, nil
	}
	employees := make([]internal.Employee, 0)
	for _, value := range m {
		employees = append(employees, value)
	}
	offset--
	if float64(len(employees))/float64(limit) <= float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.NotFound()
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
	return internal.Position{}, errors.NotFound()
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
	return internal.Employee{}, errors.NotFound()
}

func (t Serv) DeletePosition(id string) error {
	return t.repo.DeletePosition(id)
}

func (t Serv) DeleteEmployee(id string) error {
	return t.repo.DeleteEmployee(id)
}

func (t Serv) UpdatePosition(p *internal.Position) error {
	if p.ID.String() == uuid.Nil.String() {
		return errors.BadRequest()
	}
	return t.repo.UpdatePosition(p)
}

func (t Serv) UpdateEmployee(e *internal.Employee) error {
	if e.ID.String() == uuid.Nil.String() {
		return errors.BadRequest()
	}
	return t.repo.UpdateEmployee(e)
}
