package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/service"
	"net/http"
	"regexp"
)

var (
	positionsRe      = regexp.MustCompile(`^\/positions[\/]*$`)
	positionRe       = regexp.MustCompile(`^\/position\/(\d+)$`)
	createPositionRe = regexp.MustCompile(`^\/position$`)

	employeesRe      = regexp.MustCompile(`^\/employees[\/]*$`)
	employeeRe       = regexp.MustCompile(`^\/employee\/(\d+)$`)
	createEmployeeRe = regexp.MustCompile(`^\/employee$`)
)

type Handler struct {
	service service.ServiceHandler
}

func NewHandler(ser *service.Service) *Handler {
	return &Handler{service: ser}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Println(r.URL.Path)
	switch {
	case r.Method == http.MethodGet && positionsRe.MatchString(r.URL.Path):
		h.getPositions(w, r)
		return
	case r.Method == http.MethodGet && employeesRe.MatchString(r.URL.Path):
		h.getEmployees(w, r)
		return

	case r.Method == http.MethodGet && positionRe.MatchString(r.URL.Path):
		h.getPosition(w, r)
		return
	case r.Method == http.MethodGet && employeeRe.MatchString(r.URL.Path):
		h.getEmployee(w, r)
		return

	case r.Method == http.MethodPost && createPositionRe.MatchString(r.URL.Path):
		h.createPosition(w, r)
		return
	case r.Method == http.MethodPost && createEmployeeRe.MatchString(r.URL.Path):
		h.createEmployee(w, r)
		return

	case r.Method == http.MethodDelete && positionRe.MatchString(r.URL.Path):
		h.deletePosition(w, r)
		return
	case r.Method == http.MethodDelete && employeeRe.MatchString(r.URL.Path):
		h.deleteEmployee(w, r)
		return

	case r.Method == http.MethodPut && createPositionRe.MatchString(r.URL.Path):
		h.updatePosition(w, r)
		return
	case r.Method == http.MethodPut && createEmployeeRe.MatchString(r.URL.Path):
		h.updateEmployee(w, r)
		return
	default:
		fmt.Println("Error: Not Found!")
		http.NotFound(w, r)
		return
	}
}

func (h *Handler) getPositions(w http.ResponseWriter, r *http.Request) {
	positions:=h.service.GetPositions()
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployees(w http.ResponseWriter, r *http.Request) {
	employees := h.service.GetEmployees()
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPosition(w http.ResponseWriter, r *http.Request) {
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	p, err := h.service.GetPosition(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Position not found"))
		if err != nil {
			http.Error(w, "Error: bad request", http.StatusBadRequest)
		}
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployee(w http.ResponseWriter, r *http.Request) {
	matches := employeeRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	e, err := h.service.GetEmployee(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("employee not found"))
		if err != nil {
			http.Error(w, "error: bad request", http.StatusBadRequest)
		}
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createPosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	err:=h.service.CreatePosition(p)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	err:=h.service.CreateEmployee(e)
	if err != nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	err := h.service.UpdatePosition(p)
	if err != nil {
		http.Error(w, "error: no content", http.StatusNoContent)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "error with servere", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	err := h.service.UpdateEmployee(e)
	if err != nil {
		http.Error(w, "error: no content", http.StatusNoContent)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deletePosition(w http.ResponseWriter, r *http.Request) {
	deleteID := positionRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeletePosition(deleteID)
	if err != nil {
		http.Error(w, "error: no content", http.StatusNoContent)
	}
	jsonBytes, err := json.Marshal(deleteID)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	deleteID := employeeRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeleteEmployee(deleteID)
	if err != nil {
		http.Error(w, "error: no content", http.StatusNoContent)
	}
	jsonBytes, err := json.Marshal(deleteID)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: bad request", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
