package errors

var (
	badRequest                = newError("bad request")
	notFound                  = newError("not found")
	positionIsExists          = newError("position is exists")
	employeeIsExists          = newError("employee is exists")
	statusInternalServerError = newError("internal server error")
)

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

func StatusInternalServerError() error {
	return statusInternalServerError
}
