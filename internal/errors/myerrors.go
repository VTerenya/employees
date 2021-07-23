package errors

var (
	badRequest       = newError("bad request")
	notFound         = newError("not found")
	positionIsExists = newError("position is exists")
	employeeIsExists = newError("employee is exists")
)

type myError struct {
	description string
}

func newError(desc string) *myError {
	return &myError{description: desc}
}

func (m myError) Error() string {
	return m.description
}

func BadRequest() error {
	return badRequest
}

func NotFound() error {
	return notFound
}

func PositionIsExists() error {
	return positionIsExists
}

func EmployeeIsExists() error {
	return employeeIsExists
}
