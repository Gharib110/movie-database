package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (app *Application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := &AppStatus{
		Status:      "Available",
		Environment: app.ConfigApp.HostName + ":" + strconv.Itoa(app.ConfigApp.Port),
		Version:     "1.0.0",
	}

	out, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		log.Fatalln("Error in marshaling currentStatus : " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		w.Write([]byte("Error occurred : " + err.Error() + "\n"))
		log.Println(err.Error())
		return
	}
}
