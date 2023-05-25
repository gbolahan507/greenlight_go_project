package main

import (
	// "fmt"
	"database/sql"
	"errors"
	"fmt"
	model "greenlight_gbolahan/internal/data"
	"greenlight_gbolahan/internal/validator"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string        `json:"title"`
		Year    int32         `json:"year"`
		Runtime model.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie := &model.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if model.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movie/%d", movie.ID))

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrorRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrorRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	var input struct {
		Title   *string        `json:"title"`
		Year    *int32         `json:"year"`
		Runtime *model.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}
	v := validator.New()

	if model.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(*movie)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.errorEditingConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrorRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie deleted successfully"}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string
		Genres []string
		filter model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")

	input.Genres = app.readCSV(qs, "genres", []string{})

	input.filter.Page = app.readInt(qs, "page", 1, v)

	input.filter.PageSize = app.readInt(qs, "page_size", 20, v)

	input.filter.Sort = app.readString(qs, "sort", "id")

	// Add the supported sort values for this endpoint to the sort safelist.
	input.filter.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if model.ValidateFilters(v, input.filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, err := app.models.Movies.GetAllMovies(input.Title, input.Genres, input.filter)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
