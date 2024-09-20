package main

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

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

	error := encodeJson(w, struct {
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
