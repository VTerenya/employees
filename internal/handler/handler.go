package handler

import (
	"encoding/json"
	errs "errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	posRegex = regexp.MustCompile(`[\\/]position[\\/](\S+)$`)
	empRegex = regexp.MustCompile(`[\\/]employee[\\/](\S+)$`)
)

type Hand struct {
	service Service
}

func NewHandler(service Service) *Hand {
	return &Hand{service: service}
}

func (h *Hand) GetPositions(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limitQuery := query["limit"][0]
	limit, err := strconv.ParseInt(limitQuery, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offsetQuery := query["offset"][0]
	offset, err := strconv.ParseInt(offsetQuery, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	positions, err := h.service.GetPositions(r.Context(), int(limit), int(offset))
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
}

func (h *Hand) GetEmployees(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limitQuery := query["limit"][0]
	limit, err := strconv.ParseInt(limitQuery, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offsetQuery := query["offset"][0]
	offset, err := strconv.ParseInt(offsetQuery, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	employees, err := h.service.GetEmployees(r.Context(), int(limit), int(offset))
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
}

func (h *Hand) GetPosition(w http.ResponseWriter, r *http.Request) {
	matches := posRegex.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.service.GetPosition(r.Context(), matches[1])
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
}

func (h *Hand) GetEmployee(w http.ResponseWriter, r *http.Request) {
	matches := empRegex.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	e, err := h.service.GetEmployee(r.Context(), matches[1])
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
	err := h.service.CreatePosition(r.Context(), &p)
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
}

func (h *Hand) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if e.LasName == "" || e.FirstName == "" || e.PositionID == uuid.Nil {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	err := h.service.CreateEmployee(r.Context(), &e)
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
}

func (h *Hand) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdatePosition(r.Context(), &p)
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
}

func (h *Hand) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdateEmployee(r.Context(), &e)
	if err != nil {
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errs.Is(err, errors.PositionIsNotExists()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
}

func (h *Hand) DeletePosition(w http.ResponseWriter, r *http.Request) {
	deleteID := posRegex.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeletePosition(r.Context(), deleteID)
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
}

func (h *Hand) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	deleteID := empRegex.FindStringSubmatch(r.URL.Path)[1]
	err := h.service.DeleteEmployee(r.Context(), deleteID)
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
}
