package main

import (
	"encoding/json"
	"github.com/DapperBlondie/movie-server/src/models"
	"github.com/julienschmidt/httprouter"
	zerolog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"strconv"
	"time"
)

// MoviePayload use for getting the payload data from post request
type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

// JsonResp use for give the response to client
type JsonResp struct {
	OK bool `json:"ok"`
}

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

// getAllGenres use for get all genres that available in the database
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

// getAllMoviesByGenre use for get all movies by specifying the
func (app *Application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		//app.errorJSON(w, err)
		return
	}

	movies, err := app.models.DB.AllByGenre(genreID)
	if err != nil {
		//app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		//app.errorJSON(w, err)
		return
	}
}

// editMovie use for inserting an movie into database
func (app *Application) editMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		zerolog.Error().Msg(err.Error() + " error occurred in decoding the body of request.")
		return
	}

	var movie models.Movie

	movie.ID, err = strconv.Atoi(payload.ID)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, _ = time.Parse("2006-07-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.ID == 0 {
		err = app.models.DB.InsertMovie(&movie)

		if err != nil {
			zerolog.Error().Msg(err.Error() + " occurred in editMovie handler")
			return
		}
	} else {
		err = app.models.DB.UpdateMovie(&movie)
		if err != nil {
			zerolog.Error().Msg(err.Error() + " occurred in editMovie handler")
			return
		}
	}

	ok := &JsonResp{OK: true}
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		zerolog.Error().Msg(err.Error() + " occurred in writing the response to client")
		return
	}
}

func (app *Application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	err = app.models.DB.DeleteMovie(id)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	ok := &JsonResp{OK: true}
	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}
}

func (app *Application) searchMovie(w http.ResponseWriter, r *http.Request) {

}
