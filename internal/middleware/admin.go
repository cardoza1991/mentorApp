package middleware

import (
	"database/sql"
	"net/http"
)

// IsAdmin checks if a user has admin privileges
func IsAdmin(db *sql.DB, userID interface{}) bool {
	// If userID is nil, user is not logged in
	if userID == nil {
		return false
	}

	// Convert userID to int
	id, ok := userID.(int)
	if !ok {
		return false
	}

	var isAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE id = $1", id).Scan(&isAdmin)
	if err != nil {
		return false
	}
	return isAdmin
}

// RequireAdmin middleware ensures the user is an admin
func RequireAdmin(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve user ID from the context
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "Access denied: user not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if the user is an admin
			var isAdmin bool
			err := db.QueryRow("SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin)
			if err != nil || !isAdmin {
				http.Error(w, "Access denied: not an admin", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireApproved middleware ensures the user has an approved account
func RequireApproved(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve user ID from the context
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "Access denied: user not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if the user is approved
			var isApproved bool
			err := db.QueryRow("SELECT is_approved FROM users WHERE id = $1", userID).Scan(&isApproved)
			if err != nil || !isApproved {
				http.Error(w, "Access denied: account not approved", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
