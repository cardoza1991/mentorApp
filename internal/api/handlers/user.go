package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"mentorApp/internal/api/handlers/common"
	"mentorApp/internal/services"
)

type UserHandler struct {
	service services.IUserService
}

func NewUserHandler(service services.IUserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterMentee(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Bio       string `json:"bio"`
		Skills    string `json:"skills"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	input := services.RegisterUserInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		Skills:    req.Skills,
		IsMentor:  false, // Ensures this user is a mentee
	}

	user, err := h.service.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Bio       string `json:"bio"`
		Skills    string `json:"skills"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	input := services.RegisterUserInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		Skills:    req.Skills,
	}

	user, err := h.service.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
	})
}

// RegisterMentor handles mentor registration
func (h *UserHandler) RegisterMentor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string  `json:"username"`
		Email       string  `json:"email"`
		Password    string  `json:"password"`
		FirstName   string  `json:"first_name"`
		LastName    string  `json:"last_name"`
		Bio         string  `json:"bio"`
		Skills      string  `json:"skills"`
		Rate        float64 `json:"rate"`
		Experience  string  `json:"experience"`
		Specialties string  `json:"specialties"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	input := services.RegisterMentorInput{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Bio:         req.Bio,
		Skills:      req.Skills,
		Rate:        req.Rate,
		Experience:  req.Experience,
		Specialties: req.Specialties,
	}

	user, err := h.service.RegisterMentor(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login handles user authentication
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	user, err := h.service.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cookieName := "session_token"
	if user.IsAdmin {
		cookieName = "admin_session_token"
	} else if user.IsMentor {
		cookieName = "mentor_session_token"
	} else {
		cookieName = "mentee_session_token"
	}

	// Create session token
	sessionToken := fmt.Sprintf("session_%d_%s", user.Id, uuid.New().String())

	// Set session cookie with unique name
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	common.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"id":        user.Id,
		"username":  user.Username,
		"is_mentor": user.IsMentor,
		"is_admin":  user.IsAdmin,
	})
}

// Logout handles user logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear all possible session cookies
	cookieNames := []string{"session_token", "admin_session_token", "mentor_session_token", "mentee_session_token"}

	for _, name := range cookieNames {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   -1,
		})
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := h.service.VerifyEmail(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Email verified successfully"})
}

// RequestPasswordReset initiates password reset process
func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password reset email sent"})
}

// ResetPassword completes the password reset process
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.service.ResetPassword(r.Context(), token, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
