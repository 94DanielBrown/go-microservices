package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	fmt.Println(requestPayload)

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, errors.New("Can't find email"), http.StatusBadRequest)
		return
	}
	fmt.Println(user)

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		fmt.Println("Wrong Password")
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
