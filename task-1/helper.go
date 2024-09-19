package main

import (
	"encoding/json"
	"io"
	"net/http"
)

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
