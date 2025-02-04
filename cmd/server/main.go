package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"


	"mentorApp/internal/api/handlers"
	"mentorApp/internal/api/routes"
	"mentorApp/internal/db"
	"mentorApp/internal/middleware"
	"mentorApp/internal/repository"
	"mentorApp/internal/services"
	"mentorApp/pkg/utils/email"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[NEXUS] ", log.LstdFlags)
	logger.Println("Starting server...")

	// Get database configuration from environment or use defaults
	dbConfig := getDatabaseConfig()

	// Construct database URL for migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	// Run migrations
	if err := db.MigrateDB(dbURL); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database connection
	db, err := initDatabase(dbConfig)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Could not ping database: %v", err)
	}
	logger.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	mentorshipRepo := repository.NewMentorshipRepository(db)

	// Initialize services
	emailSvc := email.NewEmailService("noreply@nexusmentors.org")
	userService := services.NewUserService(userRepo, profileRepo, emailSvc)
	mentorshipService := services.NewMentorshipService(mentorshipRepo, profileRepo, userRepo)

	// Initialize templates with recursive glob
	var allTemplates []string
	templatePaths := []string{
		"views/*.html",
		"views/admin/*.html",
		"views/mentor/*.html",
		"views/mentee/*.html",
		"views/jobs/*.html",
	}

	for _, pattern := range templatePaths {
		files, err := filepath.Glob(pattern)
		if err != nil {
			logger.Fatalf("Failed to glob pattern %s: %v", pattern, err)
		}
		allTemplates = append(allTemplates, files...)
	}

	// Parse all templates
	templates, err := template.ParseFiles(allTemplates...)
	if err != nil {
		logger.Fatalf("Failed to parse templates: %v", err)
	}

	// Log all registered templates for debugging
	for _, t := range templates.Templates() {
		logger.Printf("Registered template: %s", t.Name())
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	mentorshipHandler := handlers.NewMentorshipHandler(mentorshipService)
	profileHandler := handlers.NewProfileHandler(userService)
	homeHandler := handlers.NewHomeHandler(userService, mentorshipService)
	adminHandler := handlers.NewAdminHandler(db, userRepo, profileRepo, templates)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(middleware.LoggingMiddleware)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Setup routes
	routes.SetupRoutes(r, userHandler, mentorshipHandler, profileHandler, homeHandler, adminHandler)

	// Server configuration
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", getEnv("HOST", "10.20.0.1"), getEnv("PORT", "8080")),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown setup
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	// Start server
	logger.Printf("Server is ready to handle requests at :%s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5433"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "nexus"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func initDatabase(config DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
