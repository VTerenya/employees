package errors

type Errors struct {
	description string
}

func newError(desc string) *Errors {
	return &Errors{description: desc}
}

func (m Errors) Error() string {
	return m.description
}

func BadRequest() error {
	return newError("bad request")
}

func NotFound() error {
	return newError("not found")
}

func PositionIsExists() error {
	return newError("position is exists")
}

func EmployeeIsExists() error {
	return newError("employee is exists")
}

func StatusInternalServerError() error {
	return newError("internal server error")
}

func PositionIsNotExists() error {
	return newError("position is not exists")
}
