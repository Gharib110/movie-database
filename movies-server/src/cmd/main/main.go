package main

import (
	"context"
	"database/sql"
	"github.com/DapperBlondie/movie-server/src/models"
	_ "github.com/lib/pq"
	zerolog "github.com/rs/zerolog/log"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     int
	HostName string
	db       struct {
		DSN string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type Application struct {
	ConfigApp *Config
	Logger    *log.Logger
	models    *models.Models
}

func main() {
	run()
	return
}

func run() {
	config := &Config{
		Port:     8080,
		HostName: "localhost",
		db:       struct{ DSN string }{DSN: "postgres://postgre:alireza1380##@localhost:5720/my_movies?sslmode=disable"},
	}

	db, err := openDB(config)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			zerolog.Fatal().Msg("Error in closing the db object : " + err.Error())
			return
		}
	}(db)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &Application{
		ConfigApp: config,
		Logger:    logger,
		models:    models.NewModel(db),
	}

	srv := &http.Server{
		Addr:              app.ConfigApp.HostName + ":" + strconv.Itoa(app.ConfigApp.Port),
		Handler:           app.routes(),
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute,
	}

	log.Println("App is listening on localhost:8080 ...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Error in serving the application : " + err.Error())
		return
	}
}

func openDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.DSN)
	if err != nil {
		zerolog.Fatal().Msg("Error occurred in openDB : " + err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		zerolog.Fatal().Msg("Error in pinging the db in openDB : " + err.Error())
		return nil, err
	}

	return db, nil
}
