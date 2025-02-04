package models

import (
	"time"
)

type MentorshipProgram struct {
	ID          int       `json:"id"`
	MentorID    int       `json:"mentor_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	Price       float64   `json:"price"`
	MaxMentees  int       `json:"max_mentees"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MentorshipRequest struct {
	ID        int       `json:"id"`
	MentorID  int       `json:"mentor_id"`
	MenteeID  int       `json:"mentee_id"`
	ProgramID int       `json:"program_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MentorshipSession struct {
	ID        int       `json:"id"`
	RequestID int       `json:"request_id"`
	Title     string    `json:"title"` // Added Title field
	Topic     string    `json:"topic"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Add these session status constants
var SessionStatus = struct {
	Scheduled string
	Completed string
	Cancelled string
	NoShow    string
}{
	Scheduled: "scheduled",
	Completed: "completed",
	Cancelled: "cancelled",
	NoShow:    "no_show",
}
