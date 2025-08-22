// Package main provides the entry point for the Mockbuster Movie API.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/rxbenefits/go-hw/docs"
	"github.com/rxbenefits/go-hw/internal/database"
	"github.com/rxbenefits/go-hw/internal/handlers"
	"github.com/rxbenefits/go-hw/internal/repository"
	"github.com/rxbenefits/go-hw/internal/service"
	"github.com/rxbenefits/go-hw/internal/util"
)

const (
	readTimeout  = 15 * time.Second
	writeTimeout = 15 * time.Second
	idleTimeout  = 60 * time.Second
)

// @title Mockbuster Movie API.
// @version 1.0

// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080.
// @BasePath /.
// @schemes http.

func main() {
	// Initialize database connection.
	config := util.InitConfig()
	db, err := database.InitDB(
		database.WithDBHost(config.DBHost),
		database.WithDBPort(config.DBPort),
		database.WithDBUser(config.DBUser),
		database.WithDBPassword(config.DBPassword),
		database.WithDBName(config.DBName),
	)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories.
	filmRepo := repository.NewFilmRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Run database migrations.
	if migrationErr := database.RunMigrations(db.DB, "migrations"); migrationErr != nil {
		slog.Error("Failed to run database migrations", "error", migrationErr)
		db.Close() //nolint:gosec // Exiting the program anyways
		os.Exit(1) //nolint:gocritic // Running the db.Close() before os.Exit
	}

	// Initialize services with dependency injection.
	filmService := service.NewFilmService(filmRepo)
	commentService := service.NewCommentService(commentRepo, filmRepo)

	// Initialize handlers with services.
	filmHandler := handlers.NewFilmHandler(filmService, commentService)

	// Initialize router.
	r := mux.NewRouter()

	// API routes.
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", handlers.APIInfoHandler).Methods("GET")

	// Film routes.
	api.HandleFunc("/films", filmHandler.GetFilms).Methods("GET")
	api.HandleFunc("/films/{id}", filmHandler.GetFilmByID).Methods("GET")
	api.HandleFunc("/categories", filmHandler.GetCategories).Methods("GET")

	// Comment routes.
	api.HandleFunc("/films/{id}/comments", filmHandler.AddComment).Methods("POST")
	api.HandleFunc("/films/{id}/comments", filmHandler.GetComments).Methods("GET")

	// Welcome route.
	r.HandleFunc("/", handlers.WelcomeHandler).Methods("GET")

	// Swagger documentation.
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Set Swagger info
	docs.SwaggerInfo.Title = "Mockbuster Movie API"
	docs.SwaggerInfo.Description = "A RESTful API for the Mockbuster DVD rental business"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// CORS middleware.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Apply CORS middleware.
	handler := c.Handler(r)

	// Get port from environment or use default.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure server with timeouts.
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	slog.Info("Starting Mockbuster Movie API server", "port", port)
	slog.Info("API Documentation available", "url", "http://localhost:"+port+"/swagger/")
	slog.Info("API Base URL", "url", "http://localhost:"+port+"/api/v1")

	if serveErr := server.ListenAndServe(); serveErr != nil {
		slog.Error("Failed to start server", "error", serveErr)
		err = db.Close()
		if err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
		os.Exit(1)
	}
}
