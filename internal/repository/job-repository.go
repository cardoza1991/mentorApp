package repository

import (
	"context"
	"database/sql"
	"mentorApp/internal/models"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
        INSERT INTO jobs (
            title, company, location, description, salary_range,
            job_type, experience_level, remote_policy, contact_email,
            status, is_featured, created_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		job.Title, job.Company, job.Location, job.Description,
		job.SalaryRange, job.JobType, job.ExperienceLevel,
		job.RemotePolicy, job.ContactEmail, job.Status,
		job.IsFeatured, job.CreatedBy,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id int) (*models.Job, error) {
	job := &models.Job{}
	query := `
        SELECT id, title, company, location, description,
               salary_range, job_type, experience_level, remote_policy,
               contact_email, status, is_featured, created_by,
               created_at, updated_at
        FROM jobs
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &job.Title, &job.Company, &job.Location,
		&job.Description, &job.SalaryRange, &job.JobType,
		&job.ExperienceLevel, &job.RemotePolicy, &job.ContactEmail,
		&job.Status, &job.IsFeatured, &job.CreatedBy,
		&job.CreatedAt, &job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return job, err
}

func (r *JobRepository) List(ctx context.Context, filters map[string]interface{}) ([]*models.Job, error) {
	query := `
        SELECT id, title, company, location, description,
               salary_range, job_type, experience_level, remote_policy,
               contact_email, status, is_featured, created_by,
               created_at, updated_at
        FROM jobs
        WHERE status = $1`

	args := []interface{}{models.JobStatus.Active}

	// Add filters
	if loc, ok := filters["location"].(string); ok && loc != "" {
		query += ` AND location = $` + string(len(args)+1)
		args = append(args, loc)
	}

	query += ` ORDER BY is_featured DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		err := rows.Scan(
			&job.ID, &job.Title, &job.Company, &job.Location,
			&job.Description, &job.SalaryRange, &job.JobType,
			&job.ExperienceLevel, &job.RemotePolicy, &job.ContactEmail,
			&job.Status, &job.IsFeatured, &job.CreatedBy,
			&job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *JobRepository) Update(ctx context.Context, job *models.Job) error {
	query := `
        UPDATE jobs
        SET title = $1, company = $2, location = $3,
            description = $4, salary_range = $5,
            job_type = $6, experience_level = $7,
            remote_policy = $8, contact_email = $9,
            status = $10, is_featured = $11,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $12`

	result, err := r.db.ExecContext(ctx, query,
		job.Title, job.Company, job.Location,
		job.Description, job.SalaryRange,
		job.JobType, job.ExperienceLevel,
		job.RemotePolicy, job.ContactEmail,
		job.Status, job.IsFeatured, job.ID)

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

func (r *JobRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE jobs SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, models.JobStatus.Closed, id)
	return err
}

func (r *JobRepository) CreateApplication(ctx context.Context, app *models.JobApplication) error {
	query := `
        INSERT INTO job_applications (
            job_id, user_id, status, cover_letter, resume_url
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id, applied_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		app.JobID, app.UserID, app.Status,
		app.CoverLetter, app.ResumeURL,
	).Scan(&app.ID, &app.AppliedAt, &app.UpdatedAt)
}

func (r *JobRepository) GetApplications(ctx context.Context, jobID int) ([]*models.JobApplication, error) {
	query := `
        SELECT id, job_id, user_id, status, cover_letter,
               resume_url, applied_at, updated_at
        FROM job_applications
        WHERE job_id = $1
        ORDER BY applied_at DESC`

	rows, err := r.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []*models.JobApplication
	for rows.Next() {
		app := &models.JobApplication{}
		err := rows.Scan(
			&app.ID, &app.JobID, &app.UserID,
			&app.Status, &app.CoverLetter, &app.ResumeURL,
			&app.AppliedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}
