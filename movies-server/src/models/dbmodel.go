package models

import (
	"context"
	"database/sql"
	"fmt"
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

// GenresAll use for getting all genres in the database
func (d *DBModel) GenresAll() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	query := `SELECT id,genre_name,created_at,updated_at FROM genres ORDER BY genre_name`

	rows, err := d.DB.QueryContext(ctx, query)
	if err != nil {
		zerolog.Error().Msg(err.Error() + " occurred in getting all of genres")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zerolog.Error().Msg(err.Error() + " occurred in closing rows")
			return
		}
	}(rows)

	var genres []*Genre
	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			zerolog.Error().Msg(err.Error() + " occurred in scanning values")
			return nil, err
		}
		genres = append(genres, &g)
	}

	return genres, err
}

// AllByGenre get all movies that available in the database
func (d *DBModel) AllByGenre(genre ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, rating, runtime, mpaa_rating,
				created_at, updated_at from movies  %s order by title`, where)

	rows, err := d.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zerolog.Error().Msg(err.Error() + " in closing the Rows")
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
			return nil, err
		}

		// get genres, if any
		genreQuery := `select
			mg.id, mg.movie_id, mg.genre_id, g.genre_name
		from
			movies_genres mg
			left join genres g on (g.id = mg.genre_id)
		where
			mg.movie_id = $1
		`

		genreRows, _ := d.DB.QueryContext(ctx, genreQuery, movie.ID)

		genres := make(map[int]string)
		for genreRows.Next() {
			var mg MovieGenre
			err := genreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.ID] = mg.Genre.GenreName
		}
		err = genreRows.Close()
		if err != nil {
			return nil, err
		}

		movie.MovieGenre = genres
		movies = append(movies, &movie)

	}
	return movies, nil
}

// InsertMovie use for inserting a movie information into the database
func (d *DBModel) InsertMovie(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	stmt := `INSERT INTO movies (title,description,year,release_date,runtime,rating,mpaa_rating,created_at,updated_at) VALUES 
($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := d.DB.ExecContext(ctx, stmt,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		zerolog.Error().Msg(err.Error() + " occurred in InsertMovie function")
		return err
	}

	return nil
}
