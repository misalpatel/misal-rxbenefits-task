// Package models provides data structures for the Mockbuster API.
package models

import (
	"time"
)

// Film represents a movie in the database.
type Film struct {
	FilmID          int       `json:"film_id"                    db:"film_id"`
	Title           string    `json:"title"                      db:"title"            validate:"required"`
	Description     *string   `json:"description,omitempty"      db:"description"`
	ReleaseYear     *int      `json:"release_year,omitempty"     db:"release_year"`
	LanguageID      int       `json:"language_id"                db:"language_id"`
	RentalDuration  int       `json:"rental_duration"            db:"rental_duration"`
	RentalRate      float64   `json:"rental_rate"                db:"rental_rate"`
	Length          *int      `json:"length,omitempty"           db:"length"`
	ReplacementCost float64   `json:"replacement_cost"           db:"replacement_cost"`
	Rating          string    `json:"rating"                     db:"rating"`
	LastUpdate      time.Time `json:"last_update"                db:"last_update"`
	SpecialFeatures []string  `json:"special_features,omitempty" db:"special_features"`
	Categories      []string  `json:"categories,omitempty"`
	Actors          []string  `json:"actors,omitempty"`
}

// FilmListResponse represents the response for listing films.
type FilmListResponse struct {
	Films []Film `json:"films"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// FilmFilters represents filters for film search.
type FilmFilters struct {
	Title    string `json:"title,omitempty"`
	Rating   string `json:"rating,omitempty"`
	Category string `json:"category,omitempty"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// Comment represents a customer comment on a film.
type Comment struct {
	ID           int       `json:"id"            db:"id"`
	FilmID       int       `json:"film_id"       db:"film_id"       validate:"required"`
	CustomerName string    `json:"customer_name" db:"customer_name" validate:"required"`
	Comment      string    `json:"comment"       db:"comment"       validate:"required"`
	CreatedAt    time.Time `json:"created_at"    db:"created_at"`
}

// CommentRequest represents the request to add a comment.
type CommentRequest struct {
	CustomerName string `json:"customer_name" validate:"required"`
	Comment      string `json:"comment"       validate:"required"`
}

// Category represents a film category.
type Category struct {
	CategoryID int    `json:"category_id" db:"category_id"`
	Name       string `json:"name"        db:"name"`
}

// Actor represents a film actor.
type Actor struct {
	ActorID   int    `json:"actor_id"   db:"actor_id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name"  db:"last_name"`
}

// WelcomeResponse represents the welcome message response.
type WelcomeResponse struct {
	Message string `json:"message" example:"Welcome to Mockbuster Movie API!"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"             example:"Failed to retrieve films"`
	Details string `json:"details,omitempty" example:"database connection failed"`
}

// APIInfoResponse represents the API information response.
type APIInfoResponse struct {
	Name          string   `json:"name"          example:"Mockbuster Movie API"`
	Version       string   `json:"version"       example:"1.0"`
	Description   string   `json:"description"   example:"A RESTful API for the Mockbuster DVD rental business"`
	Endpoints     []string `json:"endpoints"     example:"GET /api/v1/films"`
	Documentation string   `json:"documentation" example:"http://localhost:8080/swagger/"`
}
