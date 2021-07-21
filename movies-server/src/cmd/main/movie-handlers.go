package main

import (
	"github.com/DapperBlondie/movie-server/src/models"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	movieId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Println("Error in getting the movie id from req context : " + err.Error())
		return
	}

	movie := &models.Movie{
		ID:          movieId,
		Title:       "None",
		Description: "None of my concern",
		Year:        2021,
		ReleaseDate: time.Date(2021, 01, 01, 1, 20, 23, 0, time.UTC),
		Runtime:     100,
		Rating:      5,
		MPAARating:  "PG-13",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		log.Println("Error in movie-handlers : " + err.Error())
		return
	}
}

func (app *Application) getAllMovies(w http.ResponseWriter, r *http.Request) {

}
