package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var error error
	db, error = sql.Open("postgres", "user=postgres password=pass dbname=test sslmode=disable")

	if error != nil {
		log.Fatal(error)
	}

	http.HandleFunc("GET /users", getUsers)
	http.HandleFunc("POST /users", createUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, error := db.Query("SELECT users_id, name FROM users")
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	var responsePayload []User

	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		responsePayload = append(responsePayload, user)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Provide Content-Type application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user struct {
		Name string `json:"name"`
	}

	error := json.NewDecoder(r.Body).Decode(&user)

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	if user.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	username := user.Name

	time.Sleep(5 * time.Second) // Simulate a long database operation

	var userId int
	error = db.QueryRow(
		"INSERT INTO users (name) VALUES ($1) RETURNING users_id",
		username,
	).Scan(&userId)

	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Id int `json:"id"`
	}{userId})
}
