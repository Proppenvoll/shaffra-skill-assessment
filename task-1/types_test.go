package main

import "testing"

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
