package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	positionsRe      = regexp.MustCompile(`^\/positions[\/]*$`)
	positionRe       = regexp.MustCompile(`^\/position\/(\d+)$`)
	createPositionRe = regexp.MustCompile(`^\/position$`)
)

type Position struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"'`
	Salary string `json:"salary"`
}
type Employee struct {
	ID, FirstName, LasName string
	PositionID             *Position
}

type datastore struct {
	m map[string]Position
	*sync.RWMutex
}

var positions []Position

type myHandler struct {
	data *datastore
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.RawQuery)
	switch {
	case r.Method == http.MethodGet && positionsRe.MatchString(r.URL.Path):
		h.getPositions(w, r)
		return
	case r.Method == http.MethodGet && positionRe.MatchString(r.URL.Path):
		h.getPosition(w, r)
		return
	case r.Method == http.MethodPost && createPositionRe.MatchString(r.URL.Path):

		h.createPosition(w, r)
		return
	case r.Method == http.MethodDelete && positionRe.MatchString(r.URL.Path):
		h.deletePosition(w, r)
		return
	default:
		fmt.Println("Error: Not Found!")
		http.NotFound(w, r)
		return
	}
}

func (h *myHandler) getPositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.data.RLock()
	users := make([]Position, 0, len(h.data.m))
	for _, v := range h.data.m {
		users = append(users, v)
	}
	h.data.RUnlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *myHandler) getPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	matches := positionRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	h.data.RLock()
	u, ok := h.data.m[matches[1]]
	h.data.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *myHandler) createPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.data.Lock()
	var position Position
	position.ID = strings.Trim(strings.Split(r.URL.RawQuery,"&")[0],"id=")
	position.Name = strings.Trim(strings.Split(r.URL.RawQuery,"&")[1],"name=")
	position.Salary = strings.Trim(strings.Split(r.URL.RawQuery,"&")[2],"salary=")
	h.data.m[position.ID] = position
	h.data.Unlock()
	jsonBytes, err := json.Marshal(position)
	if err != nil {
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *myHandler) updatePosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var position Position
	inputs:=strings.Split(r.URL.RawQuery, "&")
	position.ID = strings.Trim(strings.Split(r.URL.RawQuery,"&")[0],"id=")
	if len(inputs) < 3{
		inputString:=strings.Split(r.URL.RawQuery,"&")[1]
		if inputString == "name=" {

		}else if inputString == "salary="{

		}

	}

	position.Name = strings.Trim(strings.Split(r.URL.RawQuery,"&")[1],"name=")
	position.Salary = strings.Trim(strings.Split(r.URL.RawQuery,"&")[2],"salary=")
	for _, pos := range positions {
		if pos.ID == position.ID {
			pos.Name = position.Name
			pos.Salary = position.Salary
			h.data.Lock()
			h.data.m[pos.ID] = position
			h.data.Unlock()
			jsonBytes, err := json.Marshal(position)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *myHandler) deletePosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var position Position
	if err := json.NewDecoder(r.Body).Decode(&position); err!=nil{
		http.Error(w, "Error with server!", http.StatusInternalServerError)
		return
	}
	for key, pos := range h.data.m {
		if pos.ID == position.ID {
			h.data.Lock()
			delete(h.data.m, key)
			h.data.Unlock()
			jsonBytes, err := json.Marshal(position)
			if err != nil {
				http.Error(w, "Error with server!", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page"))
}

func main() {
	mux := http.NewServeMux()
	myH := &myHandler{
		data: &datastore{
			m: map[string]Position{
				"1": {ID: "1", Name: "Nova", Salary: "100"},
			},
			RWMutex: &sync.RWMutex{},
	},
	}
	//mux.HandleFunc("/", home)
	mux.Handle("/positions", myH)
	mux.Handle("/position/", myH)
	mux.Handle("/position", myH)
	http.ListenAndServe("localhost:8080", mux)
}
