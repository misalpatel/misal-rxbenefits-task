package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rxbenefits/go-hw/internal/models"
	"github.com/rxbenefits/go-hw/internal/repository"
	"github.com/rxbenefits/go-hw/internal/service"
)

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) AddComment(filmID int, commentReq models.CommentRequest) (*models.Comment, error) {
	args := m.Called(filmID, commentReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByFilmID(filmID int) ([]models.Comment, error) {
	args := m.Called(filmID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func TestCommentService_AddComment(t *testing.T) {
	tests := []struct {
		name           string
		filmID         int
		commentReq     models.CommentRequest
		filmExists     bool
		filmError      error
		mockResponse   *models.Comment
		mockError      error
		expectedResult *models.Comment
		expectedError  string
	}{
		{
			name:   "successful comment addition",
			filmID: 1,
			commentReq: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			filmExists: true,
			mockResponse: &models.Comment{
				ID:           1,
				FilmID:       1,
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			expectedResult: &models.Comment{
				ID:           1,
				FilmID:       1,
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
		},
		{
			name:   "film not found",
			filmID: 999,
			commentReq: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			filmError:     repository.ErrFilmNotFound,
			expectedError: "film not found",
		},
		{
			name:   "invalid film ID",
			filmID: 0,
			commentReq: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "Great movie!",
			},
			expectedError: "invalid film ID",
		},
		{
			name:   "empty customer name",
			filmID: 1,
			commentReq: models.CommentRequest{
				CustomerName: "",
				Comment:      "Great movie!",
			},
			expectedError: "customer name is required",
		},
		{
			name:   "empty comment",
			filmID: 1,
			commentReq: models.CommentRequest{
				CustomerName: "John Doe",
				Comment:      "",
			},
			expectedError: "comment text is required",
		},
		{
			name:   "customer name too long",
			filmID: 1,
			commentReq: models.CommentRequest{
				CustomerName: "This is a very long customer name that exceeds the maximum allowed length of 100 characters for the customer name field",
				Comment:      "Great movie!",
			},
			expectedError: "customer name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFilmRepo := new(MockFilmRepository)
			mockCommentRepo := new(MockCommentRepository)
			commentService := service.NewCommentService(mockCommentRepo, mockFilmRepo)

			// Setup film existence check if filmID is valid
			if tt.filmID > 0 && tt.expectedError != "customer name is required" &&
				tt.expectedError != "comment text is required" &&
				tt.expectedError != "customer name too long" {
				if tt.filmExists {
					mockFilmRepo.On("GetFilmByID", tt.filmID).Return(&models.Film{FilmID: tt.filmID}, tt.filmError)
					if tt.filmError == nil {
						mockCommentRepo.On("AddComment", tt.filmID, tt.commentReq).Return(tt.mockResponse, tt.mockError)
					}
				} else {
					mockFilmRepo.On("GetFilmByID", tt.filmID).Return(nil, tt.filmError)
				}
			}

			result, err := commentService.AddComment(context.Background(), tt.filmID, tt.commentReq)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockFilmRepo.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)
		})
	}
}

func TestCommentService_GetCommentsByFilmID(t *testing.T) {
	tests := []struct {
		name           string
		filmID         int
		filmExists     bool
		filmError      error
		mockResponse   []models.Comment
		mockError      error
		expectedResult []models.Comment
		expectedError  string
	}{
		{
			name:       "successful retrieval",
			filmID:     1,
			filmExists: true,
			mockResponse: []models.Comment{
				{ID: 1, FilmID: 1, CustomerName: "John Doe", Comment: "Great movie!"},
				{ID: 2, FilmID: 1, CustomerName: "Jane Smith", Comment: "Loved it!"},
			},
			expectedResult: []models.Comment{
				{ID: 1, FilmID: 1, CustomerName: "John Doe", Comment: "Great movie!"},
				{ID: 2, FilmID: 1, CustomerName: "Jane Smith", Comment: "Loved it!"},
			},
		},
		{
			name:          "film not found",
			filmID:        999,
			filmError:     repository.ErrFilmNotFound,
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
			mockFilmRepo := new(MockFilmRepository)
			mockCommentRepo := new(MockCommentRepository)
			commentService := service.NewCommentService(mockCommentRepo, mockFilmRepo)

			if tt.filmID > 0 {
				if tt.filmExists {
					mockFilmRepo.On("GetFilmByID", tt.filmID).Return(&models.Film{FilmID: tt.filmID}, tt.filmError)
					if tt.filmError == nil {
						mockCommentRepo.On("GetCommentsByFilmID", tt.filmID).Return(tt.mockResponse, tt.mockError)
					}
				} else {
					mockFilmRepo.On("GetFilmByID", tt.filmID).Return(nil, tt.filmError)
				}
			}

			result, err := commentService.GetCommentsByFilmID(context.Background(), tt.filmID)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockFilmRepo.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)
		})
	}
}
