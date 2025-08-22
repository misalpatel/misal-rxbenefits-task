// Package service provides business logic services for the Mockbuster API.
package service

import (
	"context"

	"github.com/rxbenefits/go-hw/internal/models"
)

// FilmService defines the interface for film-related business operations.
type FilmService interface {
	// GetFilms retrieves films with optional filtering and pagination.
	GetFilms(ctx context.Context, filters models.FilmFilters) (*models.FilmListResponse, error)

	// GetFilmByID retrieves a specific film by its ID.
	GetFilmByID(ctx context.Context, filmID int) (*models.Film, error)

	// GetCategories retrieves all available film categories.
	GetCategories(ctx context.Context) ([]models.Category, error)
}

// CommentService defines the interface for comment-related business operations.
type CommentService interface {
	// AddComment adds a new comment to a film.
	AddComment(ctx context.Context, filmID int, commentReq models.CommentRequest) (*models.Comment, error)

	// GetCommentsByFilmID retrieves all comments for a specific film.
	GetCommentsByFilmID(ctx context.Context, filmID int) ([]models.Comment, error)
}
