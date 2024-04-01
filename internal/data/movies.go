package data

import "time"

type Movie struct {
	ID        int64     // Unique integer ID for each movie
	CreatedAt time.Time // Timestamp from when the movie is added to the database
	Title     string    // Movie title
	Year      int32     // Movie release year
	Runtime   int32     // Movie runtime/duration (in minutes)
	Genres    []string  // Slices of genres for the movie (romance, comedy, etc)
	Version   int32     // The version number starts at 1, and will be incremented each time the movie information is updated
}
