package main

import (
	"fmt"
	"net/http"
)

// A generic helper for logging an error message.
// TODO: upgrade to structured logging, and record additional information
// about the request including the HTTP method and url
func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error message to the client
// with a given status code.
// Note: I'm using "any" as the type for the "message" parameter rather than a string type,
// as this will give me more flexibility over the values that I can include in the parameter.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return
	// an error, then log it, and fallback to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method will be used when the application encounter an unexpected
// problem at runtime. It logs the detailed error message, then uses the errorResponse() helper to
// send a 500 Internal Server Error status code and JSON response (containing a generic error message) to
// the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)

}

// The notFoundResponse() method will be used to send a 404 Not Found Status code and
// JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowed() method will be used to send a 405 Method Not Allowed
// status code and JSON response tok the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter has the type map[string]string, which is exactly
// the same as the errors map contained in the Validator type.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
