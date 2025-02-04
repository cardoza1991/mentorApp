package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mentorApp/internal/api/handlers/common"
	"mentorApp/internal/models"
	"mentorApp/internal/services"

	"github.com/go-chi/chi/v5"
)

type ProfileHandler struct {
	service services.IUserService
}

func NewProfileHandler(service services.IUserService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	profile, err := h.service.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		Bio            string  `json:"bio"`
		Skills         string  `json:"skills"`
		Rate           float64 `json:"rate,omitempty"`
		TimeZone       string  `json:"timezone"`
		Availability   string  `json:"availability"`
		LinkedIn       string  `json:"linkedin"`
		Github         string  `json:"github"`
		ProfilePicture string  `json:"profile_picture"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	profile := &models.Profile{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Bio:            req.Bio,
		Skills:         req.Skills,
		Rate:           req.Rate,
		Timezone:       req.TimeZone,
		LinkedIn:       req.LinkedIn,
		Github:         req.Github,
		ProfilePicture: req.ProfilePicture,
	}

	if err := h.service.UpdateUserProfile(r.Context(), userID, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.service.GetPublicProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, profile)
}

// GetSettings returns user's profile settings
func (h *ProfileHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	settings, err := h.service.GetProfileSettings(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, settings)
}

// UpdateSettings handles settings updates
func (h *ProfileHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req models.ProfileSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	if err := h.service.UpdateProfileSettings(r.Context(), userID, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

// GetNotificationSettings returns user's notification preferences
func (h *ProfileHandler) GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	settings, err := h.service.GetNotificationSettings(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, settings)
}

// UpdateNotificationSettings handles notification preferences updates
func (h *ProfileHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	var req models.NotificationSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	if err := h.service.UpdateNotificationSettings(r.Context(), userID, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Notification settings updated successfully"})
}

// ApproveProfile handles profile approval requests
func (h *ProfileHandler) ApproveProfile(w http.ResponseWriter, r *http.Request) {
	// Only admin users should be able to approve profiles
	if !isAdmin(r.Context()) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.service.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile.IsApproved = true
	if err := h.service.UpdateUserProfile(r.Context(), userID, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Profile approved successfully",
	})
}

// CheckJobBoardAccess verifies if a user can access the job board
func (h *ProfileHandler) CheckJobBoardAccess(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	profile, err := h.service.GetUserProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !profile.IsApproved {
		http.Error(w, "Profile requires approval to access job board", http.StatusForbidden)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]bool{
		"hasAccess": true,
	})
}

// Helper function to check if user is admin
func isAdmin(ctx context.Context) bool {
	// Implement your admin check logic here
	// For now, return false as placeholder
	return false
}
