package handler

import (
	"encoding/json"
	errs "errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/errors"
	"github.com/shopspring/decimal"
)

var (
	positionsRe      = regexp.MustCompile(`[\\/]positions\?limit=(\d+)&offset=(\d+)$`)
	positionRe       = regexp.MustCompile(`[\\/]position[\\/](\S+)$`)
	createPositionRe = regexp.MustCompile(`[\\/]position$`)

	employeesRe      = regexp.MustCompile(`[\\/]employees\?limit=(\d+)&offset=(\d+)$`)
	employeeRe       = regexp.MustCompile(`[\\/]employee[\\/](\S+)$`)
	createEmployeeRe = regexp.MustCompile(`[\\/]employee$`)
)

type Hand struct {
	service Service
}

func NewHandler(service Service) *Hand {
	return &Hand{service: service}
}

func (h *Hand) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && positionsRe.MatchString(r.URL.String()):
		h.GetPositions(w, r)
		return
	case r.Method == http.MethodGet && employeesRe.MatchString(r.URL.String()):
		h.GetEmployees(w, r)
		return

	case r.Method == http.MethodGet && positionRe.MatchString(r.URL.String()):
		h.GetPosition(w, r)
		return
	case r.Method == http.MethodGet && employeeRe.MatchString(r.URL.String()):
		h.GetEmployee(w, r)
		return

	case r.Method == http.MethodPost && createPositionRe.MatchString(r.URL.Path):
		h.CreatePosition(w, r)
		return
	case r.Method == http.MethodPost && createEmployeeRe.MatchString(r.URL.Path):
		h.CreateEmployee(w, r)
		return

	case r.Method == http.MethodDelete && positionRe.MatchString(r.URL.String()):
		h.DeletePosition(w, r)
		return
	case r.Method == http.MethodDelete && employeeRe.MatchString(r.URL.String()):
		h.DeleteEmployee(w, r)
		return

	case r.Method == http.MethodPut && createPositionRe.MatchString(r.URL.String()):
		h.UpdatePosition(w, r)
		return
	case r.Method == http.MethodPut && createEmployeeRe.MatchString(r.URL.String()):
		h.UpdateEmployee(w, r)
		return
	default:
		fmt.Println("errors: not found")
		http.NotFound(w, r)
		return
	}
}

func (h *Hand) GetPositions(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limit, err := strconv.ParseInt(query["limit"][0], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offset, err := strconv.ParseInt(query["offset"][0], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	positions, err := h.service.GetPositions(int(limit), int(offset))
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) GetEmployees(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limit, err := strconv.ParseInt(query["limit"][0], 10, 64)
	fmt.Println(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offset, err := strconv.ParseInt(query["offset"][0], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	employees, err := h.service.GetEmployees(int(limit), int(offset))
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) GetPosition(w http.ResponseWriter, r *http.Request) {
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.service.GetPosition(matches[1])
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) GetEmployee(w http.ResponseWriter, r *http.Request) {
	matches := employeeRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	e, err := h.service.GetEmployee(matches[1])
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p.Salary == decimal.Zero || p.Name == "" {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	err := h.service.CreatePosition(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonBytes, err := json.Marshal(p.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Hand) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if e.LasName == "" || e.FirstName == "" || e.Position == "" {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	err := h.service.CreateEmployee(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonBytes, err := json.Marshal(e.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Hand) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdatePosition(&p)
	if err != nil {
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdateEmployee(&e)
	if err != nil {
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Hand) DeletePosition(w http.ResponseWriter, r *http.Request) {
	deleteID := positionRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeletePosition(deleteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(internal.Position{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Hand) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	deleteID := employeeRe.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeleteEmployee(deleteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(internal.Employee{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
