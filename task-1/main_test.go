package main

import (
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
