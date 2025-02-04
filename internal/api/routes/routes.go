package routes

import (
	"mentorApp/internal/api/handlers"
	"mentorApp/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	r *chi.Mux,
	userHandler *handlers.UserHandler,
	mentorshipHandler *handlers.MentorshipHandler,
	profileHandler *handlers.ProfileHandler,
	homeHandler *handlers.HomeHandler,
	adminHandler *handlers.AdminHandler,
) {
	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		// Core pages
		r.Get("/", homeHandler.Get)
		r.Get("/jobs", homeHandler.GetJobBoard)
		r.Get("/mentee/register", homeHandler.GetMenteeRegistration)
		r.Get("/mentor/register", homeHandler.GetMentorRegistration)
		r.Get("/login", homeHandler.GetLogin)

		// Auth endpoints
		r.Post("/auth/register", userHandler.Register)
		r.Post("/auth/login", userHandler.Login)
		r.Get("/auth/logout", userHandler.Logout)

		// Registration endpoints
		r.Post("/register/mentee", userHandler.RegisterMentee)
		r.Post("/register/mentor", userHandler.RegisterMentor)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Profile routes
		r.Get("/profile", profileHandler.GetProfile)
		r.Put("/profile", profileHandler.UpdateProfile)

		// Mentor routes
		r.Route("/mentor", func(r chi.Router) {
			r.Use(middleware.MentorRequired)
			r.Get("/dashboard", homeHandler.GetMentorDashboard) // Make sure this exists
			r.Get("/programs", mentorshipHandler.ListMentorPrograms)
			r.Post("/programs", mentorshipHandler.CreateProgram)
			r.Get("/requests", mentorshipHandler.ListMentorshipRequests)
			r.Put("/requests/{requestId}", mentorshipHandler.RespondToRequest)
			
		})

		// Mentee routes
		r.Route("/mentee", func(r chi.Router) {
			r.Get("/programs", mentorshipHandler.ListAvailablePrograms)
			r.Post("/request/{programId}", mentorshipHandler.RequestMentorship)
			r.Get("/sessions", mentorshipHandler.ListMenteeSessions)
		})
	})

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		// Public admin setup
		r.Get("/setup", adminHandler.Setup)
		r.Post("/setup", adminHandler.Setup)

		// Protected admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Use(middleware.RequireAdmin(adminHandler.DB()))

			// Main admin dashboard
			r.Get("/", adminHandler.Dashboard)
			r.Get("/dashboard", adminHandler.Dashboard)

			// User management
			r.Get("/users", adminHandler.ListUsers)
			r.Get("/profiles", adminHandler.ListProfiles)
			r.Post("/mentors/{userId}/approve", adminHandler.ApproveMentor)

			// Job management
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", adminHandler.ListJobs)
				r.Post("/", adminHandler.CreateJob)
				r.Delete("/{jobId}", adminHandler.DeleteJob)
				r.Post("/{jobId}/feature", adminHandler.FeatureJob)
			})
		})
	})

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
}
