package main

import (
	"encoding/json"
	"fmt"
	"github.com/DapperBlondie/movie-server/src/models"
	"github.com/graphql-go/graphql"
	zerolog "github.com/rs/zerolog/log"
	"io"
	"log"
	"net/http"
	"strings"
)

var movies []*models.Movie

// graphql schema definition
var fields = graphql.Fields{
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get all movies",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Search movies by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var theList []*models.Movie
			search, ok := params.Args["titleContains"].(string)
			if ok {
				for _, currentMovie := range movies {
					if strings.Contains(currentMovie.Title, search) {
						log.Println("Found one")
						theList = append(theList, currentMovie)
					}
				}
			}
			return theList, nil
		},
	},
}

var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

func (app *Application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	movies, _ = app.models.DB.GetAll()

	q, err := io.ReadAll(r.Body)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	query := string(q)

	log.Println(query)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		http.Error(w, "Failed to create graphql schema", http.StatusInternalServerError)
		zerolog.Error().Msg(err.Error())
		log.Println(err)
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		http.Error(w, fmt.Sprintf("failed: %+v", resp.Errors), http.StatusInternalServerError)
		zerolog.Error().Msg(fmt.Sprintf("failed: %+v", resp.Errors))
		return
	}

	j, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

	return
}
