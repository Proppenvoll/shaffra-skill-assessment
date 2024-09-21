package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
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

func TestGetLoggingHandler(t *testing.T) {
	handler := getLoggingHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	http.Get(server.URL)
	log.SetOutput(os.Stdout)

	expected := regexp.MustCompile(`\d+\/\d+\/\d+ \d+:\d+:\d+ \w+ \/ took \w*`)
	result := logBuffer.String()

	if !expected.MatchString(result) {
		t.Errorf(wantGot, expected.String(), result)
	}
}
