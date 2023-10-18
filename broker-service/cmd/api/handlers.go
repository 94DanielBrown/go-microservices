package main

import (
	"fmt"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	fmt.Printf("Broker has been hit\n")

	_ = app.writeJSON(w, http.StatusOK, payload)
}
