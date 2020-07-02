package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Movie struct
type Movie struct {
	ID          string `json:"id"`
	Name        string `json:"name" form:"name"`
	ImdbRating  string `json:"imdbrating form:"imdbrating"`
	ReleaseDate string `json:"releasedate" form:"releasedate"`
}

// Response struct
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var movie Movie
var movies []Movie

//Connect DB
func connect() *sql.DB {
	db, err := sql.Open("mysql", "u4560401_simple:Aa3einga101!!@tcp(tempatneduh.com:3306)/u4560401_simple")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

//Get ALL Movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	rows, err := db.Query("Select id, name, imdb_rating, release_date from movie")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Name, &movie.ImdbRating, &movie.ReleaseDate); err != nil {
			log.Fatal(err.Error())

		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movie)
		}
	}

}

//Get SINGLE Movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	params := mux.Vars(r)
	inputMovieID := params["id"]

	rows := db.QueryRow("Select id, name, imdb_rating, release_date from movie where id=?", inputMovieID)

	if err := rows.Scan(&movie.ID, &movie.Name, &movie.ImdbRating, &movie.ReleaseDate); err != nil {
		log.Fatal(err.Error())

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)

}

//Update SINGLE Movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	var response Response
	db := connect()
	defer db.Close()
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	inputMovieID := params["id"]

	name := r.FormValue("name")
	imdbRating := r.FormValue("imdbrating")
	releaseDate := r.FormValue("releasedate")

	_, err = db.Exec("UPDATE movie SET name=?, imdb_rating=?, release_date=? where id=?", name, imdbRating, releaseDate, inputMovieID)
	if err != nil {
		log.Print(err)
	}

	// if err := rows.Scan(&movie.ID, &movie.Name, &movie.ImdbRating, &movie.ReleaseDate); err != nil {
	// 	log.Fatal(err.Error())

	// }

	response.Status = 1
	response.Message = "Updated"
	log.Print("Data Updated")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//DELETE SINGLE Movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	var response Response
	db := connect()
	defer db.Close()
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}
	id := r.FormValue("id")
	// params := mux.Vars(r)
	// inputMovieID := params["id"]

	_, err = db.Exec("DELETE from movie where id=?", id)
	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Deleted"
	log.Print("Data Deleted")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func main() {

	//intial router
	router := mux.NewRouter()

	//API endpoints
	router.HandleFunc("/v1/movie", getMovies).Methods("GET")
	router.HandleFunc("/v1/movie/{id}", getMovie).Methods("GET")
	router.HandleFunc("/v1/movie/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/v1/movie/{id}", deleteMovie).Methods("DELETE")
	fmt.Println("Running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
