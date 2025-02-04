package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"mentorApp/internal/services"
)

type HomeHandler struct {
	templates     *template.Template
	userService   services.IUserService
	mentorService services.IMentorshipService
}

func NewHomeHandler(userService services.IUserService, mentorService services.IMentorshipService) *HomeHandler {
	// Create a new template instance
	tmpl := template.New("")

	// Parse templates once
	var allFiles []string

	// Add base templates
	files, err := filepath.Glob("views/*.html")
	if err != nil {
		log.Fatalf("Failed to glob base templates: %v", err)
	}
	allFiles = append(allFiles, files...)

	// Add templates from subdirectories
	subdirs := []string{"mentor", "mentee", "jobs"}
	for _, subdir := range subdirs {
		files, err := filepath.Glob(fmt.Sprintf("views/%s/*.html", subdir))
		if err != nil {
			log.Fatalf("Failed to glob %s templates: %v", subdir, err)
		}
		allFiles = append(allFiles, files...)
	}

	// Parse all templates at once
	tmpl, err = tmpl.ParseFiles(allFiles...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Debug: log templates only once
	log.Printf("Loaded templates:")
	for _, t := range tmpl.Templates() {
		log.Printf("- %s", t.Name())
	}

	return &HomeHandler{
		templates:     tmpl,
		userService:   userService,
		mentorService: mentorService,
	}
}

func (h *HomeHandler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	// List all available templates for debugging
	templates := []string{}
	for _, t := range h.templates.Templates() {
		templates = append(templates, t.Name())
	}
	log.Printf("Available templates: %v", templates)
	log.Printf("Attempting to render template: %s", templateName)

	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template error for %s: %v", templateName, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Get renders the home page
func (h *HomeHandler) Get(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Website": "NEXUS Mentorship Platform",
		"Email":   "contact@nexusmentors.org",
	}

	if userID := r.Context().Value("userID"); userID != nil {
		featuredMentors, err := h.mentorService.GetFeaturedMentors(r.Context())
		if err == nil {
			data["FeaturedMentors"] = featuredMentors
		}
	}

	h.renderTemplate(w, "index.html", data)
}

func (h *HomeHandler) GetJobBoard(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.mentorService.GetAvailableJobs(r.Context())
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Website":     "NEXUS Mentorship Platform",
		"Jobs":        jobs,
		"LastUpdated": time.Now().Format(time.RFC822),
	}

	h.renderTemplate(w, "job_board.html", data)
}

func (h *HomeHandler) GetMenteeDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	dataChan := make(chan map[string]interface{}, 4)
	errChan := make(chan error, 4)

	go func() {
		profile, err := h.userService.GetUserProfile(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"Profile": profile}
	}()

	go func() {
		mentorships, err := h.mentorService.GetActiveMentorships(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"ActiveMentorships": mentorships}
	}()

	go func() {
		sessions, err := h.mentorService.GetUpcomingSessions(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"UpcomingSessions": sessions}
	}()

	go func() {
		mentors, err := h.mentorService.GetRecommendedMentors(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"RecommendedMentors": mentors}
	}()

	dashboardData := make(map[string]interface{})
	dashboardData["Website"] = "NEXUS Mentorship Platform"

	for i := 0; i < 4; i++ {
		select {
		case data := <-dataChan:
			for k, v := range data {
				dashboardData[k] = v
			}
		case err := <-errChan:
			http.Error(w, "Error fetching dashboard data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.renderTemplate(w, "mentee_dashboard.html", dashboardData)
}

// Update the HomeHandler to use the correct template path
func (h *HomeHandler) GetMentorDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	dataChan := make(chan map[string]interface{}, 5)
	errChan := make(chan error, 5)

	go func() {
		profile, err := h.userService.GetUserProfile(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"Profile": profile}
	}()

	go func() {
		programs, err := h.mentorService.GetMentorPrograms(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"Programs": programs}
	}()

	go func() {
		sessions, err := h.mentorService.GetUpcomingSessions(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"UpcomingSessions": sessions}
	}()

	go func() {
		requests, err := h.mentorService.GetPendingRequests(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"PendingRequests": requests}
	}()

	go func() {
		analytics, err := h.mentorService.GetMentorAnalytics(r.Context(), userID)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- map[string]interface{}{"Analytics": analytics}
	}()

	dashboardData := make(map[string]interface{})
	dashboardData["Website"] = "NEXUS Mentorship Platform"

	for i := 0; i < 5; i++ {
		select {
		case data := <-dataChan:
			for k, v := range data {
				dashboardData[k] = v
			}
		case err := <-errChan:
			log.Printf("Error fetching dashboard data: %v", err)
		}
	}

	// Debug template name
	log.Printf("Attempting to render template: mentor_dashboard.html")

	if err := h.templates.ExecuteTemplate(w, "mentor_dashboard.html", dashboardData); err != nil {
		log.Printf("Template error: %v", err)
		// List all available templates for debugging
		var templateNames []string
		for _, t := range h.templates.Templates() {
			templateNames = append(templateNames, t.Name())
		}
		log.Printf("Available templates: %v", templateNames)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	h.renderTemplate(w, "mentor_dashboard.html", dashboardData)
}

func (h *HomeHandler) GetMenteeRegistration(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Website": "NEXUS Mentorship Platform",
		"Title":   "Mentee Registration",
	}

	h.renderTemplate(w, "mentee_registration.html", data)
}

func (h *HomeHandler) GetMentorRegistration(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Website":     "NEXUS Mentorship Platform",
		"Title":       "Mentor Registration",
		"Specialties": h.mentorService.GetAvailableSpecialties(r.Context()),
	}

	h.renderTemplate(w, "mentor_registration.html", data)
}

// GetLogin handles the login page
func (h *HomeHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Website":    "NEXUS Mentorship Platform",
		"Registered": r.URL.Query().Get("registered") == "true",
	}

	h.renderTemplate(w, "login.html", data)
}

func (h *HomeHandler) renderError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	data := map[string]interface{}{
		"Error":   message,
		"Website": "NEXUS Mentorship Platform",
	}
	h.renderTemplate(w, "error.html", data)
}
