package main

import (
	"fmt"
	"net/http"
	"time"

	"delsanchez.gl/internal/data"
	"delsanchez.gl/internal/validator"
)

// for "POST /v1/movies" endpoint
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP response body.
	// Note: The field name and types of the anonymous struct are a subset of the Movie struct
	// that I created. This struct will be the *target decode destination*
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	// Using readJSON() helper to decode the request body into the input struct.
	// If this returns an error, we send the client an error message along with
	// 400 Bad Request Status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new validator instance.
	v := validator.New()

	// Use the Check() method to execute the validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line I "check that the title is not equal to the empty string"
	// In the second line, I "check that the length of the title is less than or equal to 500 bytes long"
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain atleast 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Here, I'm using the Unique() helper method to check that all the values in the
	// input.Genres slices are unique
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	// Use the Valid() helper method to see if any of the checks failed. If they did,
	// then use the failedValidationResponse() helper to send a response to the client, passing
	// the v.Errors map.
	if v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
