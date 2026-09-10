package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"` // *Director is a pointer to an instance of the Director struct
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applications/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.IntN(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content-Type", "applications/json")

	// params
	params := mux.Vars(r)

	// loop over the movies, range
	// delete the movie with the ID that you've sent
	// add a new movie - the movie that we send in the http request body
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = strconv.Itoa(rand.IntN(100000000))
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) // all the movies that come after this index now replace the movie at this index
			break
		}
	}

	// return all remaining movies
	json.NewEncoder(w).Encode(movies)
}

func main() {

	// create new router
	r := mux.NewRouter()

	// append some sample movies to our movies struct
	movies = append(movies, Movie{
		ID:    "1",
		ISBN:  "123456",
		Title: "Movie1",
		Director: &Director{
			FirstName: "John",
			LastName:  "Doe"}})
	movies = append(movies, Movie{
		ID:    "2",
		ISBN:  "45678",
		Title: "Movie 2",
		Director: &Director{
			FirstName: "Steve",
			LastName:  "Smith"}})

	// handle all API routes
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// start server
	fmt.Printf("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
