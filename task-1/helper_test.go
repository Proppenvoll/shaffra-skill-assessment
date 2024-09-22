package main

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeJson(t *testing.T) {
	type Test struct {
		Test int `json:"test"`
	}

	t.Run("returns a result and no error", func(t *testing.T) {
		reader := strings.NewReader("{\"test\":1}\n")
		comparison := Test{1}
		result, error := decodeJson[Test](reader)

		if error != nil {
			t.Error()
		}

		if !reflect.DeepEqual(comparison, *result) {
			t.Errorf(wantGot, comparison, *result)
		}
	})

	t.Run("returns no result and an error", func(t *testing.T) {
		reader := strings.NewReader("{malformed}\n")
		result, error := decodeJson[Test](reader)

		if error == nil {
			t.Errorf(wantGot, nil, error)
		}

		if result != nil {
			t.Errorf(wantGot, nil, result)
		}
	})
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

	a := httptest.NewRecorder()

	if error = encodeJson(a, new(complex64)); error == nil {
		t.Errorf(wantGot, "not nil", error)
	}
}
