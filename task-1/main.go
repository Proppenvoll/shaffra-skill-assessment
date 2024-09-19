package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func getUsersPattern(verb string) string {
	return fmt.Sprintf("%s /users", verb)
}

func getUserByIdPattern(verb string, id string) string {
	usersPattern := getUsersPattern(verb)
	return fmt.Sprintf("%s/{%s}", usersPattern, id)
}

func getLoggingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		end := time.Since(start)
		log.Printf("%s %s took %v\n", r.Method, r.URL, end)
	})
}

type UserWithoutId struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type User struct {
	Id int `json:"id"`
	UserWithoutId
}

func (userWithoutId *UserWithoutId) validate() error {
	if userWithoutId.Name == "" {
		return errors.New("Missing name entry")
	}

	if _, error := mail.ParseAddress(userWithoutId.Email); error != nil {
		return errors.New("Invalid email address")
	}

	if userWithoutId.Age == 0 {
		return errors.New("Invalid age")
	}

	return nil
}

func decodeJson[T any](reader io.Reader) (*T, error) {
	result := new(T)
	if decodeError := json.NewDecoder(reader).Decode(result); decodeError != nil {
		return nil, decodeError
	}
	return result, nil
}

func endcodeJson(writer http.ResponseWriter, payload any) error {
	writer.Header().Set("Content-Type", "application/json")
	if error := json.NewEncoder(writer).Encode(payload); error != nil {
		return error
	}
	return nil
}

func main() {
	db, error := sql.Open(
		"postgres",
		"postgresql://admin:admin@localhost:5432?sslmode=disable",
	)

	if error != nil {
		log.Fatal(error)
	}

	const pathId = "id"

	http.HandleFunc(
		getUsersPattern("POST"),
		func(w http.ResponseWriter, r *http.Request) {
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
		},
	)

	http.HandleFunc(
		getUserByIdPattern("GET", pathId),
		func(w http.ResponseWriter, r *http.Request) {
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
		},
	)

	http.HandleFunc(
		getUserByIdPattern("PUT", pathId),
		func(w http.ResponseWriter, r *http.Request) {
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
				sql.Named("id", user.Id),
				sql.Named("name", user.Name),
				sql.Named("email", user.Email),
				sql.Named("age", user.Age),
			)

			if error != nil {
				http.Error(w, error.Error(), http.StatusInternalServerError)
			}
		},
	)

	http.HandleFunc(
		getUserByIdPattern("DELETE", pathId),
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue(pathId)

			_, error := db.Exec("DELETE FROM app_user WHERE app_user_id = $1", id)

			if error != nil {
				http.Error(w, error.Error(), http.StatusInternalServerError)
			}
		},
	)

	serverAddress := ":8080"
	log.Println("Starting server on", serverAddress)
	loggingHandler := getLoggingHandler(http.DefaultServeMux)
	log.Fatalln(http.ListenAndServe(serverAddress, loggingHandler))
}
