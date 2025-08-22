// Package service provides business logic services for the Mockbuster API.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
)

// filmServiceImpl implements the FilmService interface.
type filmServiceImpl struct {
	filmRepo repository.FilmRepositoryInterface
}

// NewFilmService creates a new film service with the given repository.
func NewFilmService(filmRepo repository.FilmRepositoryInterface) FilmService {
	return &filmServiceImpl{
		filmRepo: filmRepo,
	}
}

// GetFilms retrieves films with optional filtering and pagination.
func (s *filmServiceImpl) GetFilms(_ context.Context, filters models.FilmFilters) (*models.FilmListResponse, error) {
	if err := s.validateFilters(filters); err != nil {
		slog.Warn("Invalid filters provided", "filters", filters, "error", err)
		return nil, err
	}

	s.applyDefaultPagination(&filters)

	films, err := s.filmRepo.GetFilms(filters)
	if err != nil {
		slog.Error("Failed to retrieve films from repository", "filters", filters, "error", err)
		return nil, err
	}

	slog.Info("Successfully retrieved films", "count", len(films.Films), "total", films.Total)
	return films, nil
}

// GetFilmByID retrieves a specific film by its ID.
func (s *filmServiceImpl) GetFilmByID(_ context.Context, filmID int) (*models.Film, error) {
	if filmID <= 0 {
		slog.Warn("Invalid film ID provided", "filmID", filmID)
		return nil, errors.New("invalid film ID")
	}

	film, err := s.filmRepo.GetFilmByID(filmID)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			slog.Warn("Film not found", "filmID", filmID)
			return nil, err
		}
		slog.Error("Failed to retrieve film from repository", "filmID", filmID, "error", err)
		return nil, err
	}

	slog.Info("Successfully retrieved film", "filmID", filmID, "title", film.Title)
	return film, nil
}

// GetCategories retrieves all available film categories.
func (s *filmServiceImpl) GetCategories(_ context.Context) ([]models.Category, error) {
	categories, err := s.filmRepo.GetCategories()
	if err != nil {
		slog.Error("Failed to retrieve categories from repository", "error", err)
		return nil, err
	}

	slog.Info("Successfully retrieved categories", "count", len(categories))
	return categories, nil
}

// validateFilters validates the provided filters.
func (s *filmServiceImpl) validateFilters(filters models.FilmFilters) error {
	if filters.Page < 1 {
		return errors.New("page must be greater than 0")
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	if filters.Rating != "" {
		validRatings := map[string]bool{
			"G": true, "PG": true, "PG-13": true, "R": true, "NC-17": true,
		}
		if !validRatings[filters.Rating] {
			return errors.New("invalid rating provided")
		}
	}

	return nil
}

// applyDefaultPagination applies default pagination values if not provided.
func (s *filmServiceImpl) applyDefaultPagination(filters *models.FilmFilters) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
}
