package services

import (
	"context"
	"errors"
	"time"

	"mentorApp/internal/repository"

	"mentorApp/internal/models"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// CreateProfile creates a new profile for a user
func (s *ProfileService) CreateProfile(ctx context.Context, userID int, profile *models.Profile) error {
	// Verify user exists
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	profile.UserId = userID
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	return s.profileRepo.CreateProfile(ctx, profile)
}

// GetProfile retrieves a user's profile
func (s *ProfileService) GetProfile(ctx context.Context, userID int) (*models.Profile, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}

// UpdateProfile updates a user's profile information
func (s *ProfileService) UpdateProfile(ctx context.Context, userID int, updates *models.ProfileUpdate) error {
	// Verify profile exists
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	// Apply updates only if the pointer is not nil
	if updates.FirstName != nil {
		profile.FirstName = *updates.FirstName
	}
	if updates.LastName != nil {
		profile.LastName = *updates.LastName
	}
	if updates.Bio != nil {
		profile.Bio = *updates.Bio
	}
	if updates.Skills != nil {
		profile.Skills = *updates.Skills
	}
	if updates.TimeZone != nil {
		profile.Timezone = *updates.TimeZone
	}
	if updates.Rate != nil {
		profile.Rate = *updates.Rate
	}
	if updates.Availability != nil {
		// Handle availability update
	}
	if updates.ProfilePicture != nil {
		// Handle profile picture update
	}

	profile.UpdatedAt = time.Now()

	return s.profileRepo.UpdateProfile(ctx, profile)
}

// GetPublicProfile retrieves a user's public profile information
func (s *ProfileService) GetPublicProfile(ctx context.Context, userID int) (*models.PublicProfile, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// Get user's ratings and reviews
	ratings, err := s.profileRepo.GetUserRatings(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate average rating
	var averageRating float64
	if len(ratings) > 0 {
		var total float64
		for _, r := range ratings {
			total += r.Rating
		}
		averageRating = total / float64(len(ratings))
	}

	return &models.PublicProfile{
		FirstName:      profile.FirstName,
		LastName:       profile.LastName,
		Bio:            profile.Bio,
		Skills:         profile.Skills,
		Rate:           profile.Rate,
		ProfilePicture: profile.ProfilePicture,
		AverageRating:  averageRating,
		RatingCount:    len(ratings),
	}, nil
}

// UpdateProfileSettings updates a user's profile settings
func (s *ProfileService) UpdateProfileSettings(ctx context.Context, userID int, settings *models.ProfileSettings) error {
	// Verify profile exists
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	profile.NotificationPreferences = settings.NotificationPreferences
	profile.PrivacySettings = settings.PrivacySettings
	profile.UpdatedAt = time.Now()

	return s.profileRepo.UpdateProfile(ctx, profile)
}

// AddProfileExperience adds a new experience entry to a profile
func (s *ProfileService) AddProfileExperience(ctx context.Context, userID int, experience *models.Experience) error {
	// Verify profile exists
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	experience.UserID = userID
	experience.CreatedAt = time.Now()

	return s.profileRepo.AddExperience(ctx, experience)
}

// UpdateProfileExperience updates an existing experience entry
func (s *ProfileService) UpdateProfileExperience(ctx context.Context, userID int, experienceID int, updates *models.Experience) error {
	// Verify ownership
	existing, err := s.profileRepo.GetExperience(ctx, experienceID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return errors.New("experience not found or unauthorized")
	}

	updates.ID = experienceID
	updates.UserID = userID
	updates.UpdatedAt = time.Now()

	return s.profileRepo.UpdateExperience(ctx, updates)
}

// DeleteProfileExperience removes an experience entry
func (s *ProfileService) DeleteProfileExperience(ctx context.Context, userID int, experienceID int) error {
	// Verify ownership
	existing, err := s.profileRepo.GetExperience(ctx, experienceID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return errors.New("experience not found or unauthorized")
	}

	return s.profileRepo.DeleteExperience(ctx, experienceID)
}

// AddProfileEducation adds a new education entry to a profile
func (s *ProfileService) AddProfileEducation(ctx context.Context, userID int, education *models.Education) error {
	// Verify profile exists
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	education.UserID = userID
	education.CreatedAt = time.Now()

	return s.profileRepo.AddEducation(ctx, education)
}

// UpdateAvailability updates a user's availability schedule
func (s *ProfileService) UpdateAvailability(ctx context.Context, userID int, availability *models.Availability) error {
	// Verify profile exists
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	availability.UserID = userID
	availability.UpdatedAt = time.Now()

	return s.profileRepo.UpdateAvailability(ctx, availability)
}
