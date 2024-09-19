package main

import (
	"errors"
	"net/mail"
)

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
