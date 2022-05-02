package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database Connection
const (
	host = "127.0.0.1"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname = "postgres"
)

// Model 
type Person struct {
	gorm.Model
	Name string `json:"name"`
	Nickname string `json:"nickname"`
}


func main(){
	http.HandleFunc("/people", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))	
}

func GETHandler(w http.ResponseWriter, r *http.Request){
	db := OpenConnection()

	var people []Person
	
	db.Find(&people)

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)
}

func POSTHandler(w http.ResponseWriter, r *http.Request){
	db := OpenConnection()
	var person Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	db.Create(&person)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

}

func OpenConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s " + 
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Person{})

	return db
}