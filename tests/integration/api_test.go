package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/rxbenefits/go-hw/internal/handlers"
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

type IntegrationTestSuite struct {
	suite.Suite

	mockFilmRepo    *MockFilmRepository
	mockCommentRepo *MockCommentRepository
	router          *mux.Router
	filmHandler     *handlers.FilmHandler
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Initialize mock repositories
	suite.mockFilmRepo = new(MockFilmRepository)
	suite.mockCommentRepo = new(MockCommentRepository)

	// Initialize services with mock repositories
	filmService := service.NewFilmService(suite.mockFilmRepo)
	commentService := service.NewCommentService(suite.mockCommentRepo, suite.mockFilmRepo)

	// Initialize handlers
	suite.filmHandler = handlers.NewFilmHandler(filmService, commentService)

	// Setup router
	suite.router = mux.NewRouter()
	api := suite.router.PathPrefix("/api/v1").Subrouter()

	// Film routes
	api.HandleFunc("/films", suite.filmHandler.GetFilms).Methods("GET")
	api.HandleFunc("/films/{id}", suite.filmHandler.GetFilmByID).Methods("GET")
	api.HandleFunc("/categories", suite.filmHandler.GetCategories).Methods("GET")

	// Comment routes
	api.HandleFunc("/films/{id}/comments", suite.filmHandler.AddComment).Methods("POST")
	api.HandleFunc("/films/{id}/comments", suite.filmHandler.GetComments).Methods("GET")

	// Welcome route
	suite.router.HandleFunc("/", handlers.WelcomeHandler).Methods("GET")
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up mocks
	suite.mockFilmRepo.AssertExpectations(suite.T())
	suite.mockCommentRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Reset mock expectations before each test
	suite.mockFilmRepo.ExpectedCalls = nil
	suite.mockCommentRepo.ExpectedCalls = nil
}

func (suite *IntegrationTestSuite) TestWelcomeEndpoint() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response models.WelcomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("Welcome to Mockbuster Movie API!", response.Message)
}

func (suite *IntegrationTestSuite) TestGetFilms() {
	// Setup mock expectations
	expectedFilters := models.FilmFilters{
		Page:  1,
		Limit: 5,
	}
	mockResponse := &models.FilmListResponse{
		Films: []models.Film{
			{FilmID: 1, Title: "Test Film 1", Rating: "PG"},
			{FilmID: 2, Title: "Test Film 2", Rating: "G"},
		},
		Total: 2,
		Page:  1,
		Limit: 5,
	}
	suite.mockFilmRepo.On("GetFilms", expectedFilters).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/films?page=1&limit=5", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response models.FilmListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response.Films, 2)
	suite.Equal(1, response.Page)
	suite.Equal(5, response.Limit)
	suite.Equal(2, response.Total)
}

func (suite *IntegrationTestSuite) TestGetFilmsWithFilters() {
	// Setup mock expectations
	expectedFilters := models.FilmFilters{
		Title:  "Academy",
		Rating: "PG",
		Page:   1,
		Limit:  10,
	}
	mockResponse := &models.FilmListResponse{
		Films: []models.Film{
			{FilmID: 1, Title: "Academy Dinosaur", Rating: "PG"},
		},
		Total: 1,
		Page:  1,
		Limit: 10,
	}
	suite.mockFilmRepo.On("GetFilms", expectedFilters).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/films?title=Academy&rating=PG", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response models.FilmListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Verify filtering works
	suite.Len(response.Films, 1)
	suite.Contains(response.Films[0].Title, "Academy")
	suite.Equal("PG", response.Films[0].Rating)
}

func (suite *IntegrationTestSuite) TestGetFilmByID() {
	// Setup mock expectations
	filmID := 1
	description := "A test film"
	releaseYear := 2023
	length := 120
	lastUpdate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockFilm := &models.Film{
		FilmID:          1,
		Title:           "Test Film",
		Description:     &description,
		ReleaseYear:     &releaseYear,
		LanguageID:      1,
		RentalDuration:  3,
		RentalRate:      4.99,
		Length:          &length,
		ReplacementCost: 19.99,
		Rating:          "PG",
		SpecialFeatures: []string{"Trailers", "Commentaries"},
		LastUpdate:      lastUpdate,
	}
	suite.mockFilmRepo.On("GetFilmByID", filmID).Return(mockFilm, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/films/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response models.Film
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(1, response.FilmID)
	suite.Equal("Test Film", response.Title)
}

func (suite *IntegrationTestSuite) TestGetFilmByIDNotFound() {
	// Setup mock expectations
	filmID := 99999
	suite.mockFilmRepo.On("GetFilmByID", filmID).Return(nil, repository.ErrFilmNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/films/99999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "99999"})
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("Film not found", response.Error)
}

func (suite *IntegrationTestSuite) TestGetCategories() {
	// Setup mock expectations
	mockCategories := []models.Category{
		{CategoryID: 1, Name: "Action"},
		{CategoryID: 2, Name: "Comedy"},
		{CategoryID: 3, Name: "Drama"},
	}
	suite.mockFilmRepo.On("GetCategories").Return(mockCategories, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response, 3)

	// Verify category structure
	for i, category := range response {
		suite.Equal(mockCategories[i].CategoryID, category.CategoryID)
		suite.Equal(mockCategories[i].Name, category.Name)
	}
}

func (suite *IntegrationTestSuite) TestAddAndGetComments() {
	filmID := 1

	// Setup mock expectations for film existence check
	mockFilm := &models.Film{FilmID: 1, Title: "Test Film"}
	suite.mockFilmRepo.On("GetFilmByID", filmID).Return(mockFilm, nil)

	// Setup mock expectations for adding comment
	commentReq := models.CommentRequest{
		CustomerName: "Integration Test User",
		Comment:      "This is a test comment from integration test",
	}
	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockComment := &models.Comment{
		ID:           1,
		FilmID:       filmID,
		CustomerName: commentReq.CustomerName,
		Comment:      commentReq.Comment,
		CreatedAt:    createdAt,
	}
	suite.mockCommentRepo.On("AddComment", filmID, commentReq).Return(mockComment, nil)

	// First, add a comment
	requestBody, _ := json.Marshal(commentReq)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/films/"+strconv.Itoa(filmID)+"/comments",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(filmID)})
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	var addResponse models.Comment
	err := json.Unmarshal(w.Body.Bytes(), &addResponse)
	suite.Require().NoError(err)
	suite.Equal(1, addResponse.ID)
	suite.Equal(filmID, addResponse.FilmID)
	suite.Equal(commentReq.CustomerName, addResponse.CustomerName)
	suite.Equal(commentReq.Comment, addResponse.Comment)

	// Setup mock expectations for getting comments
	mockComments := []models.Comment{*mockComment}
	suite.mockCommentRepo.On("GetCommentsByFilmID", filmID).Return(mockComments, nil)

	// Now, get comments for the film
	req = httptest.NewRequest(http.MethodGet, "/api/v1/films/"+strconv.Itoa(filmID)+"/comments", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(filmID)})
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var getResponse []models.Comment
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	suite.Require().NoError(err)
	suite.Len(getResponse, 1)

	// Verify our comment is in the list
	suite.Equal(addResponse.ID, getResponse[0].ID)
	suite.Equal(commentReq.CustomerName, getResponse[0].CustomerName)
	suite.Equal(commentReq.Comment, getResponse[0].Comment)
}

func (suite *IntegrationTestSuite) TestAddCommentToNonExistentFilm() {
	filmID := 99999

	// Setup mock expectations for film not found
	suite.mockFilmRepo.On("GetFilmByID", filmID).Return(nil, repository.ErrFilmNotFound)

	commentReq := models.CommentRequest{
		CustomerName: "Test User",
		Comment:      "This should fail",
	}

	requestBody, _ := json.Marshal(commentReq)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/films/"+strconv.Itoa(filmID)+"/comments",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(filmID)})
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("Film not found", response.Error)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
