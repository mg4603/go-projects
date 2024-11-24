package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type errorResponseWriter struct {
	Recorder *httptest.ResponseRecorder
}

func (e *errorResponseWriter) Header() http.Header {
	return e.Recorder.Header()
}

func (e *errorResponseWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("Forced write error")
}

func (e *errorResponseWriter) WriteHeader(statusCode int) {
	e.Recorder.WriteHeader(statusCode)
}

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

			handler := SetJSONContentType(http.HandlerFunc(GetMovies))
			handler.ServeHTTP(rec, req)

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
func TestGetMovie(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody: `{
							"1":{"id":"1","isbn":"438227","title":"Movie One",
								"director":{"first_name":"John","last_name":"Doe"}}}`,
		},
		{
			name:           "Non-existent ID",
			id:             "3",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 not found\n",
		},
		{
			name:           "Invalid ID format",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movies = getDummyMovies()

			req := httptest.NewRequest(http.MethodGet, "/movies/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rec := httptest.NewRecorder()

			handler := SetJSONContentType(http.HandlerFunc(GetMovie))
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody {
				ok, err := compareJson(tt.expectedBody, rec.Body.String())
				if err != nil {
					t.Error(err)
				}

				if !ok {
					t.Errorf("Expected body %q, got%q", tt.expectedBody, rec.Body.String())
				}
			}

		})
	}
}

func TestRespondInternalServerError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Handles internal server error",
			err:            fmt.Errorf("test error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "An error occurred while processing your request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			RespondInternalServerError(rec, tt.err)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status%v, got %v", tt.expectedStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
			}

		})
	}
}

func TestCreateMovie(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedMovie  *Movie
	}{
		{
			name:           "Valid movie creation",
			body:           `{"isbn":"12345","title":"Test Movie","director":{"first_name":"John","last_name":"Doe"}}`,
			expectedStatus: http.StatusOK,
			expectedMovie: &Movie{
				Isbn:  "12345",
				Title: "Test Movie",
				Director: &Director{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
		},
		{
			name:           "Invalid json in body",
			body:           `{"isbn":"12345","title":"Test Movie","director":{`,
			expectedStatus: http.StatusInternalServerError,
			expectedMovie:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movies = getDummyMovies()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(tt.body))

			handler := SetJSONContentType(http.HandlerFunc(CreateMovie))
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedMovie != nil {
				var got Movie

				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if got.Title != tt.expectedMovie.Title ||
					got.Isbn != tt.expectedMovie.Isbn ||
					got.Director.FirstName != tt.expectedMovie.Director.FirstName ||
					got.Director.LastName != tt.expectedMovie.Director.LastName {

					t.Errorf("Expected movie %+v, got %+v", tt.expectedMovie, got)
				}
			}
		})
	}
}
