package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mentorApp/internal/models"
)

var (
	ErrProgramNotFound = errors.New("mentorship program not found")
	ErrProgramFull     = errors.New("mentorship program is full")
	ErrInvalidStatus   = errors.New("invalid status transition")
)

type MentorshipRepository struct {
	db *sql.DB
}

func NewMentorshipRepository(db *sql.DB) *MentorshipRepository {
	return &MentorshipRepository{
		db: db,
	}
}

// CreateProgram creates a new mentorship program
func (r *MentorshipRepository) CreateProgram(ctx context.Context, program *models.MentorshipProgram) error {
	query := `
        INSERT INTO mentorship_programs (mentor_id, title, description, duration, price, max_mentees, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		program.MentorID,
		program.Title,
		program.Description,
		program.Duration,
		program.Price,
		program.MaxMentees,
		program.Status,
	).Scan(&program.ID, &program.CreatedAt, &program.UpdatedAt)
}

// GetProgram retrieves a mentorship program by ID
func (r *MentorshipRepository) GetProgram(ctx context.Context, id int) (*models.MentorshipProgram, error) {
	program := &models.MentorshipProgram{}
	query := `
        SELECT id, mentor_id, title, description, duration, price, max_mentees, status, created_at, updated_at
        FROM mentorship_programs
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID,
		&program.MentorID,
		&program.Title,
		&program.Description,
		&program.Duration,
		&program.Price,
		&program.MaxMentees,
		&program.Status,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrProgramNotFound
	}
	return program, err
}

// ListMentorPrograms lists all programs by a mentor
func (r *MentorshipRepository) ListMentorPrograms(ctx context.Context, mentorID int) ([]*models.MentorshipProgram, error) {
	query := `
        SELECT id, mentor_id, title, description, duration, price, max_mentees, status, created_at, updated_at
        FROM mentorship_programs
        WHERE mentor_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []*models.MentorshipProgram
	for rows.Next() {
		program := &models.MentorshipProgram{}
		if err := rows.Scan(
			&program.ID,
			&program.MentorID,
			&program.Title,
			&program.Description,
			&program.Duration,
			&program.Price,
			&program.MaxMentees,
			&program.Status,
			&program.CreatedAt,
			&program.UpdatedAt,
		); err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}
	return programs, nil
}

// CreateRequest creates a new mentorship request
func (r *MentorshipRepository) CreateRequest(ctx context.Context, request *models.MentorshipRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var maxMentees int
	err = tx.QueryRowContext(ctx, "SELECT max_mentees FROM mentorship_programs WHERE id = $1", request.ProgramID).Scan(&maxMentees)
	if err == sql.ErrNoRows {
		return ErrProgramNotFound
	}
	if err != nil {
		return err
	}

	// Create request
	query := `
        INSERT INTO mentorship_requests (mentor_id, mentee_id, program_id, status, message)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		request.MentorID,
		request.MenteeID,
		request.ProgramID,
		request.Status,
		request.Message,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetRequest retrieves a request by ID
func (r *MentorshipRepository) GetRequest(ctx context.Context, requestID int) (*models.MentorshipRequest, error) {
	request := &models.MentorshipRequest{}
	query := `
        SELECT id, mentor_id, mentee_id, program_id, status, message, created_at, updated_at
        FROM mentorship_requests
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&request.ID,
		&request.MentorID,
		&request.MenteeID,
		&request.ProgramID,
		&request.Status,
		&request.Message,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return request, err
}

// ListMenteeRequests lists all requests by a mentee
func (r *MentorshipRepository) ListMenteeRequests(ctx context.Context, menteeID int) ([]*models.MentorshipRequest, error) {
	query := `
        SELECT id, mentor_id, mentee_id, program_id, status, message, created_at, updated_at
        FROM mentorship_requests
        WHERE mentee_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, menteeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.MentorshipRequest
	for rows.Next() {
		request := &models.MentorshipRequest{}
		if err := rows.Scan(
			&request.ID,
			&request.MentorID,
			&request.MenteeID,
			&request.ProgramID,
			&request.Status,
			&request.Message,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

// ListActivePrograms lists all active mentorship programs
func (r *MentorshipRepository) ListActivePrograms(ctx context.Context) ([]*models.MentorshipProgram, error) {
	query := `
        SELECT id, mentor_id, title, description, duration, price, max_mentees, status, created_at, updated_at
        FROM mentorship_programs 
        WHERE status = 'active' 
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []*models.MentorshipProgram
	for rows.Next() {
		program := &models.MentorshipProgram{}
		if err := rows.Scan(
			&program.ID,
			&program.MentorID,
			&program.Title,
			&program.Description,
			&program.Duration,
			&program.Price,
			&program.MaxMentees,
			&program.Status,
			&program.CreatedAt,
			&program.UpdatedAt,
		); err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}
	return programs, nil
}

// CreateSessionFeedback creates a new feedback entry for a session
func (r *MentorshipRepository) CreateSessionFeedback(ctx context.Context, feedback *models.SessionFeedback) error {
	query := `
        INSERT INTO session_feedback (session_id, user_id, rating, comment, created_at)
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
        RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		feedback.SessionID,
		feedback.UserID,
		feedback.Rating,
		feedback.Comment,
	).Scan(&feedback.ID)
}

// CreateSession creates a new mentorship session
func (r *MentorshipRepository) CreateSession(ctx context.Context, session *models.MentorshipSession) error {
	query := `
        INSERT INTO mentorship_sessions (request_id, title, start_time, end_time, status, notes)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		session.RequestID,
		session.Title,
		session.StartTime,
		session.EndTime,
		session.Status,
		session.Notes,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSession retrieves a session by ID
func (r *MentorshipRepository) GetSession(ctx context.Context, sessionID int) (*models.MentorshipSession, error) {
	query := `
        SELECT id, request_id, title, start_time, end_time, status, notes, created_at, updated_at
        FROM mentorship_sessions
        WHERE id = $1`

	session := &models.MentorshipSession{}
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.RequestID,
		&session.Title,
		&session.StartTime,
		&session.EndTime,
		&session.Status,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// UpdateSessionStatus updates the status of a session
func (r *MentorshipRepository) UpdateSessionStatus(ctx context.Context, sessionID int, status string) error {
	query := `
        UPDATE mentorship_sessions
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("session not found")
	}

	return nil
}

// ListSessionsByRequest retrieves all sessions for a mentorship request
func (r *MentorshipRepository) ListSessionsByRequest(ctx context.Context, requestID int) ([]*models.MentorshipSession, error) {
	query := `
        SELECT id, request_id, title, start_time, end_time, status, notes, created_at, updated_at
        FROM mentorship_sessions
        WHERE request_id = $1
        ORDER BY start_time ASC`

	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.MentorshipSession
	for rows.Next() {
		session := &models.MentorshipSession{}
		err := rows.Scan(
			&session.ID,
			&session.RequestID,
			&session.Title,
			&session.StartTime,
			&session.EndTime,
			&session.Status,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// GetAverageRating calculates the average rating for a mentor
func (r *MentorshipRepository) GetAverageRating(ctx context.Context, mentorID int) (float64, error) {
	query := `
        SELECT COALESCE(AVG(sf.rating), 0)
        FROM session_feedback sf
        JOIN mentorship_sessions ms ON sf.session_id = ms.id
        JOIN mentorship_requests mr ON ms.request_id = mr.id
        WHERE mr.mentor_id = $1 AND sf.rating IS NOT NULL`

	var avgRating float64
	err := r.db.QueryRowContext(ctx, query, mentorID).Scan(&avgRating)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avgRating, nil
}

// UpdateProgramStatus updates the status of a mentorship program
func (r *MentorshipRepository) UpdateProgramStatus(ctx context.Context, programID int, status string) error {
	query := `
        UPDATE mentorship_programs
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, programID)
	if err != nil {
		return fmt.Errorf("failed to update program status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("program not found")
	}

	return nil
}

// Helper function to validate status transitions
func isValidStatusTransition(current, new string) bool {
	transitions := map[string][]string{
		"pending":   {"approved", "rejected"},
		"approved":  {"completed", "cancelled"},
		"rejected":  {},
		"completed": {},
		"cancelled": {},
	}

	validTransitions, exists := transitions[current]
	if !exists {
		return false
	}

	for _, validStatus := range validTransitions {
		if validStatus == new {
			return true
		}
	}
	return false
}

// GetSessionFeedback retrieves feedback for a specific session
func (r *MentorshipRepository) GetSessionFeedback(ctx context.Context, sessionID int) ([]*models.SessionFeedback, error) {
	query := `
        SELECT id, session_id, user_id, rating, comment, created_at
        FROM session_feedback
        WHERE session_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session feedback: %w", err)
	}
	defer rows.Close()

	var feedbacks []*models.SessionFeedback
	for rows.Next() {
		feedback := &models.SessionFeedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.SessionID,
			&feedback.UserID,
			&feedback.Rating,
			&feedback.Comment,
			&feedback.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan feedback: %w", err)
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feedback: %w", err)
	}

	return feedbacks, nil
}

// UpdateRequestStatus updates the status of a mentorship request
func (r *MentorshipRepository) UpdateRequestStatus(ctx context.Context, requestID int, status string) error {
	query := `
        UPDATE mentorship_requests 
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, requestID)
	if err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("request not found")
	}

	return nil
}
