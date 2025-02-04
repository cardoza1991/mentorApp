package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mentorApp/internal/models"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

// CreateProfile creates a new profile
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	query := `INSERT INTO profiles (
        user_id, first_name, last_name, bio, skills, experience, linkedin,
        github, twitter, rate, available, timezone, profile_picture,
        notification_preferences, privacy_settings
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
    ) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		profile.UserId,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.Skills,
		profile.Experience,
		profile.LinkedIn,
		profile.Github,
		profile.Twitter,
		profile.Rate,
		profile.Available,
		profile.Timezone,
		profile.ProfilePicture,
		profile.NotificationPreferences,
		profile.PrivacySettings,
	).Scan(&profile.Id)
}

// GetProfileByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	query := `SELECT id, user_id, first_name, last_name, bio, skills, experience,
        linkedin, github, twitter, rate, available, timezone, profile_picture,
        notification_preferences, privacy_settings, created_at, updated_at
        FROM profiles WHERE user_id = $1`

	profile := &models.Profile{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.Id,
		&profile.UserId,
		&profile.FirstName,
		&profile.LastName,
		&profile.Bio,
		&profile.Skills,
		&profile.Experience,
		&profile.LinkedIn,
		&profile.Github,
		&profile.Twitter,
		&profile.Rate,
		&profile.Available,
		&profile.Timezone,
		&profile.ProfilePicture,
		&profile.NotificationPreferences,
		&profile.PrivacySettings,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return profile, err
}

// UpdateProfile updates profile information
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	query := `UPDATE profiles SET
        first_name = $1,
        last_name = $2,
        bio = $3,
        skills = $4,
        experience = $5,
        linkedin = $6,
        github = $7,
        twitter = $8,
        rate = $9,
        available = $10,
        timezone = $11,
        profile_picture = $12,
        notification_preferences = $13,
        privacy_settings = $14,
        updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $15`

	result, err := r.db.ExecContext(ctx, query,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.Skills,
		profile.Experience,
		profile.LinkedIn,
		profile.Github,
		profile.Twitter,
		profile.Rate,
		profile.Available,
		profile.Timezone,
		profile.ProfilePicture,
		profile.NotificationPreferences,
		profile.PrivacySettings,
		profile.UserId,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetUserRatings gets all ratings for a user
func (r *ProfileRepository) GetUserRatings(ctx context.Context, userID int) ([]*models.Rating, error) {
	query := `SELECT id, from_user_id, rating, comment, created_at 
              FROM ratings WHERE to_user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*models.Rating
	for rows.Next() {
		rating := &models.Rating{}
		err := rows.Scan(
			&rating.Id,
			&rating.FromUserId,
			&rating.Rating,
			&rating.Comment,
			&rating.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

// AddExperience adds a new experience entry
func (r *ProfileRepository) AddExperience(ctx context.Context, experience *models.Experience) error {
	query := `INSERT INTO experience (user_id, title, company, start_date, end_date, description, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		experience.UserID,
		experience.Title,
		experience.Company,
		experience.StartDate,
		experience.EndDate,
		experience.Description,
	).Scan(&experience.ID)
}

// GetExperience gets a specific experience entry
func (r *ProfileRepository) GetExperience(ctx context.Context, experienceID int) (*models.Experience, error) {
	query := `SELECT id, user_id, title, company, start_date, end_date, description, created_at, updated_at
              FROM experience WHERE id = $1`

	exp := &models.Experience{}
	err := r.db.QueryRowContext(ctx, query, experienceID).Scan(
		&exp.ID,
		&exp.UserID,
		&exp.Title,
		&exp.Company,
		&exp.StartDate,
		&exp.EndDate,
		&exp.Description,
		&exp.CreatedAt,
		&exp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return exp, err
}

func (r *ProfileRepository) UpdateExperience(ctx context.Context, experience *models.Experience) error {
	query := `UPDATE experience 
              SET title = $1, company = $2, start_date = $3, end_date = $4, description = $5, updated_at = CURRENT_TIMESTAMP
              WHERE id = $6 AND user_id = $7`

	result, err := r.db.ExecContext(ctx, query,
		experience.Title,
		experience.Company,
		experience.StartDate,
		experience.EndDate,
		experience.Description,
		experience.ID,
		experience.UserID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteExperience removes an experience entry
func (r *ProfileRepository) DeleteExperience(ctx context.Context, experienceID int) error {
	query := `DELETE FROM experience WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, experienceID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// AddEducation adds a new education entry
func (r *ProfileRepository) AddEducation(ctx context.Context, education *models.Education) error {
	query := `INSERT INTO education (user_id, degree, institution, start_date, end_date)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		education.UserID,
		education.Degree,
		education.Institution,
		education.StartDate,
		education.EndDate).Scan(&education.ID)
}

// SearchProfiles searches for profiles based on filters
func (r *ProfileRepository) SearchProfiles(ctx context.Context, filters map[string]interface{}) ([]*models.Profile, error) {
	query := `SELECT p.* FROM profiles p
              INNER JOIN users u ON p.user_id = u.id
              WHERE u.is_mentor = true AND u.is_approved = true`

	var args []interface{}
	argCount := 1

	if skills, ok := filters["skills"].(string); ok && skills != "" {
		query += fmt.Sprintf(` AND p.skills ILIKE $%d`, argCount)
		args = append(args, "%"+skills+"%")
		argCount++
	}

	if minRate, ok := filters["rate_min"].(float64); ok {
		query += fmt.Sprintf(` AND p.rate >= $%d`, argCount)
		args = append(args, minRate)
		argCount++
	}

	if maxRate, ok := filters["rate_max"].(float64); ok {
		query += fmt.Sprintf(` AND p.rate <= $%d`, argCount)
		args = append(args, maxRate)
		argCount++
	}

	if timezone, ok := filters["timezone"].(string); ok && timezone != "" {
		query += fmt.Sprintf(` AND p.timezone = $%d`, argCount)
		args = append(args, timezone)
		argCount++
	}

	query += ` ORDER BY p.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*models.Profile
	for rows.Next() {
		profile := &models.Profile{}
		if err := rows.Scan(
			&profile.Id,
			&profile.UserId,
			&profile.FirstName,
			&profile.LastName,
			&profile.Bio,
			&profile.Skills,
			&profile.Rate,
			&profile.Available,
			&profile.Timezone,
			&profile.ProfilePicture,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// SetAvailability sets availability for a profile
func (r *ProfileRepository) SetAvailability(ctx context.Context, availability *models.Availability) error {
	query := `INSERT INTO availability (profile_id, day_of_week, start_time, end_time)
              VALUES ($1, $2, $3, $4)
              ON CONFLICT (profile_id, day_of_week) DO UPDATE
              SET start_time = EXCLUDED.start_time, end_time = EXCLUDED.end_time`

	_, err := r.db.ExecContext(ctx, query,
		availability.ProfileID,
		availability.DayOfWeek,
		availability.StartTime,
		availability.EndTime,
	)
	return err
}

// GetProfileAvailability gets all availability slots for a profile
func (r *ProfileRepository) GetProfileAvailability(ctx context.Context, profileID int) ([]models.Availability, error) {
	query := `SELECT profile_id, day_of_week, start_time, end_time, created_at, updated_at
              FROM availability WHERE profile_id = $1 ORDER BY day_of_week, start_time`

	rows, err := r.db.QueryContext(ctx, query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availability []models.Availability
	for rows.Next() {
		slot := models.Availability{}
		if err := rows.Scan(
			&slot.ProfileID,
			&slot.DayOfWeek,
			&slot.StartTime,
			&slot.EndTime,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		); err != nil {
			return nil, err
		}
		availability = append(availability, slot)
	}
	return availability, nil
}

// UpdateAvailability updates a user's availability
func (r *ProfileRepository) UpdateAvailability(ctx context.Context, availability *models.Availability) error {
	query := `UPDATE availability 
              SET day_of_week = $1, start_time = $2, end_time = $3, updated_at = CURRENT_TIMESTAMP
              WHERE profile_id = $4`

	result, err := r.db.ExecContext(ctx, query,
		availability.DayOfWeek,
		availability.StartTime,
		availability.EndTime,
		availability.ProfileID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
