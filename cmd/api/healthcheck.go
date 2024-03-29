package main

import (
	"fmt"
	"net/http"
)

// A handler which writes a plain-text response with information about
// the application status, operating environment, and version.
// Important: healthCheckHandler is implemented as a method on the application struct -
// an idiomatic way to make dependencies available to the handler w/o resorting to
// global variable or closures. any dependency that the healthcheckHandler needs can simply
// be included as a field in the application struct when we initialize it in main().
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version %s", version)
}
