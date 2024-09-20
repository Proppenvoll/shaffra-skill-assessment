package main

import (
	"net/http"
	"strconv"
)

func createUser(queryCreateUser queryCreateUser) func(http.ResponseWriter, *http.Request) {
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

		userId, error := queryCreateUser(*payload)

		if error != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		error = encodeJson(w, struct {
			Id int `json:"id"`
		}{userId})
	}
}

func getUser(pathId string, queryGetUser queryGetUser) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(pathId)
		user, error := queryGetUser(id)

		if error != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		encodeJson(w, user)
	}
}

func replaceUser(
	pathId string,
	queryReplaceUser queryReplaceUser,
) func(http.ResponseWriter, *http.Request) {
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
		error = queryReplaceUser(user)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteUser(
	pathId string,
	queryDeleteUser queryDeleteUser,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(pathId)
		error := queryDeleteUser(id)

		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
		}
	}
}
