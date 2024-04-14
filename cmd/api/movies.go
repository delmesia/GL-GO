package main

import (
	"fmt"
	"net/http"
	"time"

	"delsanchez.gl/internal/data"
)

// for "POST /v1/movies" endpoint
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP response body.
	// Note: The field name and types of the anonymous struct are a subset of the Movie struct
	// that I created. This struct will be the *target decode destination*
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	/************* version 1
	// Initialize a new json decoder instance which reads from the request body,
	and then use the Decode() method to decode the body contents into the input struct.
	When I call Decode(), I've used the address-of operator to the input struct
	as the target decode destination.

	err := json.NewDecoder(r.Body).Decode(&input)
	*****************************/
	// Using readJSON() helper to decode the request body into the input struct.
	// If this returns an error, we send the client an error message along with
	// 400 Bad Request Status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

// for "GET /v1/movies/:id" endpoint. For now, we'll retrieve the interpolated id
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Create a new instance of the Movie struct, containing the ID extracted from the URL
	// and some dummy data. No year field yet.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	// Encode the struct to JSON and send it as HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
