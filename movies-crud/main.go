package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

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

func RespondInternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal Server Error: %v", err)
	http.Error(w, "An error occurred while processing your request", http.StatusInternalServerError)
}

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "application/json")
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

func GetMovies(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		RespondInternalServerError(w, err)
		return
	}
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id string
	if id = params["id"]; !regexp.MustCompile(`^\d+$`).MatchString(id) {
		log.Printf("Invalid ID format: %v", id)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	movie, ok := movies[id]
	if !ok {
		log.Printf("Movie with id %v does not exist", id)
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	prepared_obj := map[string]Movie{id: movie}
	if err := json.NewEncoder(w).Encode(prepared_obj); err != nil {
		RespondInternalServerError(w, err)
		return
	}

}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		RespondInternalServerError(w, err)
		return
	}

	id := strconv.Itoa(rand.Intn(1000000))
	movie.Id = id
	movies[id] = movie

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		RespondInternalServerError(w, err)
		return
	}
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var movie Movie

	var id string

	if id = params["id"]; !regexp.MustCompile(`^\d+$`).MatchString(id) {
		log.Printf("Invalid ID format: %v", id)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if _, ok := movies[id]; !ok {
		log.Printf("Movie with id: %v does not exist", id)
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		RespondInternalServerError(w, err)
		return
	}

	movie.Id = id
	movies[id] = movie

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		RespondInternalServerError(w, err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	r.HandleFunc("/movies/", CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	// r.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	//

	r.HandleFunc("/movies", GetMovies).Methods("GET").Handler(SetJSONContentType(http.HandlerFunc(GetMovies)))
	r.HandleFunc("/movies/{id}", GetMovie).Methods("GET").Handler(SetJSONContentType(http.HandlerFunc(GetMovie)))
	r.HandleFunc("/movies", CreateMovie).Methods("POST").Handler(SetJSONContentType(http.HandlerFunc(CreateMovie)))
	r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT").Handler(SetJSONContentType(http.HandlerFunc(UpdateMovie)))
	fmt.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
