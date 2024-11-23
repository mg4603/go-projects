package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies = make(map[string]Movie)

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getDummyMovies() map[string]Movie {
	return map[string]Movie{
		"1": {Id: "1", Isbn: "438227", Title: "Movie One",
			Director: &Director{FirstName: "John", LastName: "Doe"}},
		"2": {Id: "2", Isbn: "45455", Title: "Movie Two",
			Director: &Director{FirstName: "Jane", LastName: "Doe"}},
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", GetMovies).Methods("GET")
	// r.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	// r.HandleFunc("/movies/", CreateMovie).Methods("POST")
	// r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	// r.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	//
	// r.HandleFunc("/movies", GetMovies).Methods("GET").Handler(SetJSONContentType(http.HandlerFunc(GetMovies)))
	//
	fmt.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
