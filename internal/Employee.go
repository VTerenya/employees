package internal

type Employee interface{}

type EmployeeResponse struct {
	Employee
	FirstName  string    `json:"firstName"`
	LasName    string    `json:"lasName"`
	PositionID *Position `json:"position"`
}

type EmployeeRequest struct {
	Employee
	ID         string    `json:"id"`
	FirstName  string    `json:"firstName"`
	LasName    string    `json:"lasName"`
	PositionID *Position `json:"position"`
}
