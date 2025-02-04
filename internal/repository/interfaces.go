// File: internal/repository/interfaces.go

package repository

import (
	"context"
	"mentorApp/internal/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByVerificationToken(ctx context.Context, token string) (*models.User, error)
	GetUserByResetToken(ctx context.Context, token string) (*models.User, error)
	CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile) error
	GetAllUsers(ctx context.Context) ([]*models.User, error) // Add this line
}

type IProfileRepository interface {
	CreateProfile(ctx context.Context, profile *models.Profile) error
	GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error)
	UpdateProfile(ctx context.Context, profile *models.Profile) error
	GetUserRatings(ctx context.Context, userID int) ([]*models.Rating, error)
	AddExperience(ctx context.Context, experience *models.Experience) error
	GetExperience(ctx context.Context, experienceID int) (*models.Experience, error)
	UpdateExperience(ctx context.Context, experience *models.Experience) error
	DeleteExperience(ctx context.Context, experienceID int) error
	SearchProfiles(ctx context.Context, filters map[string]interface{}) ([]*models.Profile, error)
	SetAvailability(ctx context.Context, availability *models.Availability) error
	GetProfileAvailability(ctx context.Context, profileID int) ([]models.Availability, error)
	UpdateAvailability(ctx context.Context, availability *models.Availability) error
}

type IMentorshipRepository interface {
	// Program Management
	CreateProgram(ctx context.Context, program *models.MentorshipProgram) error
	GetProgram(ctx context.Context, id int) (*models.MentorshipProgram, error)
	ListMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error)
	ListActivePrograms(ctx context.Context) ([]*models.MentorshipProgram, error)
	UpdateProgramStatus(ctx context.Context, programID int, status string) error

	// Request Management
	CreateRequest(ctx context.Context, request *models.MentorshipRequest) error
	GetRequest(ctx context.Context, requestID int) (*models.MentorshipRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID int, status string) error
	ListMenteeRequests(ctx context.Context, menteeID int) ([]*models.MentorshipRequest, error)

	// Session Management
	CreateSession(ctx context.Context, session *models.MentorshipSession) error
	GetSession(ctx context.Context, sessionID int) (*models.MentorshipSession, error)
	UpdateSessionStatus(ctx context.Context, sessionID int, status string) error
	ListSessionsByRequest(ctx context.Context, requestID int) ([]*models.MentorshipSession, error)

	// Feedback Management
	CreateSessionFeedback(ctx context.Context, feedback *models.SessionFeedback) error
	GetSessionFeedback(ctx context.Context, sessionID int) ([]*models.SessionFeedback, error)
	GetAverageRating(ctx context.Context, mentorID int) (float64, error)
}

type IJobRepository interface {
	CreateJob(ctx context.Context, job *models.Job) error
	GetJob(ctx context.Context, jobID int) (*models.Job, error)
	ListJobs(ctx context.Context, filters map[string]interface{}) ([]*models.Job, error)
	UpdateJob(ctx context.Context, job *models.Job) error
	DeleteJob(ctx context.Context, jobID int) error
	CreateJobApplication(ctx context.Context, application *models.JobApplication) error
	GetJobApplication(ctx context.Context, applicationID int) (*models.JobApplication, error)
	ListJobApplications(ctx context.Context, jobID int) ([]*models.JobApplication, error)
}
