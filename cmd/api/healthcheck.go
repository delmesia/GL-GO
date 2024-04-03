package main

import (
	"net/http"
)

/*
****************************************************************************************
// A handler which writes a plain-text response with information about
// the application status, operating environment, and version.
// Important: healthCheckHandler is implemented as a method on the application struct -
// an idiomatic way to make dependencies available to the handler w/o resorting to
// global variable or closures. any dependency that the healthcheckHandler needs can simply
// be included as a field in the application struct when we initialize it in main().

	func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
		// Create a fixed-format JSON response from a string. I'm using a raw string
		// literal (enclosed with backticks) so that I can include double-quote character
		// in the JSON without needing to escape them. We also use the %q verb to
		// wrap the interpolated values in double-quotes.
		js := `{"status": "available, "environment": %q, "version": %q}`
		js = fmt.Sprintf(js, app.config.env, version)

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(js))
	}

************************************************************************************
*/
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an envelope map containing the data for the response.
	// The way I constructed this means the environment and version data will now
	// be nested under a system_info key in the JSON response.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
