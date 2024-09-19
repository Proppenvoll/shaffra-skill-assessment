package main

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const wantGot = "want: %v, got: %v"

func TestGetUsersPattern(t *testing.T) {
	patterns := []struct {
		verb       string
		comparison string
	}{
		{"GET", "GET /users"},
		{"POST", "POST /users"},
	}

	for _, pattern := range patterns {
		userPattern := getUsersPattern(pattern.verb)
		if userPattern != pattern.comparison {
			t.Errorf(wantGot, pattern.comparison, userPattern)
		}
	}
}

func TestGetUserByIdPattern(t *testing.T) {
	patterns := []struct {
		verb       string
		id         string
		comparison string
	}{
		{"GET", "id1", "GET /users/{id1}"},
		{"POST", "id2", "POST /users/{id2}"},
	}

	for _, pattern := range patterns {
		userByIdPattern := getUserByIdPattern(pattern.verb, pattern.id)
		if userByIdPattern != pattern.comparison {
			t.Errorf(wantGot, pattern.comparison, userByIdPattern)
		}
	}
}

func TestDecodeJson(t *testing.T) {
	t.Run("returns a result and no error", func(t *testing.T) {
		reader := strings.NewReader("{\"test\":1}\n")

		type Test struct {
			Test int `json:"test"`
		}

		comparison := Test{1}

		result, error := decodeJson[Test](reader)

		if error != nil {
			t.Error()
		}

		if !reflect.DeepEqual(comparison, *result) {
			t.Errorf(wantGot, comparison, *result)
		}
	})

	t.Run("returns no result and an error", func(t *testing.T) {})
}

func TestEncodeJson(t *testing.T) {
	w := httptest.NewRecorder()

	error := endcodeJson(w, struct {
		Test int `json:"test"`
	}{1})

	if error != nil {
		t.Error()
	}

	contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
	if contentType != expectedContentType {
		t.Errorf(wantGot, expectedContentType, contentType)
	}

	jsonString := fmt.Sprintf("%s", w.Body)
	expectedJsonString := "{\"test\":1}\n"
	if jsonString != expectedJsonString {
		t.Errorf(wantGot, expectedJsonString, jsonString)
	}
}

func TestValidate(t *testing.T) {
	userWithoutId := UserWithoutId{}

	error := userWithoutId.validate().Error()
	expected := "Missing name entry"
	if error != expected {
		t.Errorf(wantGot, expected, error)
	}

	userWithoutId.Name = "name"

	error = userWithoutId.validate().Error()
	expected = "Invalid email address"
	if error != expected {
		t.Errorf(wantGot, expected, error)
	}

	userWithoutId.Email = "a@b.de"

	error = userWithoutId.validate().Error()
	expected = "Invalid age"
	if error != expected {
		t.Errorf(wantGot, expected, error)
	}

	userWithoutId.Age = 3

	if error := userWithoutId.validate(); error != nil {
		t.Error()
	}
}
