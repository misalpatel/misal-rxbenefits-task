package repository

import "errors"

// ErrFilmNotFound is returned when a film is not found in the database.
var ErrFilmNotFound = errors.New("film not found")
