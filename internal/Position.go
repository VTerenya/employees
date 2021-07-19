package internal

type Position interface{}

type PositionRequest struct{
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary float32 `json:"salary"`
}

type PositionResponse struct {
	Name   string  `json:"name"`
	Salary float32 `json:"salary"`
}
