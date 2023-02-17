package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := envelope{
		"status":     "available",
		"enviroment": app.config.env,
		"version":    version,
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"healthCheck": data}, nil)
	if err != nil {
		app.logger.Print(err)
		app.serverErrorResponse(w, r, err)
		return
	}
}
