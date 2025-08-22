// Package repository provides data access layer for the Mockbuster API.
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rxbenefits/go-hw/internal/database"
	"github.com/rxbenefits/go-hw/internal/models"
)

// CommentRepository handles database operations for comments.
type CommentRepository struct {
	db *database.DB
}

// NewCommentRepository creates a new comment repository.
func NewCommentRepository(db *database.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// AddComment adds a new comment to a film.
func (r *CommentRepository) AddComment(filmID int, commentReq models.CommentRequest) (*models.Comment, error) {
	var filmExists bool
	err := r.db.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM film WHERE film_id = $1)", filmID).
		Scan(&filmExists)
	if err != nil {
		return nil, fmt.Errorf("error checking film existence: %w", err)
	}
	if !filmExists {
		return nil, ErrFilmNotFound
	}

	query := `
		INSERT INTO film_comments (film_id, customer_name, comment, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, film_id, customer_name, comment, created_at
	`

	var comment models.Comment
	now := time.Now()
	err = r.db.QueryRowContext(context.Background(), query, filmID, commentReq.CustomerName, commentReq.Comment, now).
		Scan(
			&comment.ID, &comment.FilmID, &comment.CustomerName, &comment.Comment, &comment.CreatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("error inserting comment: %w", err)
	}

	return &comment, nil
}

// GetCommentsByFilmID retrieves all comments for a specific film.
func (r *CommentRepository) GetCommentsByFilmID(filmID int) ([]models.Comment, error) {
	var filmExists bool
	err := r.db.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM film WHERE film_id = $1)", filmID).
		Scan(&filmExists)
	if err != nil {
		return nil, fmt.Errorf("error checking film existence: %w", err)
	}
	if !filmExists {
		return nil, ErrFilmNotFound
	}

	query := `
		SELECT id, film_id, customer_name, comment, created_at
		FROM film_comments
		WHERE film_id = $1
		ORDER BY created_at DESC
	`

	rows, queryErr := r.db.QueryContext(context.Background(), query, filmID)
	if queryErr != nil {
		return nil, fmt.Errorf("error querying comments: %w", queryErr)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		scanErr := rows.Scan(&comment.ID, &comment.FilmID, &comment.CustomerName, &comment.Comment, &comment.CreatedAt)
		if scanErr != nil {
			return nil, fmt.Errorf("error scanning comment: %w", scanErr)
		}
		comments = append(comments, comment)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating comments: %w", rowsErr)
	}

	return comments, nil
}
