package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VTerenya/employees/internal/employee"
	"github.com/VTerenya/employees/internal/position"
	"log"
	"net/http"
	"regexp"
	"sync"
)

var (
	positionsRe      = regexp.MustCompile(`^\/positions[\/]*$`)
	positionRe       = regexp.MustCompile(`^\/position\/(\d+)$`)
	createPositionRe = regexp.MustCompile(`^\/position$`)

	employeesRe      = regexp.MustCompile(`^\/employees[\/]*$`)
	employeeRe       = regexp.MustCompile(`^\/employee\/(\d+)$`)
	createEmployeeRe = regexp.MustCompile(`^\/employee$`)
)

type Database struct {
	Positions map[string]position.Position
	Employees map[string]employee.Employee
	*sync.RWMutex
}

type Handler struct {
	Data *Database
}

func newDatabase() *Database {
	var p position.Position
	p.Salary = 500
	p.Name = "Worker"
	p.ID = "1"
	return &Database{
		Positions: map[string]position.Position{
			"1": p,
		},
		Employees: map[string]employee.Employee{
			"1": {"1", "Nova", "Tern", &p},
		},
		RWMutex: &sync.RWMutex{},
	}
}
func NewHandler() *Handler {
	return &Handler{Data: newDatabase()}
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
	w.Header().Set("Content-Type", "application/json")
	h.Data.RLock()
	positions := make([]position.Position, 0, len(h.Data.Positions))
	for _, v := range h.Data.Positions {
		positions = append(positions, v)
	}
	h.Data.RUnlock()
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.Data.RLock()
	employees := make([]employee.Employee, 0, len(h.Data.Employees))
	for _, v := range h.Data.Employees {
		employees = append(employees, v)
	}
	h.Data.RUnlock()
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) getPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	h.Data.RLock()
	p, ok := h.Data.Positions[matches[1]]
	h.Data.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("user not found"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) getEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	matches := employeeRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	h.Data.RLock()
	e, ok := h.Data.Employees[matches[1]]
	h.Data.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("user not found"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) createPosition(w http.ResponseWriter, r *http.Request) {
	var p position.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	h.Data.Lock()
	h.Data.Positions[p.ID] = p
	h.Data.Unlock()
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var e employee.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	h.Data.Lock()
	h.Data.Employees[e.ID] = e
	h.Data.Unlock()
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, er := w.Write(jsonBytes)
	if er != nil {
		log.Fatal(er)
	}
}

func (h *Handler) updatePosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p position.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	for key, value := range h.Data.Positions {
		if value.ID == p.ID {
			h.Data.Positions[key] = p
			jsonBytes, err := json.Marshal(p)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, er := w.Write(jsonBytes)
			if er != nil {
				log.Fatal(er)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e employee.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	for key, value := range h.Data.Positions {
		if value.ID == e.ID {
			h.Data.Employees[key] = e
			jsonBytes, err := json.Marshal(e)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, er := w.Write(jsonBytes)
			if er != nil {
				log.Fatal(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deletePosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deleteID := positionRe.FindStringSubmatch(r.URL.Path)[1]
	for key, pos := range h.Data.Positions {
		if pos.ID == deleteID {
			h.Data.Lock()
			delete(h.Data.Positions, key)
			h.Data.Unlock()
			jsonBytes, err := json.Marshal(deleteID)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, er := w.Write(jsonBytes)
			if er != nil {
				log.Fatal(er)
			}
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deleteID := employeeRe.FindStringSubmatch(r.URL.Path)[1]
	for key, value := range h.Data.Employees {
		if value.ID == deleteID {
			h.Data.Lock()
			delete(h.Data.Employees, key)
			h.Data.Unlock()
			jsonBytes, err := json.Marshal(deleteID)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, er := w.Write(jsonBytes)
			if er != nil {
				log.Fatal(er)
			}
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
