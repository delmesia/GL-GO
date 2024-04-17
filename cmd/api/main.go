package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// A string containing the application version number. Later,
// this number will be generated automatically at build time.
const version = "1.0.0"

// A config struct that will hold all the configuration settings of the application.
// For now, the configuration setting will be the network port that we want the server
// to listen on, and then name of the current operating system environment.
// ===================================================================================
// Add a db struct field to hold the configuration settings for the database
// connection pool. For now this only holds the DSN, which we will read in from the
// command-line flag.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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
	// Read the DSN value from the db-dsn command-line flag into the config struct.
	// Default is development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:2031@localhost/greenlight", "PostgreSQL DSN")
	flag.Parse()

	// Initialize a new logger which writes a message to stdout stream.
	// prefixed with the current date and time.

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Call the openDB() helper function to create the connection pool
	// passing in the config struct,. If this returns an error, log it and
	// exit the application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// Defer a call to db.Close() so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// if there are no errors above, log a message to say that
	// the connection pool has been successfully established.
	logger.Printf("database connection pool has been established.")

	// Instance of application struct containing config struct and the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a HTTP server with sensible timeout settings, which listens on the port
	// provided in the config struct and use the servemux created above as the handler.
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		// Using the httprouter instance returned by app.routes() as the server handler
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// This will start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from
	// the config struct
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// A context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing
	// in the context created above as a parameter. If the connection couldn't
	//  be established successfully within 5 seconds, then this will return an error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
