package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     int
	HostName string
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type Application struct {
	ConfigApp *Config
	Logger    *log.Logger
}

func main() {
	run()
	return
}

func run() {
	config := &Config{
		Port:     8080,
		HostName: "localhost",
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &Application{
		ConfigApp: config,
		Logger:    logger,
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
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Error in serving the application : " + err.Error())
		return
	}
}
