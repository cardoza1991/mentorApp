package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"mentorApp/internal/models"
	"mentorApp/internal/repository"
	"mentorApp/pkg/utils/email"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	emailSvc    *email.EmailService
}

func NewUserService(
	userRepo *repository.UserRepository,
	profileRepo *repository.ProfileRepository,
	emailSvc *email.EmailService,
) IUserService {
	return &UserService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		emailSvc:    emailSvc,
	}
}

// RegisterUser handles the complete user registration process
func (s *UserService) RegisterUser(ctx context.Context, input RegisterUserInput) (*models.User, error) {
	// Check if email is already registered
	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	if err := validateEmailDomain(input.Email); err != nil {
		return nil, err
	}

	// Create user with verification token
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		IsMentor: input.IsMentor,
		VerificationToken: sql.NullString{
			String: uuid.New().String(),
			Valid:  true,
		},
	}

	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}

	// Create user and profile in transaction
	err = s.userRepo.CreateUserWithProfile(ctx, user, &models.Profile{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Bio:       input.Bio,
		Skills:    input.Skills,
	})
	if err != nil {
		return nil, err
	}

	// Send verification email
	if user.VerificationToken.Valid {
		go s.emailSvc.SendVerificationEmail(user.Email, user.VerificationToken.String)
	} else {
		log.Printf("Verification token is invalid for user: %s", user.Email)
	}

	return user, nil
}

// RegisterMentor handles mentor-specific registration
func (s *UserService) RegisterMentor(ctx context.Context, input RegisterMentorInput) (*models.User, error) {
	// Check if email already registered
	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	if err := validateEmailDomain(input.Email); err != nil {
		return nil, err
	}

	// Create user - note we're setting EmailVerified to true to skip verification
	// but IsApproved to false to require admin approval
	user := &models.User{
		Username:      input.Username,
		Email:         input.Email,
		IsMentor:      true,
		EmailVerified: true,  // Skip email verification
		IsApproved:    false, // Require admin approval
	}

	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}

	profile := &models.Profile{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Bio:       input.Bio,
		Skills:    input.Skills,
		Rate:      input.Rate,
	}

	// Create user and profile
	err = s.userRepo.CreateUserWithProfile(ctx, user, profile)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	log.Printf("Attempting to authenticate user with email: %s", email)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		log.Printf("No user found with email: %s", email)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("Found user: ID=%d, Email=%s, IsAdmin=%v", user.Id, user.Email, user.IsAdmin) // Avoid logging sensitive info

	valid := user.ValidatePassword(password)
	log.Printf("Password validation result: %v", valid)

	if !valid {
		return nil, errors.New("invalid credentials")
	}

	// Skip email verification for admin users
	if !user.IsAdmin && !user.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return user, nil
}

// VerifyEmail verifies a user's email address
func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.GetUserByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid verification token")
	}

	user.EmailVerified = true
	user.VerificationToken = sql.NullString{
		String: "",
		Valid:  false,
	}
	return s.userRepo.UpdateUser(ctx, user)
}

// RequestPasswordReset initiates the password reset process
func (s *UserService) RequestPasswordReset(ctx context.Context, emailAddr string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	resetToken := uuid.New().String()
	user.ResetToken = sql.NullString{
		String: resetToken,
		Valid:  true,
	}
	expiry := time.Now().Add(24 * time.Hour)
	user.ResetTokenExpiry = &expiry

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	if user.ResetToken.Valid && user.ResetTokenExpiry != nil {
		if err := s.emailSvc.SendPasswordResetEmail(user.Email, user.ResetToken.String); err != nil {
			log.Printf("Failed to send password reset email: %v", err)
		}
	} else {
		log.Printf("Reset token is invalid for user: %s", user.Email)
	}

	return nil
}

// ResetPassword completes the password reset process
func (s *UserService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.userRepo.GetUserByResetToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil || user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		return errors.New("invalid or expired reset token")
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	user.ResetToken = sql.NullString{
		String: "",
		Valid:  false,
	}
	user.ResetTokenExpiry = nil

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) GetPublicProfile(ctx context.Context, userID int) (*models.PublicProfile, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	publicProfile := &models.PublicProfile{
		FirstName:      profile.FirstName,
		LastName:       profile.LastName,
		Bio:            profile.Bio,
		Skills:         profile.Skills,
		Rate:           profile.Rate,
		ProfilePicture: profile.ProfilePicture,
	}

	return publicProfile, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, userID int, profile *models.Profile) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	profile.UserId = userID
	return s.profileRepo.UpdateProfile(ctx, profile)
}

// GetUserProfile retrieves a user's complete profile
func (s *UserService) GetUserProfile(ctx context.Context, userID int) (*models.Profile, error) {
	return s.profileRepo.GetProfileByUserID(ctx, userID)
}

// Settings Management

func (s *UserService) GetProfileSettings(ctx context.Context, userID int) (*models.ProfileSettings, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	return &models.ProfileSettings{
		NotificationPreferences: profile.NotificationPreferences,
		PrivacySettings:         profile.PrivacySettings,
	}, nil
}

func (s *UserService) UpdateProfileSettings(ctx context.Context, userID int, settings *models.ProfileSettings) error {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	profile.NotificationPreferences = settings.NotificationPreferences
	profile.PrivacySettings = settings.PrivacySettings

	return s.profileRepo.UpdateProfile(ctx, profile)
}

func (s *UserService) GetNotificationSettings(ctx context.Context, userID int) (*models.NotificationSettings, error) {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// Assuming NotificationPreferences is a JSON string
	var ns models.NotificationSettings
	if profile.NotificationPreferences != "" {
		if err := json.Unmarshal([]byte(profile.NotificationPreferences), &ns); err != nil {
			return nil, err
		}
	}

	return &ns, nil
}

func (s *UserService) UpdateNotificationSettings(ctx context.Context, userID int, settings *models.NotificationSettings) error {
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	profile.NotificationPreferences = string(data)

	return s.profileRepo.UpdateProfile(ctx, profile)
}

// CreateSession is a placeholder; implement session storage logic as needed.
func (s *UserService) CreateSession(ctx context.Context, userID int) (string, error) {
	return fmt.Sprintf("session_%d_%s", userID, uuid.New().String()), nil
}

func validateEmailDomain(email string) error {
	// Already validated in User model, but we can keep consistent checks
	return nil
}
