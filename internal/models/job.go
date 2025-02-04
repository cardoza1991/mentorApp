package models

import (
	"time"
)

// Job represents a job posting
type Job struct {
	ID              int       `json:"id"`
	Title           string    `json:"title" validate:"required"`
	Company         string    `json:"company" validate:"required"`
	Location        string    `json:"location"`
	Description     string    `json:"description" validate:"required"`
	Requirements    string    `json:"requirements"`
	SalaryRange     string    `json:"salary_range"`
	JobType         string    `json:"job_type" validate:"required,oneof=full-time part-time contract"`
	ExperienceLevel string    `json:"experience_level" validate:"required,oneof=entry mid senior"`
	RemotePolicy    string    `json:"remote_policy" validate:"required,oneof=remote hybrid on-site"`
	ContactEmail    string    `json:"contact_email,omitempty"`
	Status          string    `json:"status" validate:"required,oneof=active expired closed"`
	IsFeatured      bool      `json:"is_featured"`
	CreatedBy       int       `json:"created_by"` // Adding CreatedBy field
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Job Status constants
var JobStatus = struct {
	Active  string
	Expired string
	Closed  string
}{
	Active:  "active",
	Expired: "expired",
	Closed:  "closed",
}

// Remote Policy constants
var RemotePolicy = struct {
	Remote string
	Hybrid string
	OnSite string
}{
	Remote: "remote",
	Hybrid: "hybrid",
	OnSite: "on-site",
}

// Experience Level constants
var ExperienceLevel = struct {
	Entry  string
	Mid    string
	Senior string
}{
	Entry:  "entry",
	Mid:    "mid",
	Senior: "senior",
}

// Job Type constants
var JobType = struct {
	FullTime string
	PartTime string
	Contract string
}{
	FullTime: "full-time",
	PartTime: "part-time",
	Contract: "contract",
}

// JobApplication represents a user's application for a job
type JobApplication struct {
	ID          int       `json:"id"`
	JobID       int       `json:"job_id"`
	UserID      int       `json:"user_id"`
	Status      string    `json:"status" validate:"required,oneof=pending reviewed accepted rejected"`
	CoverLetter string    `json:"cover_letter"`
	ResumeURL   string    `json:"resume_url,omitempty"`
	AppliedAt   time.Time `json:"applied_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Application Status constants
var ApplicationStatus = struct {
	Pending  string
	Reviewed string
	Accepted string
	Rejected string
}{
	Pending:  "pending",
	Reviewed: "reviewed",
	Accepted: "accepted",
	Rejected: "rejected",
}
