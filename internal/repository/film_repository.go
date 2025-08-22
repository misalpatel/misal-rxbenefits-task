package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rxbenefits/go-hw/internal/database"
	"github.com/rxbenefits/go-hw/internal/models"
)

// FilmRepository handles database operations for films.
type FilmRepository struct {
	db *database.DB
}

// NewFilmRepository creates a new film repository.
func NewFilmRepository(db *database.DB) *FilmRepository {
	return &FilmRepository{db: db}
}

// GetFilms retrieves films with optional filters.
func (r *FilmRepository) GetFilms(filters models.FilmFilters) (*models.FilmListResponse, error) {
	r.normalizePagination(&filters)

	query, args := r.buildFilmsQuery(filters)
	films, err := r.executeFilmsQuery(query, args)
	if err != nil {
		return nil, err
	}

	total, err := r.getFilmsCount(filters)
	if err != nil {
		return nil, err
	}

	return &models.FilmListResponse{
		Films: films,
		Total: total,
		Page:  filters.Page,
		Limit: filters.Limit,
	}, nil
}

// normalizePagination sets default values for pagination parameters.
func (r *FilmRepository) normalizePagination(filters *models.FilmFilters) {
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
}

// buildFilmsQuery constructs the SQL query and arguments for fetching films.
func (r *FilmRepository) buildFilmsQuery(filters models.FilmFilters) (string, []interface{}) {
	query := `
		SELECT DISTINCT f.film_id, f.title, f.description, f.release_year, 
		       f.language_id, f.rental_duration, f.rental_rate, f.length, 
		       f.replacement_cost, f.rating, f.last_update, f.special_features
		FROM film f
		LEFT JOIN film_category fc ON f.film_id = fc.film_id
		LEFT JOIN category c ON fc.category_id = c.category_id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 0

	if filters.Title != "" {
		argCount++
		query += fmt.Sprintf(" AND f.title ILIKE $%d", argCount)
		args = append(args, "%"+filters.Title+"%")
	}

	if filters.Rating != "" {
		argCount++
		query += fmt.Sprintf(" AND f.rating = $%d", argCount)
		args = append(args, filters.Rating)
	}

	if filters.Category != "" {
		argCount++
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		args = append(args, "%"+filters.Category+"%")
	}

	offset := (filters.Page - 1) * filters.Limit
	argCount++
	query += fmt.Sprintf(" ORDER BY f.title LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filters.Limit, offset)

	return query, args
}

// executeFilmsQuery executes the query and scans the results into film objects.
func (r *FilmRepository) executeFilmsQuery(query string, args []interface{}) ([]models.Film, error) {
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying films: %w", err)
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		film, scanErr := r.scanFilm(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		films = append(films, film)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating films: %w", rowsErr)
	}

	return films, nil
}

// scanFilm scans a single film row and enriches it with categories and actors.
func (r *FilmRepository) scanFilm(rows *sql.Rows) (models.Film, error) {
	var film models.Film
	var specialFeatures sql.NullString

	scanErr := rows.Scan(
		&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear,
		&film.LanguageID, &film.RentalDuration, &film.RentalRate, &film.Length,
		&film.ReplacementCost, &film.Rating, &film.LastUpdate, &specialFeatures,
	)
	if scanErr != nil {
		return models.Film{}, fmt.Errorf("error scanning film: %w", scanErr)
	}

	if specialFeatures.Valid {
		features := strings.Trim(specialFeatures.String, "{}")
		if features != "" {
			film.SpecialFeatures = strings.Split(features, ",")
		}
	}

	categories, catErr := r.getFilmCategories(film.FilmID)
	if catErr != nil {
		return models.Film{}, catErr
	}
	film.Categories = categories

	actors, actorErr := r.getFilmActors(film.FilmID)
	if actorErr != nil {
		return models.Film{}, actorErr
	}
	film.Actors = actors

	return film, nil
}

// getFilmsCount gets the total count of films matching the filters.
func (r *FilmRepository) getFilmsCount(filters models.FilmFilters) (int, error) {
	countQuery := `
		SELECT COUNT(DISTINCT f.film_id)
		FROM film f
		LEFT JOIN film_category fc ON f.film_id = fc.film_id
		LEFT JOIN category c ON fc.category_id = c.category_id
		WHERE 1=1
	`

	countArgs := []interface{}{}
	argCount := 0

	if filters.Title != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND f.title ILIKE $%d", argCount)
		countArgs = append(countArgs, "%"+filters.Title+"%")
	}

	if filters.Rating != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND f.rating = $%d", argCount)
		countArgs = append(countArgs, filters.Rating)
	}

	if filters.Category != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		countArgs = append(countArgs, "%"+filters.Category+"%")
	}

	var total int
	err := r.db.QueryRowContext(context.Background(), countQuery, countArgs...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error counting films: %w", err)
	}

	return total, nil
}

// GetFilmByID retrieves a single film by ID.
func (r *FilmRepository) GetFilmByID(filmID int) (*models.Film, error) {
	query := `
		SELECT film_id, title, description, release_year, language_id, 
		       rental_duration, rental_rate, length, replacement_cost, 
		       rating, last_update, special_features
		FROM film 
		WHERE film_id = $1
	`

	var film models.Film
	var specialFeatures sql.NullString

	err := r.db.QueryRowContext(context.Background(), query, filmID).Scan(
		&film.FilmID, &film.Title, &film.Description, &film.ReleaseYear,
		&film.LanguageID, &film.RentalDuration, &film.RentalRate, &film.Length,
		&film.ReplacementCost, &film.Rating, &film.LastUpdate, &specialFeatures,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFilmNotFound
		}
		return nil, fmt.Errorf("error querying film: %w", err)
	}

	if specialFeatures.Valid {
		features := strings.Trim(specialFeatures.String, "{}")
		if features != "" {
			film.SpecialFeatures = strings.Split(features, ",")
		}
	}

	categories, err := r.getFilmCategories(filmID)
	if err != nil {
		return nil, err
	}
	film.Categories = categories

	actors, err := r.getFilmActors(filmID)
	if err != nil {
		return nil, err
	}
	film.Actors = actors

	return &film, nil
}

// getFilmCategories retrieves categories for a film.
func (r *FilmRepository) getFilmCategories(filmID int) ([]string, error) {
	query := `
		SELECT c.name 
		FROM category c
		JOIN film_category fc ON c.category_id = fc.category_id
		WHERE fc.film_id = $1
		ORDER BY c.name
	`

	rows, err := r.db.QueryContext(context.Background(), query, filmID)
	if err != nil {
		return nil, fmt.Errorf("error querying film categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		scanErr := rows.Scan(&category)
		if scanErr != nil {
			return nil, fmt.Errorf("error scanning category: %w", scanErr)
		}
		categories = append(categories, category)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating film categories: %w", rowsErr)
	}

	return categories, nil
}

// getFilmActors retrieves actors for a film.
func (r *FilmRepository) getFilmActors(filmID int) ([]string, error) {
	query := `
		SELECT a.first_name || ' ' || a.last_name as actor_name
		FROM actor a
		JOIN film_actor fa ON a.actor_id = fa.actor_id
		WHERE fa.film_id = $1
		ORDER BY a.last_name, a.first_name
	`

	rows, err := r.db.QueryContext(context.Background(), query, filmID)
	if err != nil {
		return nil, fmt.Errorf("error querying film actors: %w", err)
	}
	defer rows.Close()

	var actors []string
	for rows.Next() {
		var actor string
		scanErr := rows.Scan(&actor)
		if scanErr != nil {
			return nil, fmt.Errorf("error scanning actor: %w", scanErr)
		}
		actors = append(actors, actor)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating film actors: %w", rowsErr)
	}

	return actors, nil
}

// GetCategories retrieves all categories.
func (r *FilmRepository) GetCategories() ([]models.Category, error) {
	query := `SELECT category_id, name FROM category ORDER BY name`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		scanErr := rows.Scan(&category.CategoryID, &category.Name)
		if scanErr != nil {
			return nil, fmt.Errorf("error scanning category: %w", scanErr)
		}
		categories = append(categories, category)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("error iterating categories: %w", rowsErr)
	}

	return categories, nil
}
