package internal

type Employee struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"firstName"`
	LasName    string    `json:"lasName"`
	PositionID *Position `json:"position"`
}
