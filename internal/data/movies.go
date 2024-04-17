package data

import (
	"time"

	"delsanchez.gl/internal/validator"
)

type Movie struct {
	ID        int64     // Unique integer ID for each movie
	CreatedAt time.Time // Timestamp from when the movie is added to the database
	Title     string    // Movie title
	Year      int32     // Movie release year
	// Using the Runtime type instead of int32.
	Runtime Runtime  // Movie runtime/duration (in minutes)
	Genres  []string // Slices of genres for the movie (romance, comedy, etc)
	Version int32    // The version number starts at 1, and will be incremented each time the movie information is updated
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater that 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
