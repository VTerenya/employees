package handler

import (
	"encoding/json"
	"fmt"
	employee2 "github.com/VTerenya/employees/internal/employee"
	"github.com/VTerenya/employees/internal/position"
	"github.com/VTerenya/employees/internal/repository"
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
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
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
	positions := h.repo.GetPositions()
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployees(w http.ResponseWriter, r *http.Request) {
	employees := h.repo.GetEmployees()
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPosition(w http.ResponseWriter, r *http.Request) {
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	p, err := h.repo.GetPosition(matches[1])
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
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployee(w http.ResponseWriter, r *http.Request) {
	matches := employeeRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	e, err := h.repo.GetEmployee(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("user not found"))
		if err != nil {
			http.Error(w, "Error: bad request", http.StatusBadRequest)
		}
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createPosition(w http.ResponseWriter, r *http.Request) {
	var p position.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	h.repo.CreatePosition(p)
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var e employee2.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	h.repo.CreateEmployee(e)
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updatePosition(w http.ResponseWriter, r *http.Request) {
	var p position.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	err := h.repo.UpdatePosition(p)
	if err != nil {
		http.Error(w, "Error: no content!", http.StatusNoContent)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	var e employee2.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	err := h.repo.UpdateEmployee(e)
	if err != nil {
		http.Error(w, "Error: no content!", http.StatusNoContent)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deletePosition(w http.ResponseWriter, r *http.Request) {
	deleteID := positionRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.repo.DeletePosition(deleteID)
	if err != nil {
		http.Error(w, "Error: no content", http.StatusNoContent)
	}
	jsonBytes, err := json.Marshal(deleteID)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	deleteID := employeeRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.repo.DeleteEmployee(deleteID)
	if err != nil {
		http.Error(w, "Error: no content", http.StatusNoContent)
	}
	jsonBytes, err := json.Marshal(deleteID)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "Error: bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
