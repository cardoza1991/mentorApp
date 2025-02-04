package repository

import (
	"context"
	"database/sql"
	"log"
	"mentorApp/internal/models"
)

type UserRepository struct {
	db *sql.DB // Add this field
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, is_mentor, verification_token)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsMentor,
		user.VerificationToken).Scan(&user.Id)
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
        SELECT id, username, email, password_hash, role, is_mentor, is_admin, is_approved, 
               email_verified, verification_token, reset_token, reset_token_expiry, 
               created_at, updated_at
        FROM users 
        WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsMentor,
		&user.IsAdmin,
		&user.IsApproved,
		&user.EmailVerified,
		&user.VerificationToken,
		&user.ResetToken,
		&user.ResetTokenExpiry,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	log.Printf("GetUserByEmail result: ID=%d, Email=%s, IsAdmin=%v",
		user.Id, user.Email, user.IsAdmin)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, role, is_mentor, is_approved, email_verified,
              verification_token, reset_token, reset_token_expiry, created_at, updated_at
              FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsMentor,
		&user.IsApproved,
		&user.EmailVerified,
		&user.VerificationToken,
		&user.ResetToken,
		&user.ResetTokenExpiry,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `
        SELECT id, username, email, role, is_mentor, is_approved, email_verified,
               last_login_at, created_at, updated_at
        FROM users
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.IsMentor,
			&user.IsApproved,
			&user.EmailVerified,
			&lastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullTime to *time.Time
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		} else {
			user.LastLoginAt = nil
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET
              username = $1,
              email = $2,
              password_hash = $3,
              role = $4,
              is_mentor = $5,
              is_approved = $6,
              email_verified = $7,
              verification_token = $8,
              reset_token = $9,
              reset_token_expiry = $10,
              updated_at = CURRENT_TIMESTAMP
              WHERE id = $11`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsMentor,
		user.IsApproved,
		user.EmailVerified,
		user.VerificationToken,
		user.ResetToken,
		user.ResetTokenExpiry,
		user.Id)

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

// GetUserByVerificationToken retrieves a user by their verification token
func (r *UserRepository) GetUserByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, role, is_mentor, is_approved, email_verified
              FROM users WHERE verification_token = $1`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsMentor,
		&user.IsApproved,
		&user.EmailVerified)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByResetToken retrieves a user by their password reset token
func (r *UserRepository) GetUserByResetToken(ctx context.Context, token string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, reset_token_expiry
              FROM users WHERE reset_token = $1`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.ResetTokenExpiry)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// CreateUserWithProfile creates a new user and their profile in a transaction
func (r *UserRepository) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert user
	query := `INSERT INTO users (username, email, password_hash, role, is_mentor, verification_token)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsMentor,
		user.VerificationToken).Scan(&user.Id)

	if err != nil {
		return err
	}

	// Insert profile
	profile.UserId = user.Id
	query = `INSERT INTO profiles (user_id, first_name, last_name, bio, skills)
             VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, query,
		profile.UserId,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.Skills)

	if err != nil {
		return err
	}

	return tx.Commit()
}
