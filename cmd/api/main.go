package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// A string containing the application version number. Later,
// this number will be generated automatically at build time.
const version = "1.0.0"

// A config struct that will hold all the configuration settings of the application.
// For now, the configuration setting will be the network port that we want the server
// to listen on, and then name of the current operating system environment.
type config struct {
	port int
	env  string
}

// This application struct will hold the dependencies for the HTTP handlers,
// helpers, and middlewares. For now, its only the config struct and logger,
// but this will grow overtime as the project matures.

type application struct {
	config config
	logger *log.Logger
}

func main() {

	// Declaring an instance of config struct
	var cfg config

	// Read the value of the port and env command-line flag into the config struct.
	// default is using 4000 and the enviroment "development" if no
	// correspanding flags are provided.

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes a message to stdout stream.
	// prefixed with the current date and time.

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Instance of application struct containing config struct and the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// declare a new servemux and add /v1/healthcheck route which dispatches requests
	// to the healthcheckhandler method.
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthCheckHandler)

	// Declare a HTTP server with sensible timeout settings, which listens on the port
	// provided in the config struct and use the servemux created above as the handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// This will start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
