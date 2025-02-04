// File: internal/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionToken string

		// Check all possible session cookies
		if c, err := r.Cookie("admin_session_token"); err == nil {
			sessionToken = c.Value
		} else if c, err := r.Cookie("mentor_session_token"); err == nil {
			sessionToken = c.Value
		} else if c, err := r.Cookie("mentee_session_token"); err == nil {
			sessionToken = c.Value
		} else if c, err := r.Cookie("session_token"); err == nil {
			// For backward compatibility
			sessionToken = c.Value
		} else {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Extract userID from session token
		parts := strings.Split(sessionToken, "_")
		if len(parts) != 3 || parts[0] != "session" {
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}

		// Add userID to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to extract userID from context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("userID").(int)
	return userID, ok
}

func MentorRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve userID from context
		userID, ok := GetUserID(r.Context())
		if !ok {
			http.Error(w, "Access denied: user not authenticated", http.StatusUnauthorized)
			return
		}

		// Check if user is a mentor
		isMentor, err := checkUserIsMentor(userID)
		if err != nil || !isMentor {
			http.Error(w, "Mentor access required", http.StatusForbidden)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// APIKeyAuth validates API key for external services
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		if !isValidAPIKey(apiKey) {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions
func getUserIDFromSession(sessionToken string) int {
	parts := strings.Split(sessionToken, "_")
	if len(parts) != 3 {
		return 0
	}
	userID, _ := strconv.Atoi(parts[1])
	return userID
}

func checkUserIsMentor(userID int) (bool, error) {
	// Implement check against your database
	// This is a placeholder
	return true, nil
}

func isValidAPIKey(key string) bool {
	// Implement API key validation
	// This is a placeholder
	return true
}
