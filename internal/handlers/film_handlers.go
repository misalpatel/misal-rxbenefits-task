// Package handlers provides HTTP request handlers for the Mockbuster API.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
	"github.com/rxbenefits/go-hw/internal/service"
)

// FilmHandler handles HTTP requests for films.
type FilmHandler struct {
	filmService    service.FilmService
	commentService service.CommentService
	validate       *validator.Validate
}

// NewFilmHandler creates a new film handler with the given services.
// This follows the Constructor Injection pattern from the article.
func NewFilmHandler(filmService service.FilmService, commentService service.CommentService) *FilmHandler {
	return &FilmHandler{
		filmService:    filmService,
		commentService: commentService,
		validate:       validator.New(),
	}
}

// GetFilms handles GET /films.
// @Summary Get films with filtering and pagination
// @Description Retrieve a list of films with optional filtering by title, rating, and category, and pagination support
// @Tags films
// @Accept json
// @Produce json
// @Param title query string false "Filter by film title (partial match)"
// @Param rating query string false "Filter by film rating (e.g., G, PG, PG-13, R, NC-17)"
// @Param category query string false "Filter by film category"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of films per page (default: 10)"
// @Success 200 {object} []models.Film "List of films"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/films [get]
func (h *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters.
	filters := models.FilmFilters{
		Title:    r.URL.Query().Get("title"),
		Rating:   r.URL.Query().Get("rating"),
		Category: r.URL.Query().Get("category"),
	}

	// Parse pagination parameters.
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		} else {
			filters.Page = 1
		}
	} else {
		filters.Page = 1
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		} else {
			filters.Limit = 10
		}
	} else {
		filters.Limit = 10
	}

	// Get films from service.
	films, err := h.filmService.GetFilms(r.Context(), filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve films", err)
		return
	}

	respondWithJSON(w, http.StatusOK, films)
}

// GetFilmByID handles GET /films/{id}.
// @Summary Get film by ID
// @Description Retrieve detailed information about a specific film by its ID
// @Tags films
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Success 200 {object} models.Film "Film details"
// @Failure 400 {object} models.ErrorResponse "Invalid film ID"
// @Failure 404 {object} models.ErrorResponse "Film not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/films/{id} [get]
func (h *FilmHandler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	film, err := h.filmService.GetFilmByID(r.Context(), filmID)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			respondWithError(w, http.StatusNotFound, "Film not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve film", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, film)
}

// GetCategories handles GET /categories.
// @Summary Get all film categories
// @Description Retrieve a list of all available film categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} string "List of categories"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/categories [get]
func (h *FilmHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.filmService.GetCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// AddComment handles POST /films/{id}/comments.
// @Summary Add a comment to a film
// @Description Add a new comment to a specific film
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Param comment body models.CommentRequest true "Comment details"
// @Success 201 {object} models.Comment "Comment created"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or validation failed"
// @Failure 404 {object} models.ErrorResponse "Film not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/films/{id}/comments [post]
func (h *FilmHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	var commentReq models.CommentRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&commentReq); decodeErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", decodeErr)
		return
	}

	// Validate the request.
	if validateErr := h.validate.Struct(commentReq); validateErr != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", validateErr)
		return
	}

	comment, err := h.commentService.AddComment(r.Context(), filmID, commentReq)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			respondWithError(w, http.StatusNotFound, "Film not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to add comment", err)
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, comment)
}

// GetComments handles GET /films/{id}/comments.
// @Summary Get comments for a film
// @Description Retrieve all comments for a specific film
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Film ID"
// @Success 200 {array} models.Comment "List of comments"
// @Failure 400 {object} models.ErrorResponse "Invalid film ID"
// @Failure 404 {object} models.ErrorResponse "Film not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/films/{id}/comments [get]
func (h *FilmHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	comments, err := h.commentService.GetCommentsByFilmID(r.Context(), filmID)
	if err != nil {
		if errors.Is(err, repository.ErrFilmNotFound) {
			respondWithError(w, http.StatusNotFound, "Film not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve comments", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, comments)
}

// WelcomeHandler handles GET /.
// @Summary Welcome message
// @Description Get a welcome message for the Mockbuster API
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} models.WelcomeResponse "Welcome message"
// @Router / [get]
func WelcomeHandler(w http.ResponseWriter, _ *http.Request) {
	response := models.WelcomeResponse{Message: "Welcome to Mockbuster Movie API!"}
	respondWithJSON(w, http.StatusOK, response)
}

// APIInfoHandler handles GET /api/v1.
// @Summary API information
// @Description Get information about the Mockbuster API including available endpoints
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} models.APIInfoResponse "API information"
// @Router /api/v1 [get]
func APIInfoHandler(w http.ResponseWriter, _ *http.Request) {
	response := models.APIInfoResponse{
		Name:        "Mockbuster Movie API",
		Version:     "1.0",
		Description: "A RESTful API for the Mockbuster DVD rental business",
		Endpoints: []string{
			"GET /api/v1/films - List films with filtering and pagination",
			"GET /api/v1/films/{id} - Get detailed film information",
			"GET /api/v1/categories - List all available categories",
			"POST /api/v1/films/{id}/comments - Add a comment to a film",
			"GET /api/v1/films/{id}/comments - Get comments for a film",
		},
		Documentation: "http://localhost:8080/swagger/",
	}
	respondWithJSON(w, http.StatusOK, response)
}

// Helper functions.
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, writeErr := w.Write(response); writeErr != nil {
		slog.Error("Failed to write response", "error", writeErr)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	errorResponse := models.ErrorResponse{
		Error:   message,
		Details: err.Error(),
	}
	respondWithJSON(w, code, errorResponse)
}
