// Package service provides business logic services for the Mockbuster API.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
)

// commentServiceImpl implements the CommentService interface.
type commentServiceImpl struct {
	commentRepo repository.CommentRepositoryInterface
	filmRepo    repository.FilmRepositoryInterface
}

// NewCommentService creates a new comment service with the given repositories.
// This follows the Constructor Injection pattern from the article.
func NewCommentService(
	commentRepo repository.CommentRepositoryInterface,
	filmRepo repository.FilmRepositoryInterface,
) CommentService {
	return &commentServiceImpl{
		commentRepo: commentRepo,
		filmRepo:    filmRepo,
	}
}

// AddComment adds a new comment to a film.
func (s *commentServiceImpl) AddComment(
	_ context.Context,
	filmID int,
	commentReq models.CommentRequest,
) (*models.Comment, error) {
	if filmID <= 0 {
		slog.Warn("Invalid film ID provided", "filmID", filmID)
		return nil, errors.New("invalid film ID")
	}

	if err := s.validateComment(commentReq); err != nil {
		slog.Warn("Invalid comment provided", "comment", commentReq, "error", err)
		return nil, err
	}

	if _, err := s.filmRepo.GetFilmByID(filmID); err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			slog.Warn("Cannot add comment to non-existent film", "filmID", filmID)
			return nil, err
		}
		slog.Error("Failed to verify film exists", "filmID", filmID, "error", err)
		return nil, err
	}

	comment, err := s.commentRepo.AddComment(filmID, commentReq)
	if err != nil {
		slog.Error("Failed to add comment to repository", "filmID", filmID, "error", err)
		return nil, err
	}

	slog.Info("Successfully added comment", "filmID", filmID, "commentID", comment.ID)
	return comment, nil
}

// GetCommentsByFilmID retrieves all comments for a specific film.
func (s *commentServiceImpl) GetCommentsByFilmID(_ context.Context, filmID int) ([]models.Comment, error) {
	if filmID <= 0 {
		slog.Warn("Invalid film ID provided", "filmID", filmID)
		return nil, errors.New("invalid film ID")
	}

	if _, err := s.filmRepo.GetFilmByID(filmID); err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			slog.Warn("Cannot get comments for non-existent film", "filmID", filmID)
			return nil, err
		}
		slog.Error("Failed to verify film exists", "filmID", filmID, "error", err)
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByFilmID(filmID)
	if err != nil {
		slog.Error("Failed to retrieve comments from repository", "filmID", filmID, "error", err)
		return nil, err
	}

	slog.Info("Successfully retrieved comments", "filmID", filmID, "count", len(comments))
	return comments, nil
}

// validateComment validates the comment request.
func (s *commentServiceImpl) validateComment(commentReq models.CommentRequest) error {
	const (
		maxCustomerNameLength = 100
		maxCommentLength      = 1000
	)

	if commentReq.CustomerName == "" {
		return errors.New("customer name is required")
	}
	if len(commentReq.CustomerName) > maxCustomerNameLength {
		return errors.New("customer name too long (max 100 characters)")
	}

	if commentReq.Comment == "" {
		return errors.New("comment text is required")
	}
	if len(commentReq.Comment) > maxCommentLength {
		return errors.New("comment text too long (max 1000 characters)")
	}

	return nil
}
