package models

import (
	"context"
	"database/sql"
	zerolog "github.com/rs/zerolog/log"
	"time"
)

// DBModel for holding the sql.DB variable data
type DBModel struct {
	DB *sql.DB
}

// GetOne get an movie with its ID
func (d *DBModel) GetOne(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	movie := &Movie{}
	query := `SELECT id,title,description,year,release_date,runtime,rating,mpaa_rating,created_at,updated_at 
FROM movies WHERE id=$1`
	row := d.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		zerolog.Error().Msg("Error in scanning in GetOne : " + err.Error())
		return nil, err
	}

	query = `SELECT mg.id,mg.movie_id,mg.genre_id,g.genre_name
FROM 
movie_genres mg
LEFT JOIN genres g on (g.id=mg.genre_id)
WHERE mg.movie_id=$1`

	rows, err := d.DB.QueryContext(ctx, query, id)
	if err != nil {
		zerolog.Error().Msg("Error in getting movie's genres : " + err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zerolog.Error().Msg("Error in closing the rows : " + err.Error())
			return
		}
	}(rows)

	genres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			zerolog.Error().Msg("Error in scanning the rows : " + err.Error())
			return nil, err
		}
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres
	return movie, nil
}

// GetAll get all movies in the database
func (d *DBModel) GetAll() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT id,title,description,year,release_date,runtime,rating,mpaa_rating,created_at,updated_at 
FROM movies WHERE id=$1 ORDER BY title`

	rows, err := d.DB.QueryContext(ctx, query)
	if err != nil {
		zerolog.Error().Msg("Error in getting all movies : " + err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zerolog.Error().Msg("Error in closing the rows : " + err.Error())
			return
		}
	}(rows)

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			zerolog.Error().Msg("Error in scanning the data : " + err.Error())
			return nil, err
		}

		queryGenres := `SELECT mg.id,mg.movie_id,mg.genre_id,g.genre_name
FROM 
movie_genres mg
LEFT JOIN genres g on (g.id=mg.genre_id)
WHERE mg.movie_id=$1`

		rows, err := d.DB.QueryContext(ctx, queryGenres, movie.ID)
		if err != nil {
			zerolog.Error().Msg("Error in getting movie's genres : " + err.Error())
			return nil, err
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				zerolog.Error().Msg("Error in closing the rows : " + err.Error())
				return
			}
		}(rows)

		genres := make(map[int]string)
		for rows.Next() {
			var mg MovieGenre
			err := rows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				zerolog.Error().Msg("Error in scanning the rows : " + err.Error())
				return nil, err
			}
			genres[mg.ID] = mg.Genre.GenreName
		}
		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}
