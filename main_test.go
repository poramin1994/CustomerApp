package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("test_customers.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Reset the database for testing
	db.Exec("DROP TABLE IF EXISTS customers")
	db.AutoMigrate(&Customer{})

	// Insert initial data
	db.Create(&Customer{Name: "John Doe", Age: 30})
	db.Create(&Customer{Name: "Jane Smith", Age: 25})
}

func TestCreateCustomer(t *testing.T) {
	initTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/customers", createCustomer).Methods("POST")

	customer := &Customer{Name: "Test User", Age: 40}
	body, _ := json.Marshal(customer)
	req, err := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var newCustomer Customer
	json.NewDecoder(rr.Body).Decode(&newCustomer)
	if newCustomer.Name != customer.Name || newCustomer.Age != customer.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), body)
	}
}

func TestGetCustomer(t *testing.T) {
	initTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/customers/{id}", getCustomer).Methods("GET")

	req, err := http.NewRequest("GET", "/customers/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var customer Customer
	json.NewDecoder(rr.Body).Decode(&customer)
	if customer.ID != 1 || customer.Name != "John Doe" || customer.Age != 30 {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), `{"id":1,"name":"John Doe","age":30}`)
	}
}

func TestUpdateCustomer(t *testing.T) {
	initTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")

	customer := &Customer{Name: "Updated User", Age: 35}
	body, _ := json.Marshal(customer)
	req, err := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedCustomer Customer
	json.NewDecoder(rr.Body).Decode(&updatedCustomer)
	if updatedCustomer.Name != customer.Name || updatedCustomer.Age != customer.Age {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), body)
	}
}

func TestDeleteCustomer(t *testing.T) {
	initTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/customers/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
