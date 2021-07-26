package handler

import "github.com/VTerenya/employees/internal"

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
