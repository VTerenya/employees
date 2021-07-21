package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/service"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var (
	positionsRe      = regexp.MustCompile(`^\/positions\?limit=(\d+)&offset=(\d+)$`)
	positionRe       = regexp.MustCompile(`^\/position\/(\d+)$`)
	createPositionRe = regexp.MustCompile(`^\/position$`)

	employeesRe      = regexp.MustCompile(`^\/employees\?limit=(\d+)&offset=(\d+)$`)
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
		fmt.Println("error: not found")
		http.NotFound(w, r)
		return
	}
}

func (h *Handler) getPositions(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	limit, err := strconv.ParseInt(query["limit"][0], 10, 64)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	offset, err := strconv.ParseInt(query["offset"][0], 10, 64)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	positions, err := h.service.GetPositions(int(limit), int(offset))
	if err != nil {
		if err.Error() == "incorrect data" {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "error with server", http.StatusInternalServerError)
			return
		}
	}
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployees(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	limit, err := strconv.ParseInt(query["limit"][0], 10, 64)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	offset, err := strconv.ParseInt(query["offset"][0], 10, 64)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	employees, err := h.service.GetEmployees(int(limit), int(offset))
	if err != nil {
		if err.Error() == "incorrect data" {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "error with server", http.StatusInternalServerError)
			return
		}
	}
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getPosition(w http.ResponseWriter, r *http.Request) {
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "error with response", http.StatusBadRequest)
		return
	}
	p, err := h.service.GetPosition(matches[1])
	if err != nil {
		if err.Error() == "get error: no this position" {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployee(w http.ResponseWriter, r *http.Request) {
	matches := employeeRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "error with response", http.StatusBadRequest)
		return
	}
	e, err := h.service.GetEmployee(matches[1])
	if err != nil {
		if err.Error() == "get error: no this position" {
			http.Error(w, "error: not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createPosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	if p.Salary == "" || p.Name == "" {
		http.Error(w, "error with request", http.StatusBadRequest)
	}
	err := h.service.CreatePosition(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	if e.LasName == "" || e.FirstName == "" {
		http.Error(w, "error with request", http.StatusBadRequest)
	}
	err := h.service.CreateEmployee(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error: with server", http.StatusInternalServerError)
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
		http.Error(w, "error with server", http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, "error:with server", http.StatusInternalServerError)
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
		http.Error(w, "error: with server", http.StatusInternalServerError)
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
		http.Error(w, "error: with server", http.StatusInternalServerError)
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
		http.Error(w, "error: with server", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
