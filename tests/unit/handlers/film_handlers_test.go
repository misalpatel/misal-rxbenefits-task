package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rxbenefits/go-hw/internal/handlers"
	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
)

type MockFilmService struct {
	mock.Mock
}

func (m *MockFilmService) GetFilms(ctx context.Context, filters models.FilmFilters) (*models.FilmListResponse, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FilmListResponse), args.Error(1)
}

func (m *MockFilmService) GetFilmByID(ctx context.Context, filmID int) (*models.Film, error) {
	args := m.Called(ctx, filmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Film), args.Error(1)
}

func (m *MockFilmService) GetCategories(ctx context.Context) ([]models.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Category), args.Error(1)
}

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) AddComment(
	ctx context.Context,
	filmID int,
	commentReq models.CommentRequest,
) (*models.Comment, error) {
	args := m.Called(ctx, filmID, commentReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsByFilmID(ctx context.Context, filmID int) ([]models.Comment, error) {
	args := m.Called(ctx, filmID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func TestFilmHandler_GetFilms(t *testing.T) {
	tests := []struct {
		name               string
		queryParams        string
		mockResponse       *models.FilmListResponse
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:        "successful retrieval",
			queryParams: "?title=test&page=1&limit=10",
			mockResponse: &models.FilmListResponse{
				Films: []models.Film{
					{FilmID: 1, Title: "Test Film", Rating: "PG"},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &models.FilmListResponse{
				Films: []models.Film{
					{FilmID: 1, Title: "Test Film", Rating: "PG"},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			},
		},
		{
			name:               "service error",
			queryParams:        "?page=1&limit=10",
			mockError:          errors.New("database error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &models.ErrorResponse{
				Error:   "Failed to retrieve films",
				Details: "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmService := new(MockFilmService)
			mockCommentService := new(MockCommentService)
			handler := handlers.NewFilmHandler(mockFilmService, mockCommentService)

			// Setup mock expectations
			mockFilmService.On("GetFilms", mock.Anything, mock.AnythingOfType("models.FilmFilters")).
				Return(tt.mockResponse, tt.mockError)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/films"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Execute handler
			handler.GetFilms(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Convert expected response to JSON and back for comparison
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			var expectedInterface interface{}
			json.Unmarshal(expectedJSON, &expectedInterface)

			assert.Equal(t, expectedInterface, response)
			mockFilmService.AssertExpectations(t)
		})
	}
}

func TestFilmHandler_GetFilmByID(t *testing.T) {
	tests := []struct {
		name               string
		filmID             string
		mockResponse       *models.Film
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:   "successful retrieval",
			filmID: "1",
			mockResponse: &models.Film{
				FilmID: 1,
				Title:  "Test Film",
				Rating: "PG",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &models.Film{
				FilmID: 1,
				Title:  "Test Film",
				Rating: "PG",
			},
		},
		{
			name:               "film not found",
			filmID:             "999",
			mockError:          repository.ErrFilmNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &models.ErrorResponse{
				Error:   "Film not found",
				Details: "film not found",
			},
		},
		{
			name:               "invalid film ID",
			filmID:             "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &models.ErrorResponse{
				Error:   "Invalid film ID",
				Details: "strconv.Atoi: parsing \"invalid\": invalid syntax",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmService := new(MockFilmService)
			mockCommentService := new(MockCommentService)
			handler := handlers.NewFilmHandler(mockFilmService, mockCommentService)

			// Setup mock expectations only for valid film IDs
			if tt.filmID != "invalid" {
				filmID := 1
				if tt.filmID == "999" {
					filmID = 999
				}
				mockFilmService.On("GetFilmByID", mock.Anything, filmID).Return(tt.mockResponse, tt.mockError)
			}

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodGet, "/films/"+tt.filmID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.filmID})
			w := httptest.NewRecorder()

			// Execute handler
			handler.GetFilmByID(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Convert expected response to JSON and back for comparison
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			var expectedInterface interface{}
			json.Unmarshal(expectedJSON, &expectedInterface)

			assert.Equal(t, expectedInterface, response)
			mockFilmService.AssertExpectations(t)
		})
	}
}

func TestFilmHandler_AddComment(t *testing.T) {
	tests := []struct {
		name               string
		filmID             string
		requestBody        models.CommentRequest
		mockResponse       *models.Comment
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:   "successful comment addition",
			filmID: "1",
			requestBody: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			mockResponse: &models.Comment{
				ID:           1,
				FilmID:       1,
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: &models.Comment{
				ID:           1,
				FilmID:       1,
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
		},
		{
			name:   "film not found",
			filmID: "999",
			requestBody: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			mockError:          repository.ErrFilmNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &models.ErrorResponse{
				Error:   "Film not found",
				Details: "film not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmService := new(MockFilmService)
			mockCommentService := new(MockCommentService)
			handler := handlers.NewFilmHandler(mockFilmService, mockCommentService)

			// Setup mock expectations
			filmID := 1
			if tt.filmID == "999" {
				filmID = 999
			}
			mockCommentService.On("AddComment", mock.Anything, filmID, tt.requestBody).
				Return(tt.mockResponse, tt.mockError)

			// Create request body
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/films/"+tt.filmID+"/comments", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.filmID})
			w := httptest.NewRecorder()

			// Execute handler
			handler.AddComment(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Convert expected response to JSON and back for comparison
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			var expectedInterface interface{}
			json.Unmarshal(expectedJSON, &expectedInterface)

			assert.Equal(t, expectedInterface, response)
			mockCommentService.AssertExpectations(t)
		})
	}
}

func TestWelcomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handlers.WelcomeHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.WelcomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Welcome to Mockbuster Movie API!", response.Message)
}

func TestFilmHandler_GetCategories(t *testing.T) {
	tests := []struct {
		name               string
		mockResponse       []models.Category
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "successful retrieval",
			mockResponse: []models.Category{
				{CategoryID: 1, Name: "Action"},
				{CategoryID: 2, Name: "Comedy"},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: []models.Category{
				{CategoryID: 1, Name: "Action"},
				{CategoryID: 2, Name: "Comedy"},
			},
		},
		{
			name:               "service error",
			mockError:          errors.New("database error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &models.ErrorResponse{
				Error:   "Failed to retrieve categories",
				Details: "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmService := new(MockFilmService)
			mockCommentService := new(MockCommentService)
			handler := handlers.NewFilmHandler(mockFilmService, mockCommentService)

			// Setup mock expectations
			mockFilmService.On("GetCategories", mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			w := httptest.NewRecorder()

			// Execute handler
			handler.GetCategories(w, req)

			// Assert response
			require.Equal(t, tt.expectedStatusCode, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Convert expected response to JSON and back for comparison
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			var expectedInterface interface{}
			json.Unmarshal(expectedJSON, &expectedInterface)

			assert.Equal(t, expectedInterface, response)
			mockFilmService.AssertExpectations(t)
		})
	}
}

func TestFilmHandler_GetComments(t *testing.T) {
	tests := []struct {
		name               string
		filmID             string
		mockResponse       []models.Comment
		mockError          error
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:   "successful retrieval",
			filmID: "1",
			mockResponse: []models.Comment{
				{ID: 1, FilmID: 1, CustomerName: "John Doe", Comment: "Great movie!"},
				{ID: 2, FilmID: 1, CustomerName: "Jane Smith", Comment: "Loved it!"},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: []models.Comment{
				{ID: 1, FilmID: 1, CustomerName: "John Doe", Comment: "Great movie!"},
				{ID: 2, FilmID: 1, CustomerName: "Jane Smith", Comment: "Loved it!"},
			},
		},
		{
			name:               "film not found",
			filmID:             "999",
			mockError:          repository.ErrFilmNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &models.ErrorResponse{
				Error:   "Film not found",
				Details: "film not found",
			},
		},
		{
			name:               "invalid film ID",
			filmID:             "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &models.ErrorResponse{
				Error:   "Invalid film ID",
				Details: "strconv.Atoi: parsing \"invalid\": invalid syntax",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmService := new(MockFilmService)
			mockCommentService := new(MockCommentService)
			handler := handlers.NewFilmHandler(mockFilmService, mockCommentService)

			// Setup mock expectations only for valid film IDs
			if tt.filmID != "invalid" {
				filmID := 1
				if tt.filmID == "999" {
					filmID = 999
				}
				mockCommentService.On("GetCommentsByFilmID", mock.Anything, filmID).
					Return(tt.mockResponse, tt.mockError)
			}

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodGet, "/films/"+tt.filmID+"/comments", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.filmID})
			w := httptest.NewRecorder()

			// Execute handler
			handler.GetComments(w, req)

			// Assert response
			require.Equal(t, tt.expectedStatusCode, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Convert expected response to JSON and back for comparison
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			var expectedInterface interface{}
			json.Unmarshal(expectedJSON, &expectedInterface)

			assert.Equal(t, expectedInterface, response)
			mockCommentService.AssertExpectations(t)
		})
	}
}
