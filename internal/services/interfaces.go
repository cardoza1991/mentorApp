// File: internal/services/interfaces.go

package services

import (
	"context"
	"mentorApp/internal/models"
)

// IMentorshipService defines the interface for mentorship-related operations
type IMentorshipService interface {
	// Program Management
	CreateMentorshipProgram(ctx context.Context, mentorID int, program *models.MentorshipProgram) error
	GetMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error)
	GetProgramDetails(ctx context.Context, programID int) (*models.MentorshipProgram, error)
	ListMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error)
	ListAvailablePrograms(ctx context.Context) ([]*models.MentorshipProgram, error)

	// Request Management
	RequestMentorship(ctx context.Context, menteeID, programID int, message string) error
	RespondToRequest(ctx context.Context, mentorID, requestID int, approve bool) error
	ListMenteeRequests(ctx context.Context, menteeID int) ([]*models.MentorshipRequest, error)
	GetPendingRequests(ctx context.Context, mentorID int) ([]*models.MentorshipRequest, error)

	// Session Management
	ScheduleSession(ctx context.Context, userID int, session *models.MentorshipSession) error
	GetUpcomingSessions(ctx context.Context, userID int) ([]*models.MentorshipSession, error)
	SubmitSessionFeedback(ctx context.Context, feedback *models.SessionFeedback) error

	// Search and Discovery
	SearchMentors(ctx context.Context, filters map[string]interface{}) ([]*models.PublicProfile, error)
	GetFeaturedMentors(ctx context.Context) ([]*models.Profile, error)
	GetRecommendedMentors(ctx context.Context, userID int) ([]*models.Profile, error)

	// Analytics and Stats
	GetMentorshipStats(ctx context.Context, mentorID int) (map[string]interface{}, error)
	GetMentorAnalytics(ctx context.Context, mentorID int) (map[string]interface{}, error)

	// Availability Management
	UpdateAvailability(ctx context.Context, mentorID int, availability []models.Availability) error
	GetAvailability(ctx context.Context, mentorID int) ([]models.Availability, error)

	// Job Board
	GetAvailableJobs(ctx context.Context) ([]*models.Job, error)
	GetActiveMentorships(ctx context.Context, userID int) ([]*models.MentorshipRequest, error)
	GetAvailableSpecialties(ctx context.Context) []string
}

// IProfileService defines the interface for profile-related operations
type IProfileService interface {
	CreateProfile(ctx context.Context, userID int, profile *models.Profile) error
	GetProfile(ctx context.Context, userID int) (*models.Profile, error)
	UpdateProfile(ctx context.Context, userID int, updates *models.ProfileUpdate) error
	GetPublicProfile(ctx context.Context, userID int) (*models.PublicProfile, error)
	UpdateProfileSettings(ctx context.Context, userID int, settings *models.ProfileSettings) error
	AddProfileExperience(ctx context.Context, userID int, experience *models.Experience) error
	UpdateProfileExperience(ctx context.Context, userID int, experienceID int, updates *models.Experience) error
	DeleteProfileExperience(ctx context.Context, userID int, experienceID int) error
}

// IUserService defines the interface for user-related operations
type IUserService interface {
	// Existing methods
	RegisterUser(ctx context.Context, input RegisterUserInput) (*models.User, error)
	RegisterMentor(ctx context.Context, input RegisterMentorInput) (*models.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	VerifyEmail(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	// Added missing methods
	GetUserProfile(ctx context.Context, userID int) (*models.Profile, error)
	UpdateUserProfile(ctx context.Context, userID int, profile *models.Profile) error
	GetPublicProfile(ctx context.Context, userID int) (*models.PublicProfile, error)
	GetProfileSettings(ctx context.Context, userID int) (*models.ProfileSettings, error)
	UpdateProfileSettings(ctx context.Context, userID int, settings *models.ProfileSettings) error
	GetNotificationSettings(ctx context.Context, userID int) (*models.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, userID int, settings *models.NotificationSettings) error
	CreateSession(ctx context.Context, userID int) (string, error)
}

// RegisterUserInput represents the input for user registration
type RegisterUserInput struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Bio       string
	Skills    string
	IsMentor  bool
}

// RegisterMentorInput represents the input for mentor registration
type RegisterMentorInput struct {
	Username    string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Bio         string
	Skills      string
	Rate        float64
	Experience  string
	Specialties string
}
