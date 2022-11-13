package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayLoad JSONPayload
	_ = app.readJSON(w, r, &requestPayLoad)

	event := data.LogEntry{
		Name: requestPayLoad.Name,
		Data: requestPayLoad.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	res := jsonResponse{
		Error:   false,
		Message: "Log entry created",
	}
	app.writeJSON(w, http.StatusCreated, res)
}
