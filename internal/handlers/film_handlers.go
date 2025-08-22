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
func (h *FilmHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.filmService.GetCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// AddComment handles POST /films/{id}/comments.
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
func WelcomeHandler(w http.ResponseWriter, _ *http.Request) {
	response := models.WelcomeResponse{Message: "Welcome to Mockbuster Movie API!"}
	respondWithJSON(w, http.StatusOK, response)
}

// APIInfoHandler handles GET /api/v1.
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
