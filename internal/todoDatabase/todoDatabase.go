package todoDatabase

import (
	"errors"
	"github.com/VTerenya/employees/internal/database"
	"github.com/VTerenya/employees/internal/employee"
	"github.com/VTerenya/employees/internal/position"
)

type TodoDatabase struct {
	Data *database.Database
}

func NewTodoDatabase(data *database.Database) *TodoDatabase {
	return &TodoDatabase{Data: data}
}

func (t TodoDatabase) CreatePosition(p position.Position) {
	t.Data.Positions[p.ID] = p
}

func (t TodoDatabase) CreateEmployee(e employee.Employee) {
	t.Data.Employees[e.ID] = e
}

func (t TodoDatabase) GetPositions() []position.Position {
	positions := make([]position.Position, 0, len(t.Data.Positions))
	for _, v := range t.Data.Positions {
		positions = append(positions, v)
	}
	return positions
}

func (t TodoDatabase) GetEmployees() []employee.Employee {
	employees := make([]employee.Employee, 0, len(t.Data.Employees))
	for _, v := range t.Data.Employees {
		employees = append(employees, v)
	}
	return employees
}

func (t TodoDatabase) GetPosition(id string) (position.Position, error) {
	for _, value := range t.Data.Positions {
		if value.ID == id {
			return value, nil
		}
	}
	return position.Position{}, errors.New("Error: not found")
}

func (t TodoDatabase) GetEmployee(id string) (employee.Employee, error) {
	for _, value := range t.Data.Employees {
		if value.ID == id {
			return value, nil
		}
	}
	return employee.Employee{}, errors.New("Error: not found")
}

func (t TodoDatabase) DeletePosition(id string) error {
	for key, pos := range t.Data.Positions {
		if pos.ID == id {
			delete(t.Data.Positions, key)
			return nil
		}
	}
	return errors.New("Error: not found")
}

func (t TodoDatabase) DeleteEmployee(id string) error {
	for key, pos := range t.Data.Employees {
		if pos.ID == id {
			delete(t.Data.Positions, key)
			return nil
		}
	}
	return errors.New("Error: not found")
}

func (t TodoDatabase) UpdatePosition(p position.Position) error {
	for key, value := range t.Data.Positions {
		if value.ID == p.ID {
			t.Data.Positions[key] = p
			return nil
		}
	}
	return errors.New("Error: not found")
}

func (t TodoDatabase) UpdateEmployee(e employee.Employee) error {
	for key, value := range t.Data.Positions {
		if value.ID == e.ID {
			t.Data.Employees[key] = e
			return nil
		}
	}
	return errors.New("Error: not found")
}
