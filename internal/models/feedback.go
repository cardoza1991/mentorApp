package models

import (
	"time"
)

// SessionFeedback represents feedback for a mentorship session
type SessionFeedback struct {
	ID        int       `json:"id"`
	SessionID int       `json:"session_id"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// FeedbackResponse represents the response shown to users
type FeedbackResponse struct {
	ID        int       `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ValidationRules returns the validation rules for feedback
func (f *SessionFeedback) ValidationRules() []string {
	return []string{
		"Rating is required and must be between 1 and 5",
	}
}
