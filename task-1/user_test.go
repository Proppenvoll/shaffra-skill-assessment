package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"fmt"
	"testing"
)

func getSimulateHandlerInteraction(
	handler func(http.ResponseWriter, *http.Request),
) func(payload string) *httptest.ResponseRecorder {
	return func(payload string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		reader := strings.NewReader(payload)
		r := httptest.NewRequest("", "http://test", reader)
		r.SetPathValue("id", "1")
		handler(w, r)
		return w
	}
}

func TestCreateUser(t *testing.T) {
	const generatedId = 1

	handler := createUser(func(user UserWithoutId) (int, error) {
		return generatedId, nil
	})

	simulateCreateUserInteraction := getSimulateHandlerInteraction(handler)

	assertErrors := func(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int) {
		bodyLength := w.Body.Len()
		if bodyLength <= 0 {
			t.Errorf(wantGot, ">0", bodyLength)
		}

		statusCode := w.Result().StatusCode
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	}

	t.Run("responds with an error for invalid json payload", func(t *testing.T) {
		w := simulateCreateUserInteraction("invalid")
		assertErrors(t, w, 400)
	})

	t.Run("responds with an error if validation fails", func(t *testing.T) {
		w := simulateCreateUserInteraction(`{"name":"name"}`)
		assertErrors(t, w, 400)
	})

	t.Run("responds with an error if query fails", func(t *testing.T) {
		handler := createUser(func(user UserWithoutId) (int, error) {
			return 0, errors.New("test error")
		})

		simulateCreateUserInteraction := getSimulateHandlerInteraction(handler)
		w := simulateCreateUserInteraction(`{"name":"name", "email":"a@b.com", "age": 30}`)
		assertErrors(t, w, 500)
	})

	t.Run("responds with the generated id", func(t *testing.T) {
		w := simulateCreateUserInteraction(`{"name":"name", "email":"a@b.com", "age": 30}`)
		body := fmt.Sprintf("%s", w.Body)
		expectedBody := "{\"id\":1}\n"
		if body != expectedBody {
			t.Errorf(wantGot, expectedBody, body)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("responds with an error when user not found", func(t *testing.T) {
		handler := getUser("id", func(id string) (User, error) {
			return User{}, errors.New("test error")
		})

		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction("")

		body := fmt.Sprintf("%s", w.Body)
		expectedBody := "user not found\n"
		if body != expectedBody {
			t.Errorf(wantGot, expectedBody, body)
		}

		statusCode := w.Result().StatusCode
		expectedStatusCode := 404
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	})

	t.Run("responds with the user", func(t *testing.T) {
		userToGet := User{1, UserWithoutId{"name", "a@b.com", 30}}
		handler := getUser("id", func(id string) (User, error) {
			return userToGet, nil
		})

		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction("")
		body := fmt.Sprintf("%s", w.Body)

		expectedBodyRaw, _ := json.Marshal(userToGet)
		expectedBody := string(expectedBodyRaw) + "\n"

		if body != string(expectedBody) {
			t.Errorf(wantGot, expectedBody, body)
		}
	})
}

func TestReplaceUser(t *testing.T) {
	handler := replaceUser("id", func(user User) error {
		return nil
	})

	simulateHandlerInteraction := getSimulateHandlerInteraction(handler)

	assertErrors := func(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int) {
		bodyLength := w.Body.Len()
		if bodyLength <= 0 {
			t.Errorf(wantGot, ">0", bodyLength)
		}

		statusCode := w.Result().StatusCode
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	}

	t.Run("responds with error query parameter id cannot be converted to int", func(t *testing.T) {
		handler := replaceUser("fail", func(user User) error {
			return nil
		})

		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction(`{"name":"name", "email":"a@b.com", "age": 30}`)
		assertErrors(t, w, 400)
	})

	t.Run("responds with error when json is malformed", func(t *testing.T) {
		w := simulateHandlerInteraction(`{name:"name"}`)
		assertErrors(t, w, 400)
	})

	t.Run("responds with error when validation fails", func(t *testing.T) {
		w := simulateHandlerInteraction(`{"name":"name"}`)
		assertErrors(t, w, 400)
	})

	t.Run("responds with error when query fails", func(t *testing.T) {
		handler := replaceUser("id", func(user User) error {
			return errors.New("test error")
		})

		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction(`{"name":"name", "email":"a@b.com", "age": 30}`)
		assertErrors(t, w, 500)
	})

	t.Run("responds with ok", func(t *testing.T) {
		w := simulateHandlerInteraction(`{"name":"name", "email":"a@b.com", "age": 30}`)
		statusCode := w.Result().StatusCode
		expectedStatusCode := 200
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("responds with error when query fails", func(t *testing.T) {
		const testError = "test error"
		handler := deleteUser("id", func(id string) error {
			return errors.New(testError)
		})
		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction("")

		body := fmt.Sprintf("%s", w.Body)
		expectedBody := testError + "\n"
		if body != expectedBody {
			t.Errorf(wantGot, expectedBody, body)
		}

		statusCode := w.Result().StatusCode
		expectedStatusCode := 500
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	})

	t.Run("responds with ok", func(t *testing.T) {
		handler := deleteUser("id", func(id string) error {
			return nil
		})
		simulateHandlerInteraction := getSimulateHandlerInteraction(handler)
		w := simulateHandlerInteraction("")
		statusCode := w.Result().StatusCode
		expectedStatusCode := 200
		if statusCode != expectedStatusCode {
			t.Errorf(wantGot, expectedStatusCode, statusCode)
		}
	})
}
