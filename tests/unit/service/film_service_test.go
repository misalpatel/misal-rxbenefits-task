package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
	"github.com/rxbenefits/go-hw/internal/service"
)

type MockFilmRepository struct {
	mock.Mock
}

func (m *MockFilmRepository) GetFilms(filters models.FilmFilters) (*models.FilmListResponse, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FilmListResponse), args.Error(1)
}

func (m *MockFilmRepository) GetFilmByID(filmID int) (*models.Film, error) {
	args := m.Called(filmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Film), args.Error(1)
}

func (m *MockFilmRepository) GetCategories() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func TestFilmService_GetFilms(t *testing.T) {
	tests := []struct {
		name           string
		filters        models.FilmFilters
		mockResponse   *models.FilmListResponse
		mockError      error
		expectedResult *models.FilmListResponse
		expectedError  string
	}{
		{
			name: "successful retrieval with valid filters",
			filters: models.FilmFilters{
				Title:  "Test",
				Rating: "PG",
				Page:   1,
				Limit:  10,
			},
			mockResponse: &models.FilmListResponse{
				Films: []models.Film{
					{FilmID: 1, Title: "Test Film", Rating: "PG"},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			},
			expectedResult: &models.FilmListResponse{
				Films: []models.Film{
					{FilmID: 1, Title: "Test Film", Rating: "PG"},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			},
		},
		{
			name: "invalid rating filter",
			filters: models.FilmFilters{
				Rating: "INVALID",
				Page:   1,
				Limit:  10,
			},
			expectedError: "invalid rating provided",
		},
		{
			name: "invalid page number",
			filters: models.FilmFilters{
				Page:  0,
				Limit: 10,
			},
			expectedError: "page must be greater than 0",
		},
		{
			name: "invalid limit",
			filters: models.FilmFilters{
				Page:  1,
				Limit: 101,
			},
			expectedError: "limit must be between 1 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFilmRepository)
			filmService := service.NewFilmService(mockRepo)

			if tt.mockResponse != nil {
				// Normalize filters for the mock expectation
				expectedFilters := tt.filters
				if expectedFilters.Page <= 0 {
					expectedFilters.Page = 1
				}
				if expectedFilters.Limit <= 0 {
					expectedFilters.Limit = 10
				}
				mockRepo.On("GetFilms", expectedFilters).Return(tt.mockResponse, tt.mockError)
			}

			result, err := filmService.GetFilms(context.Background(), tt.filters)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFilmService_GetFilmByID(t *testing.T) {
	tests := []struct {
		name           string
		filmID         int
		mockResponse   *models.Film
		mockError      error
		expectedResult *models.Film
		expectedError  string
	}{
		{
			name:   "successful retrieval",
			filmID: 1,
			mockResponse: &models.Film{
				FilmID: 1,
				Title:  "Test Film",
				Rating: "PG",
			},
			expectedResult: &models.Film{
				FilmID: 1,
				Title:  "Test Film",
				Rating: "PG",
			},
		},
		{
			name:          "film not found",
			filmID:        999,
			mockError:     repository.ErrFilmNotFound,
			expectedError: "film not found",
		},
		{
			name:          "invalid film ID",
			filmID:        0,
			expectedError: "invalid film ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFilmRepository)
			filmService := service.NewFilmService(mockRepo)

			if tt.filmID > 0 {
				mockRepo.On("GetFilmByID", tt.filmID).Return(tt.mockResponse, tt.mockError)
			}

			result, err := filmService.GetFilmByID(context.Background(), tt.filmID)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFilmService_GetCategories(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   []models.Category
		mockError      error
		expectedResult []models.Category
		expectedError  bool
	}{
		{
			name: "successful retrieval",
			mockResponse: []models.Category{
				{CategoryID: 1, Name: "Action"},
				{CategoryID: 2, Name: "Comedy"},
			},
			expectedResult: []models.Category{
				{CategoryID: 1, Name: "Action"},
				{CategoryID: 2, Name: "Comedy"},
			},
		},
		{
			name:          "repository error",
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockFilmRepository)
			filmService := service.NewFilmService(mockRepo)

			mockRepo.On("GetCategories").Return(tt.mockResponse, tt.mockError)

			result, err := filmService.GetCategories(context.Background())

			if tt.expectedError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
