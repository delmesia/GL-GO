package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters
	// will be store in the request context. The ParamsFromContext function can
	// can retrieve a slice containing the parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())
	// ByName will get the value of the given parameter (in this case, "ID") from the slice.
	// In this project, all movies will have a unique positive integer ID, but the value
	// returned by ByName() function is always a string, so we try to convert it to
	// a base 10 integer (with a bit of size 64). If the parameter couldn't be converted
	// or is less than 1, the ID is invalid and will use http.NotFound to return a 404 Not Found response.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}
	return id, nil
}

// writeJSON helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON,
// and a header map containing any additional HTTP headers we need to include in the response.
type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	//Encode the data to JSON. returning an error if there was one.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// Append a new line to make it nice in the terminal
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the header map
	// and add each header to the http.ResponseWriter header map. It's OK if the provided header
	// is nil. Go doesn't throw an error if you try to  to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type": "application/json" header, then write the status code and
	// json response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Decode the request body into the target destination
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		// If there is an error during the decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances, Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax error in the JSON. Check for this error using errors.Is() and return
		// a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contained badly-formed JSON")

		// Catch any *json.UnmarshalTypeErrors. These occurs when the JSON value is
		// the wrong type for the target destination. If the error relates to a specific
		// field, include that in the error message to make it easier for client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty.
		// will return an plain-english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A json.InvalidUnmarshalEror error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). Catch this and panic,
		// rather than returning an error to the handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error as-is.
		default:
			return err
		}
	}
	return nil
}
