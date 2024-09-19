package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
		createUser(db),
	)

	http.HandleFunc(
		getUserByIdPattern("GET", pathId),
		getUser(pathId, db),
	)

	http.HandleFunc(
		getUserByIdPattern("PUT", pathId),
		replaceUser(pathId, db),
	)

	http.HandleFunc(
		getUserByIdPattern("DELETE", pathId),
		deleteUser(pathId, db),
	)

	serverAddress := ":8080"
	log.Println("Starting server on", serverAddress)
	loggingHandler := getLoggingHandler(http.DefaultServeMux)
	log.Fatalln(http.ListenAndServe(serverAddress, loggingHandler))
}
