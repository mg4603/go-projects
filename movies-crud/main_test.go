package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func compareJson(json1, json2 string) (bool, error) {
	var obj1, obj2 map[string]Movie
	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return false, fmt.Errorf("error parsing first JSON: %v\n", err)
	}
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return false, fmt.Errorf("error parsing second JSON: %v\n", err)
	}
	return reflect.DeepEqual(obj1, obj2), nil
}

func TestGetMovies(t *testing.T) {
	tests := []struct {
		name           string
		movies         map[string]Movie
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid movies data",
			movies:         getDummyMovies(),
			expectedStatus: http.StatusOK,
			expectedBody: `{
							"1":{"id":"1","isbn":"438227","title":"Movie One",
								"director":{"first_name":"John","last_name":"Doe"}},
							"2":{"id":"2","isbn":"45455","title":"Movie Two",
								"director":{"first_name":"Jane","last_name":"Doe"}}}`,
		},
		{
			name:           "empty movies data",
			movies:         map[string]Movie{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movies = tt.movies

			req, _ := http.NewRequest("GET", "/movies", nil)
			rec := httptest.NewRecorder()

			GetMovies(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if equal, err := compareJson(tt.expectedBody, rec.Body.String()); !equal || err != nil {
				fmt.Printf("Error parsing json: %v", err)
				t.Errorf("Expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestSetJSONContentType(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedHeader string
	}{
		{
			name: "Set content-type to application/json",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			}),
			expectedHeader: "application/json",
		},
		{
			name: "Overwrite present content-type header",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
			}),
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

			handler := SetJSONContentType(tt.handler)

			handler.ServeHTTP(rec, req)

			log.Println("Final Content-Type header:", rec.Header().Get("Content-Type"))

			if got := rec.Header().Get("Content-Type"); got != tt.expectedHeader {
				t.Errorf("Expected Content-Type %v, but got %v", tt.expectedHeader, got)
			}
		})
	}
}
