package main

import (
	"database/sql"
	"net/http"
	"strconv"
)

func createUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, error := decodeJson[UserWithoutId](r.Body)

		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		error = payload.validate()

		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		var userId int
		error = db.QueryRow(
			"INSERT INTO app_user (name, email, age) VALUES ($1, $2, $3) RETURNING app_user_id",
			payload.Name,
			payload.Email,
			payload.Age,
		).Scan(&userId)

		if error != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		endcodeJson(w, struct {
			Id int `json:"id"`
		}{userId})
	}
}

func getUser(pathId string, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(pathId)

		var user User
		error := db.QueryRow(
			"SELECT * FROM app_user WHERE app_user_id = $1",
			id,
		).Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Age,
		)
		if error != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		endcodeJson(w, user)
	}
}

func replaceUser(pathId string, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rawId := r.PathValue(pathId)
		id, error := strconv.Atoi(rawId)

		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		payload, error := decodeJson[UserWithoutId](r.Body)

		if error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		if error = payload.validate(); error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		user := User{id, *payload}

		_, error = db.Exec(
			`INSERT INTO app_user (app_user_id, name, email, age)
                 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (app_user_id) DO UPDATE SET
                 name = $2, email = $3, age = $4`,
			user.Id,
			user.Name,
			user.Email,
			user.Age,
		)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteUser(pathId string, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(pathId)

		_, error := db.Exec("DELETE FROM app_user WHERE app_user_id = $1", id)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
		}
	}
}
