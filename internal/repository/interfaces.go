package repository

import (
	"github.com/rxbenefits/go-hw/internal/models"
)

// FilmRepositoryInterface defines the interface for film-related database operations.
type FilmRepositoryInterface interface {
	// GetFilms retrieves films with optional filtering and pagination.
	GetFilms(filters models.FilmFilters) (*models.FilmListResponse, error)

	// GetFilmByID retrieves a specific film by its ID.
	GetFilmByID(filmID int) (*models.Film, error)

	// GetCategories retrieves all available film categories.
	GetCategories() ([]models.Category, error)
}

// CommentRepositoryInterface defines the interface for comment-related database operations.
type CommentRepositoryInterface interface {
	// AddComment adds a new comment to a film.
	AddComment(filmID int, commentReq models.CommentRequest) (*models.Comment, error)

	// GetCommentsByFilmID retrieves all comments for a specific film.
	GetCommentsByFilmID(filmID int) ([]models.Comment, error)
}
