package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"mentorApp/internal/api/handlers/common"
	"mentorApp/internal/models"
	"mentorApp/internal/services"

	"github.com/go-chi/chi/v5"
)

type MentorshipHandler struct {
	service services.IMentorshipService
}

func NewMentorshipHandler(service services.IMentorshipService) *MentorshipHandler {
	return &MentorshipHandler{
		service: service,
	}
}

// In mentorship.go
func (h *MentorshipHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	// Get mentor ID from context
	mentorID := r.Context().Value("userID").(int)

	// Parse request body
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Duration    string  `json:"duration"`
		MaxMentees  int     `json:"max_mentees"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		common.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// Validate input
	if req.Title == "" || req.Description == "" || req.Duration == "" || req.Price < 0 || req.MaxMentees < 1 {
		common.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "All fields are required and must be valid",
		})
		return
	}

	// Create program object
	program := &models.MentorshipProgram{
		MentorID:    mentorID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,
		MaxMentees:  req.MaxMentees,
		Status:      "active", // Set default status
	}

	// Create the program
	err := h.service.CreateMentorshipProgram(r.Context(), mentorID, program)
	if err != nil {
		log.Printf("Error creating program: %v", err)
		common.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create program: " + err.Error(),
		})
		return
	}

	// Return success response
	common.RespondJSON(w, http.StatusCreated, program)
}

func (h *MentorshipHandler) ListAvailablePrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := h.service.ListAvailablePrograms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, programs)
}

func (h *MentorshipHandler) GetProgramDetails(w http.ResponseWriter, r *http.Request) {
	programID, err := strconv.Atoi(chi.URLParam(r, "programId"))
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	program, err := h.service.GetProgramDetails(r.Context(), programID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, program)
}

func (h *MentorshipHandler) ScheduleSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestId int    `json:"request_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Topic     string `json:"topic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	session := &models.MentorshipSession{
		RequestID: req.RequestId,
		StartTime: startTime,
		EndTime:   endTime,
		Topic:     req.Topic,
	}

	if err := h.service.ScheduleSession(r.Context(), userID, session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusCreated, session)
}

func (h *MentorshipHandler) GetMentorAnalytics(w http.ResponseWriter, r *http.Request) {
	mentorID := r.Context().Value("userID").(int)

	analytics, err := h.service.GetMentorAnalytics(r.Context(), mentorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, analytics)
}

func (h *MentorshipHandler) SearchMentors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := map[string]interface{}{
		"skills":    query.Get("skills"),
		"rate_min":  query.Get("rate_min"),
		"rate_max":  query.Get("rate_max"),
		"timezone":  query.Get("timezone"),
		"available": query.Get("available"),
	}

	results, err := h.service.SearchMentors(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, results)
}

func (h *MentorshipHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Availability []models.Availability `json:"availability"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	mentorID := r.Context().Value("userID").(int)

	if err := h.service.UpdateAvailability(r.Context(), mentorID, req.Availability); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Availability updated successfully"})
}

// ListMentorPrograms lists all programs created by the mentor (logged-in user)
func (h *MentorshipHandler) ListMentorPrograms(w http.ResponseWriter, r *http.Request) {
	mentorID := r.Context().Value("userID").(int)
	programs, err := h.service.ListMentorPrograms(r.Context(), mentorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	common.RespondJSON(w, http.StatusOK, programs)
}

// ListMentorshipRequests lists all mentorship requests for the mentor
func (h *MentorshipHandler) ListMentorshipRequests(w http.ResponseWriter, r *http.Request) {
	mentorID := r.Context().Value("userID").(int)
	requests, err := h.service.GetPendingRequests(r.Context(), mentorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	common.RespondJSON(w, http.StatusOK, requests)
}

// RespondToRequest allows the mentor to approve/reject a mentee's mentorship request
func (h *MentorshipHandler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	mentorID := r.Context().Value("userID").(int)
	requestID, err := strconv.Atoi(chi.URLParam(r, "requestId"))
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Approve bool `json:"approve"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.service.RespondToRequest(r.Context(), mentorID, requestID, req.Approve); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Request updated successfully"})
}

// RequestMentorship allows a mentee to request mentorship from a specific program
func (h *MentorshipHandler) RequestMentorship(w http.ResponseWriter, r *http.Request) {
	menteeID := r.Context().Value("userID").(int)
	programID, err := strconv.Atoi(chi.URLParam(r, "programId"))
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.service.RequestMentorship(r.Context(), menteeID, programID, req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Mentorship requested successfully"})
}

// ListMenteeSessions lists all sessions for the mentee
func (h *MentorshipHandler) ListMenteeSessions(w http.ResponseWriter, r *http.Request) {
	menteeID := r.Context().Value("userID").(int)

	// Assuming you have a method to get mentee sessions in MentorshipService
	// If you do not, implement something like GetMenteeSessions or GetActiveMentorships
	sessions, err := h.service.GetUpcomingSessions(r.Context(), menteeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, sessions)
}

func (h *MentorshipHandler) SubmitSessionFeedback(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.Atoi(chi.URLParam(r, "sessionId"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	feedback := &models.SessionFeedback{
		SessionID: sessionID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := h.service.SubmitSessionFeedback(r.Context(), feedback); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusCreated, map[string]string{"message": "Feedback submitted successfully"})
}
