package main

import (
	"github.com/julienschmidt/httprouter"
	zerolog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"strconv"
)

// getOneMovie use for getting one movie from database by its ID
func (app *Application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	movieId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Println("Error in getting the movie id from req context : " + err.Error())
		return
	}

	movie, err := app.models.DB.GetOne(movieId)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		zerolog.Error().Msg("Error in movie-handlers : " + err.Error())
		return
	}
}

// getAllMovies use for getting all movies that we can find in the database
func (app *Application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAll()
	if err != nil {
		zerolog.Error().Msg("Error in getting all movies : " + err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		zerolog.Error().Msg("Error in writing the movies to response writer : " + err.Error())
		return
	}
}

func (app *Application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GenresAll()
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		zerolog.Error().Msg(err.Error() + " occurred in writing json to response writer")
		return
	}
	return
}

func (app *Application) deleteMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) updateMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) insertMovie(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) searchMovie(w http.ResponseWriter, r *http.Request) {

}
