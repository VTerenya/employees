package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VTerenya/employees/internal"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

var (
	data        *repository.Database   //nolint: gochecknoglobals
	repos       *repository.Repository //nolint: gochecknoglobals
	serv        *service.Serv          //nolint: gochecknoglobals
	handler     *Hand                  //nolint: gochecknoglobals
	positionIDs []string               //nolint: gochecknoglobals
	employeeIDs []string               //nolint: gochecknoglobals
)

func initTest() {
	data = repository.NewDataBase()
	repos = repository.NewRepo(data)
	serv = service.NewServ(repos)
	handler = NewHandler(serv)
	positionIDs = make([]string, 0)
	employeeIDs = make([]string, 0)
}

func createPosID() uuid.UUID {
	id := uuid.New()
	positionIDs = append(positionIDs, id.String())
	return id
}

func createEmpID() uuid.UUID {
	id := uuid.New()
	employeeIDs = append(employeeIDs, id.String())
	return id
}

func createTestContext(r *http.Request) *http.Request {
	ctx := r.Context()
	id := uuid.New()
	//revive:disable
	ctx = context.WithValue(ctx, "correlation_id", id.String()) //nolint:staticcheck
	//revive:enable
	r = r.WithContext(ctx)
	return r
}

func TestHand_CreateEmployee(t *testing.T) {
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	id, err := uuid.Parse(positionIDs[0])
	if err != nil {
		t.Fatalf("Error: %v ", err)
	}
	firstEmployee := internal.Employee{ID: createEmpID(), FirstName: "Bob", LasName: "Vik", PositionID: id}
	secondEmployee := internal.Employee{ID: createEmpID(), FirstName: "", LasName: "Vik", PositionID: id}
	thirdEmployee := internal.Employee{ID: createEmpID(), FirstName: "Fox", LasName: "Fok", PositionID: uuid.New()}
	jsonFirstEmployee, _ := json.Marshal(firstEmployee)
	jsonSecondEmployee, _ := json.Marshal(secondEmployee)
	jsonThirdEmployee, _ := json.Marshal(thirdEmployee)
	reader := strings.NewReader(string(jsonFirstEmployee))
	testTable := []struct {
		URL      string
		method   string
		expected int
		read     io.Reader
	}{
		{
			URL:      "http://localhost:8080/employee",
			method:   "POST",
			expected: 200,
			read:     reader,
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "POST",
			expected: 500,
			read:     strings.NewReader("s"),
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "POST",
			expected: 400,
			read:     strings.NewReader(string(jsonThirdEmployee)),
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "POST",
			expected: 400,
			read:     strings.NewReader(string(jsonSecondEmployee)),
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest("POST", "http://localhost:8080/employee", testCase.read)
		r.Header.Set("content-type", "application/json")
		r = createTestContext(r)
		if err != nil {
			t.Fatalf("could not created request: %v", err)
		}
		w := httptest.NewRecorder()
		handler.CreateEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected status created; get %v", result.StatusCode)
		}
	}
}

func TestHand_CreatePosition(t *testing.T) {
	initTest()
	firstPosition := internal.Position{Salary: decimal.New(500, 0), Name: "worker"}
	fakePosition := internal.Employee{ID: createEmpID(), FirstName: "V", LasName: "T", PositionID: uuid.New()}
	jsonFirstPosition, _ := json.Marshal(firstPosition)
	jsonFakePosition, _ := json.Marshal(fakePosition)
	reader := strings.NewReader(string(jsonFirstPosition))
	testTable := []struct {
		URL      string
		method   string
		expected int
		read     io.Reader
	}{
		{
			URL:      "http://localhost:8080/position",
			method:   "POST",
			expected: 500,
			read:     strings.NewReader("s"),
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "POST",
			expected: 200,
			read:     reader,
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "POST",
			expected: 400,
			read:     strings.NewReader(string(jsonFakePosition)),
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "POST",
			expected: 400,
			read:     strings.NewReader(string(jsonFirstPosition)),
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest("POST", "http://localhost:8080/position", testCase.read)
		r.Header.Set("content-type", "application/json")
		r = createTestContext(r)
		if err != nil {
			t.Fatalf("could not created request: %v", err)
		}
		w := httptest.NewRecorder()
		handler.CreatePosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected status created; get %v", result.StatusCode)
		}
	}
}

func TestHand_DeleteEmployeeVarsZero(t *testing.T) {
	initTest()
	posID := createPosID()
	position := internal.Position{ID: posID, Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&position)
	employee := internal.Employee{ID: createEmpID(), FirstName: "Vik", LasName: "Vok", PositionID: posID}
	repos.AddEmployee(&employee)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		employeeID string
	}{

		{
			URL:        "http://localhost:8080/employee/" + employeeIDs[0],
			method:     "DELETE",
			employeeID: employeeIDs[0],
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.DeleteEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_DeleteEmployee(t *testing.T) {
	initTest()
	posID := createPosID()
	position := internal.Position{ID: posID, Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&position)
	employee := internal.Employee{ID: createEmpID(), FirstName: "Vik", LasName: "Vok", PositionID: posID}
	repos.AddEmployee(&employee)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		employeeID string
	}{
		{
			URL:        "http://localhost:8080/employee/" + employeeIDs[0],
			method:     "DELETE",
			employeeID: employeeIDs[0],
			expected:   200,
		},
		{
			URL:        "http://localhost:8080/employee/" + uuid.New().String(),
			method:     "DELETE",
			employeeID: uuid.New().String(),
			expected:   404,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		vars := map[string]string{
			"id": testCase.employeeID,
		}
		r = createTestContext(mux.SetURLVars(r, vars))
		w := httptest.NewRecorder()
		handler.DeleteEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_DeletePositionVarsZero(t *testing.T) { //nolint: funlen
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		positionID string
	}{

		{
			URL:        "http://localhost:8080/position/" + positionIDs[0],
			method:     "DELETE",
			positionID: positionIDs[0],
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.DeletePosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}
func TestHand_DeletePosition(t *testing.T) {
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		positionID string
	}{
		{
			URL:        "http://localhost:8080/position/" + positionIDs[0],
			method:     "DELETE",
			positionID: positionIDs[0],
			expected:   200,
		},
		{
			URL:        "http://localhost:8080/position/" + uuid.New().String(),
			method:     "DELETE",
			positionID: uuid.New().String(),
			expected:   404,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		vars := map[string]string{
			"id": testCase.positionID,
		}
		r = createTestContext(mux.SetURLVars(r, vars))
		w := httptest.NewRecorder()
		handler.DeletePosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetEmployeeVarsZero(t *testing.T) {
	initTest()
	posID := createPosID()
	position := internal.Position{ID: posID, Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&position)
	employee := internal.Employee{ID: createEmpID(), FirstName: "Vik", LasName: "Vok", PositionID: posID}
	repos.AddEmployee(&employee)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		employeeID string
	}{
		{
			URL:        "http://localhost:8080/employee/" + employeeIDs[0],
			method:     "GET",
			employeeID: employeeIDs[0],
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.GetEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetEmployee(t *testing.T) {
	initTest()
	posID := createPosID()
	position := internal.Position{ID: posID, Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&position)
	employee := internal.Employee{ID: createEmpID(), FirstName: "Vik", LasName: "Vok", PositionID: posID}
	repos.AddEmployee(&employee)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		employeeID string
	}{
		{
			URL:        "http://localhost:8080/employee/" + employeeIDs[0],
			method:     "GET",
			employeeID: employeeIDs[0],
			expected:   200,
		},
		{
			URL:        "http://localhost:8080/employee/" + uuid.New().String(),
			method:     "Get",
			employeeID: uuid.New().String(),
			expected:   404,
		},
		{
			URL:        "http://localhost:8080/employee/" + "12",
			method:     "Get",
			employeeID: "12",
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		vars := map[string]string{
			"id": testCase.employeeID,
		}
		r = createTestContext(mux.SetURLVars(r, vars))
		w := httptest.NewRecorder()
		handler.GetEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetPositionVarsZero(t *testing.T) {
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		positionID string
	}{
		{
			URL:        "http://localhost:8080/position/" + positionIDs[0],
			method:     "GET",
			positionID: positionIDs[0],
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.GetPosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetPosition(t *testing.T) {
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	testTable := []struct {
		URL        string
		method     string
		expected   int
		positionID string
	}{
		{
			URL:        "http://localhost:8080/position/" + positionIDs[0],
			method:     "GET",
			positionID: positionIDs[0],
			expected:   200,
		},
		{
			URL:        "http://localhost:8080/position/" + uuid.New().String(),
			method:     "Get",
			positionID: uuid.New().String(),
			expected:   404,
		},
		{
			URL:        "http://localhost:8080/position/" + "12",
			method:     "Get",
			positionID: "12",
			expected:   400,
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		vars := map[string]string{
			"id": testCase.positionID,
		}
		r = createTestContext(mux.SetURLVars(r, vars))
		w := httptest.NewRecorder()
		handler.GetPosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetEmployees(t *testing.T) {
	initTest()
	testTable := []struct {
		URL      string
		method   string
		expected int
	}{
		{
			URL:      "http://localhost:8080/employees?limit=1&offset=1",
			method:   "GET",
			expected: 200,
		},
		{
			URL:      "http://localhost:8080/employees?limit=1&offset=1",
			method:   "GET",
			expected: 200,
		},
		{
			URL:      "http://localhost:8080/employees?limit=asd&offset=1",
			method:   "GET",
			expected: 500,
		},
		{
			URL:      "http://localhost:8080/employees?limit=1&offset=asd",
			method:   "GET",
			expected: 500,
		},
		{
			URL:      "http://localhost:8080/employees?limit=1&offset=3",
			method:   "GET",
			expected: 404,
		},
		{
			URL:      "http://localhost:8080/employees?limit=110&offset=3",
			method:   "GET",
			expected: 500,
		},
	}

	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.GetEmployees(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_GetPositions(t *testing.T) {
	initTest()
	testTable := []struct {
		URL      string
		method   string
		expected int
	}{
		{
			URL:      "http://localhost:8080/positions?limit=1&offset=1",
			method:   "GET",
			expected: 200,
		},
		{
			URL:      "http://localhost:8080/positions?limit=1&offset=1",
			method:   "GET",
			expected: 200,
		},
		{
			URL:      "http://localhost:8080/positions?limit=asd&offset=1",
			method:   "GET",
			expected: 500,
		},
		{
			URL:      "http://localhost:8080/positions?limit=1&offset=asd",
			method:   "GET",
			expected: 500,
		},
		{
			URL:      "http://localhost:8080/positions?limit=110&offset=3",
			method:   "GET",
			expected: 500,
		},
	}

	for _, testCase := range testTable {
		r, err := http.NewRequest(testCase.method, testCase.URL, nil)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		r = createTestContext(r)
		w := httptest.NewRecorder()
		handler.GetPositions(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected %v; get %v", testCase.expected, result)
		}
	}
}

func TestHand_UpdateEmployee(t *testing.T) { //nolint: funlen
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	id, err := uuid.Parse(positionIDs[0])
	if err != nil {
		t.Fatalf("Error: %v ", err)
	}
	firstEmployee := internal.Employee{ID: createEmpID(), FirstName: "Bob", LasName: "Vik", PositionID: id}
	firstEmployeeID := firstEmployee.ID
	repos.AddEmployee(&firstEmployee)
	secondEmployee := internal.Employee{ID: firstEmployeeID, FirstName: "Victor", LasName: "Vik", PositionID: id}
	thirdEmployee := internal.Employee{ID: createEmpID(), FirstName: "Fox", LasName: "Fok", PositionID: uuid.New()}
	fourthEmployee := internal.Employee{ID: uuid.Nil, FirstName: "", LasName: "Fok", PositionID: uuid.New()}
	fifthEmployee := internal.Employee{ID: firstEmployeeID, FirstName: "Chu", LasName: "Xi", PositionID: uuid.New()}
	jsonSecondEmployee, _ := json.Marshal(secondEmployee)
	jsonThirdEmployee, _ := json.Marshal(thirdEmployee)
	jsonFourthEmployee, _ := json.Marshal(fourthEmployee)
	jsonFifthEmployee, _ := json.Marshal(fifthEmployee)
	reader := strings.NewReader(string(jsonSecondEmployee))
	testTable := []struct {
		URL      string
		method   string
		expected int
		read     io.Reader
	}{
		{
			URL:      "http://localhost:8080/employee",
			method:   "PUT",
			expected: 200,
			read:     reader,
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "PUT",
			expected: 500,
			read:     strings.NewReader("s"),
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "PUT",
			expected: 404,
			read:     strings.NewReader(string(jsonThirdEmployee)),
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "PUT",
			expected: 200,
			read:     strings.NewReader(string(jsonSecondEmployee)),
		},
		{
			URL:      "localhost:8080/employee",
			method:   "PUT",
			expected: 400,
			read:     strings.NewReader(string(jsonFourthEmployee)),
		},
		{
			URL:      "http://localhost:8080/employee",
			method:   "PUT",
			expected: 400,
			read:     strings.NewReader(string(jsonFifthEmployee)),
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest("PUT", "http://localhost:8080/employee", testCase.read)
		r.Header.Set("content-type", "application/json")
		r = createTestContext(r)
		if err != nil {
			t.Fatalf("could not created request: %v", err)
		}
		w := httptest.NewRecorder()
		handler.UpdateEmployee(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected status %v; get %v", testCase.expected, result.StatusCode)
		}
	}
}

func TestHand_UpdatePosition(t *testing.T) {
	initTest()
	firstPosition := internal.Position{ID: createPosID(), Salary: decimal.New(500, 0), Name: "worker"}
	repos.AddPosition(&firstPosition)
	firstID := firstPosition.ID
	secondPosition := internal.Position{ID: firstID, Salary: decimal.New(1000, 0), Name: "worker"}
	fakePosition := internal.Employee{ID: uuid.Nil, FirstName: "V", LasName: "T", PositionID: uuid.Nil}
	fakeSecondPosition := internal.Employee{ID: createEmpID(), FirstName: "V", LasName: "T", PositionID: uuid.Nil}
	jsonSecondPosition, _ := json.Marshal(secondPosition)
	jsonFakePosition, _ := json.Marshal(fakePosition)
	jsonFakeSecondPosition, _ := json.Marshal(fakeSecondPosition)
	reader := strings.NewReader(string(jsonFakePosition))
	testTable := []struct {
		URL      string
		method   string
		expected int
		read     io.Reader
	}{
		{
			URL:      "http://localhost:8080/position",
			method:   "PUT",
			expected: 500,
			read:     strings.NewReader("s"),
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "PUT",
			expected: 400,
			read:     reader,
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "PUT",
			expected: 404,
			read:     strings.NewReader(string(jsonFakeSecondPosition)),
		},
		{
			URL:      "http://localhost:8080/position",
			method:   "PUT",
			expected: 200,
			read:     strings.NewReader(string(jsonSecondPosition)),
		},
	}
	for _, testCase := range testTable {
		r, err := http.NewRequest("PUT", "http://localhost:8080/position", testCase.read)
		r.Header.Set("content-type", "application/json")
		r = createTestContext(r)
		if err != nil {
			t.Fatalf("could not created request: %v", err)
		}
		w := httptest.NewRecorder()
		handler.UpdatePosition(w, r)
		result := w.Result()
		defer result.Body.Close()
		if result.StatusCode != testCase.expected {
			t.Fatalf("expected status created; get %v", result.StatusCode)
		}
	}
}
