package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mentorApp/internal/models"
	"mentorApp/internal/repository"
)

type MentorshipService struct {
	mentorshipRepo repository.IMentorshipRepository // Changed from *repository.MentorshipRepository
	profileRepo    repository.IProfileRepository    // Use interface instead of concrete type
	userRepo       repository.IUserRepository       // Use interface instead of concrete type
}

// Updated constructor
func NewMentorshipService(
	mentorshipRepo repository.IMentorshipRepository,
	profileRepo repository.IProfileRepository,
	userRepo repository.IUserRepository,
) IMentorshipService {
	return &MentorshipService{
		mentorshipRepo: mentorshipRepo,
		profileRepo:    profileRepo,
		userRepo:       userRepo,
	}
}

// CreateMentorshipProgram creates a new mentorship program
func (s *MentorshipService) CreateMentorshipProgram(ctx context.Context, mentorID int, program *models.MentorshipProgram) error {
	// Verify mentor exists and is approved
	user, err := s.userRepo.GetUserByID(ctx, mentorID)
	if err != nil {
		return err
	}
	if !user.IsMentor || !user.IsApproved {
		return errors.New("unauthorized: user is not an approved mentor")
	}

	program.MentorID = mentorID // Fixed from MentorId
	program.Status = "active"

	return s.mentorshipRepo.CreateProgram(ctx, program)
}

// RequestMentorship handles a mentee's request for mentorship
func (s *MentorshipService) RequestMentorship(ctx context.Context, menteeID, programID int, message string) error {
	// Verify mentee exists
	mentee, err := s.userRepo.GetUserByID(ctx, menteeID)
	if err != nil {
		return err
	}
	if mentee.IsMentor {
		return errors.New("mentors cannot request mentorship")
	}

	// Get program details
	program, err := s.mentorshipRepo.GetProgram(ctx, programID)
	if err != nil {
		return err
	}

	// Create request
	request := &models.MentorshipRequest{
		MentorID:  program.MentorID, // Fixed from MentorId
		MenteeID:  menteeID,         // Fixed from MenteeId
		ProgramID: programID,        // Fixed from ProgramId
		Status:    "pending",
		Message:   message,
	}

	return s.mentorshipRepo.CreateRequest(ctx, request)
}

// RespondToRequest handles a mentor's response to a mentorship request
func (s *MentorshipService) RespondToRequest(ctx context.Context, mentorID, requestID int, approve bool) error {
	// Verify mentor owns the request
	request, err := s.mentorshipRepo.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}
	if request.MentorID != mentorID { // Fixed from MentorId
		return errors.New("unauthorized: request belongs to different mentor")
	}

	status := "rejected"
	if approve {
		status = "approved"
	}

	return s.mentorshipRepo.UpdateRequestStatus(ctx, requestID, status)
}

// ScheduleSession schedules a new mentorship session
func (s *MentorshipService) ScheduleSession(ctx context.Context, userID int, session *models.MentorshipSession) error {
	// Get the request to verify permissions
	request, err := s.mentorshipRepo.GetRequest(ctx, session.RequestID) // Fixed from RequestId
	if err != nil {
		return err
	}

	// Verify user is part of this mentorship
	if request.MentorID != userID && request.MenteeID != userID { // Fixed from MentorId/MenteeId
		return errors.New("unauthorized: user not part of this mentorship")
	}

	// Validate session time
	if session.StartTime.Before(time.Now()) {
		return errors.New("cannot schedule session in the past")
	}
	if session.EndTime.Before(session.StartTime) {
		return errors.New("end time must be after start time")
	}

	maxSessionDuration := 4 * time.Hour
	if session.EndTime.Sub(session.StartTime) > maxSessionDuration {
		return errors.New("session duration exceeds program limits")
	}

	session.Status = "scheduled"
	return s.mentorshipRepo.CreateSession(ctx, session)
}

// ListMentorPrograms returns all programs created by a mentor
func (s *MentorshipService) ListMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error) {
	return s.mentorshipRepo.ListMentorPrograms(ctx, mentorID)
}

// ListMenteeRequests returns all mentorship requests made by a mentee
func (s *MentorshipService) ListMenteeRequests(ctx context.Context, menteeID int) ([]*models.MentorshipRequest, error) {
	return s.mentorshipRepo.ListMenteeRequests(ctx, menteeID)
}

// Add these methods to your MentorshipService struct

// GetFeaturedMentors returns a list of featured mentors
func (s *MentorshipService) GetFeaturedMentors(ctx context.Context) ([]*models.Profile, error) {
	// TODO: Implement actual featured mentors logic
	// For now, return empty slice
	return []*models.Profile{}, nil
}

// GetAvailableJobs returns a list of available jobs
func (s *MentorshipService) GetAvailableJobs(ctx context.Context) ([]*models.Job, error) {
	// TODO: Implement actual job listing logic
	// For now, return empty slice
	return []*models.Job{}, nil
}

// GetActiveMentorships returns active mentorships for a user
func (s *MentorshipService) GetActiveMentorships(ctx context.Context, userID int) ([]*models.MentorshipRequest, error) {
	// TODO: Implement actual active mentorships logic
	// For now, return empty slice
	return []*models.MentorshipRequest{}, nil
}

// GetUpcomingSessions returns upcoming sessions for a user
func (s *MentorshipService) GetUpcomingSessions(ctx context.Context, userID int) ([]*models.MentorshipSession, error) {
	// TODO: Implement actual upcoming sessions logic
	// For now, return empty slice
	return []*models.MentorshipSession{}, nil
}

// GetRecommendedMentors returns recommended mentors for a mentee
func (s *MentorshipService) GetRecommendedMentors(ctx context.Context, userID int) ([]*models.Profile, error) {
	// TODO: Implement actual mentor recommendations logic
	// For now, return empty slice
	return []*models.Profile{}, nil
}

// GetMentorPrograms returns all programs for a mentor
func (s *MentorshipService) GetMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error) {
	// Query programs for this mentor
	programs, err := s.mentorshipRepo.ListMentorPrograms(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("error fetching mentor programs: %v", err)
	}

	// Filter for active programs only
	var activePrograms []*models.MentorshipProgram
	for _, program := range programs {
		if program.Status == "active" {
			activePrograms = append(activePrograms, program)
		}
	}

	return activePrograms, nil
}

// GetPendingRequests returns pending mentorship requests
func (s *MentorshipService) GetPendingRequests(ctx context.Context, mentorID int) ([]*models.MentorshipRequest, error) {
	// TODO: Implement actual pending requests logic
	// For now, return empty slice
	return []*models.MentorshipRequest{}, nil
}

// GetMentorAnalytics returns analytics for a mentor
func (s *MentorshipService) GetMentorAnalytics(ctx context.Context, mentorID int) (map[string]interface{}, error) {
	// TODO: Implement actual analytics logic
	// For now, return empty map
	return map[string]interface{}{
		"total_sessions": 0,
		"total_mentees":  0,
		"average_rating": 0.0,
	}, nil
}

// GetAvailableSpecialties returns list of available specialties
func (s *MentorshipService) GetAvailableSpecialties(ctx context.Context) []string {
	// TODO: Implement actual specialties logic
	// For now, return some default specialties
	return []string{
		"Backend Development",
		"Frontend Development",
		"DevOps",
		"Cloud Architecture",
		"Cloud Native Technologies",
		"Mobile Development",
		"Data Science",
		"AI",
		"NLP",
		"ML",
	}
}

// ListAvailablePrograms returns all active mentorship programs
func (s *MentorshipService) ListAvailablePrograms(ctx context.Context) ([]*models.MentorshipProgram, error) {
	// Implement by filtering for active programs only
	programs, err := s.mentorshipRepo.ListMentorPrograms(ctx, 0) // 0 means all mentors
	if err != nil {
		return nil, err
	}

	var activePrograms []*models.MentorshipProgram
	for _, program := range programs {
		if program.Status == "active" {
			activePrograms = append(activePrograms, program)
		}
	}
	return activePrograms, nil
}

// GetProgramDetails returns detailed information about a specific program
func (s *MentorshipService) GetProgramDetails(ctx context.Context, programID int) (*models.MentorshipProgram, error) {
	return s.mentorshipRepo.GetProgram(ctx, programID)
}

// SearchMentors searches for mentors based on specified filters
func (s *MentorshipService) SearchMentors(ctx context.Context, filters map[string]interface{}) ([]*models.PublicProfile, error) {
	profiles, err := s.profileRepo.SearchProfiles(ctx, filters)
	if err != nil {
		return nil, err
	}

	var results []*models.PublicProfile
	for _, profile := range profiles {
		results = append(results, &models.PublicProfile{
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Bio:       profile.Bio,
			Skills:    profile.Skills,
			Rate:      profile.Rate,
		})
	}
	return results, nil
}

// UpdateAvailability updates a mentor's availability schedule
func (s *MentorshipService) UpdateAvailability(ctx context.Context, mentorID int, availability []models.Availability) error {
	// First verify the user is a mentor
	user, err := s.userRepo.GetUserByID(ctx, mentorID)
	if err != nil {
		return err
	}
	if !user.IsMentor {
		return errors.New("user is not a mentor")
	}

	// Update availability
	for _, slot := range availability {
		slot.ProfileID = mentorID // Ensure profile ID is set correctly
		if err := s.profileRepo.SetAvailability(ctx, &slot); err != nil {
			return err
		}
	}
	return nil
}

// SubmitSessionFeedback submits feedback for a mentorship session
func (s *MentorshipService) SubmitSessionFeedback(ctx context.Context, feedback *models.SessionFeedback) error {
	// Implement feedback submission logic
	return s.mentorshipRepo.CreateSessionFeedback(ctx, feedback)
}

// GetAvailability retrieves a mentor's availability schedule
func (s *MentorshipService) GetAvailability(ctx context.Context, mentorID int) ([]models.Availability, error) {
	return s.profileRepo.GetProfileAvailability(ctx, mentorID)
}

// GetMentorshipStats returns statistics about a mentor's programs and sessions
func (s *MentorshipService) GetMentorshipStats(ctx context.Context, mentorID int) (map[string]interface{}, error) {
	// This could be optimized with specific queries rather than loading all data
	programs, err := s.mentorshipRepo.ListMentorPrograms(ctx, mentorID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_programs":     len(programs),
		"active_programs":    0,
		"total_mentees":      0,
		"completed_sessions": 0,
	}

	for _, program := range programs {
		if program.Status == "active" {
			stats["active_programs"] = stats["active_programs"].(int) + 1
		}
		// Additional stats could be calculated here
	}

	return stats, nil
}
