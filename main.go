package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Customer struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("customers.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the Customer struct to create/update the customers table
	db.AutoMigrate(&Customer{})

	// Insert initial data
	db.Create(&Customer{Name: "John Doe", Age: 30})
	db.Create(&Customer{Name: "Jane Smith", Age: 25})
}

func main() {
	initDB()
	r := mux.NewRouter()
	r.HandleFunc("/customers", createCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&customer)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var customer Customer
	result := db.First(&customer, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(customer)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var customer Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Model(&Customer{}).Where("id = ?", id).Updates(customer)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	customer.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result := db.Delete(&Customer{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
