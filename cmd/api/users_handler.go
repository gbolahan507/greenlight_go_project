package main

import (
	"errors"
	model "greenlight_gbolahan/internal/data"
	"greenlight_gbolahan/internal/validator"
	"net/http"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrorDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email.

	// 	When this code is executed now, a new ‘background’ goroutine will be launched for sending
	// the welcome email. The code in this background goroutine will be executed concurrently
	// with the subsequent code in our registerUserHandler, which means we are no longer
	// waiting for the email to be sent before we return a JSON response to the client. Most likely,
	// the background goroutine will still be executing its code long after the
	// registerUserHandler has returned.

	go func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}

	}()
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	// This status code indicates that the request has been accepted for processing, but
	// the processing has not been completed because a bockground goroutine still needs to run 
	// for welcomming mail
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
