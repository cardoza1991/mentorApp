package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mentorApp/internal/api/handlers/common"
	"mentorApp/internal/models"
	"mentorApp/internal/repository"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	db          *sql.DB
	userRepo    repository.IUserRepository
	profileRepo repository.IProfileRepository
	templates   *template.Template
}

func NewAdminHandler(db *sql.DB, userRepo repository.IUserRepository, profileRepo repository.IProfileRepository, templates *template.Template) *AdminHandler {
	return &AdminHandler{
		db:          db,
		userRepo:    userRepo,
		profileRepo: profileRepo,
		templates:   templates,
	}
}

func (h *AdminHandler) DB() *sql.DB {
	return h.db
}

func (h *AdminHandler) Setup(w http.ResponseWriter, r *http.Request) {
	// Check if admin exists
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = true").Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// If admin exists, redirect to login
	if count > 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("admin_email")
		password := r.FormValue("password")

		// Validate input
		if email == "" || password == "" {
			data := map[string]interface{}{
				"Error":   "Email and password are required",
				"Website": "NEXUS Mentorship Platform",
			}
			h.templates.ExecuteTemplate(w, "admin_setup.html", data)
			return
		}

		// Validate email domain
		if !strings.HasSuffix(email, "@underground-ops.dev") {
			data := map[string]interface{}{
				"Error":   "Email must be from @underground-ops.dev domain",
				"Website": "NEXUS Mentorship Platform",
			}
			h.templates.ExecuteTemplate(w, "admin_setup.html", data)
			return
		}

		// Begin transaction
		tx, err := h.db.BeginTx(r.Context(), nil)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Create a user model to handle the password
		adminUser := &models.User{
			Username:      strings.Split(email, "@")[0],
			Email:         email,
			Role:          "admin",
			IsAdmin:       true,
			IsApproved:    true,
			EmailVerified: true,
			VerificationToken: sql.NullString{
				String: "",
				Valid:  false,
			},
		}

		// Hash the password using the user model's method
		if err := adminUser.SetPassword(password); err != nil {
			data := map[string]interface{}{
				"Error":   "Password processing error: " + err.Error(),
				"Website": "NEXUS Mentorship Platform",
			}
			h.templates.ExecuteTemplate(w, "admin_setup.html", data)
			return
		}

		// Create admin user
		var userID int
		err = tx.QueryRowContext(
			r.Context(),
			`INSERT INTO users (
				username,
				email,
				password_hash,
				role,
				is_mentor,
				is_admin,
				is_approved,
				email_verified,
				created_at,
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id`,
			adminUser.Username,
			adminUser.Email,
			adminUser.PasswordHash,
			"admin", // Explicitly set role as admin
			false,   // is_mentor
			true,    // is_admin - Make sure this is true
			true,    // is_approved
			true,    // email_verified
		).Scan(&userID)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				data := map[string]interface{}{
					"Error":   "Email already registered",
					"Website": "NEXUS Mentorship Platform",
				}
				h.templates.ExecuteTemplate(w, "admin_setup.html", data)
				return
			}
			http.Error(w, "Failed to create admin user", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to complete setup", http.StatusInternalServerError)
			return
		}

		// Redirect to login
		http.Redirect(w, r, "/login?setup=complete", http.StatusSeeOther)
		return
	}

	// Show setup form
	data := map[string]interface{}{
		"Website": "NEXUS Mentorship Platform",
	}
	h.templates.ExecuteTemplate(w, "admin_setup.html", data)
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get stats
	var stats struct {
		TotalUsers       int
		PendingApprovals int
		ActiveMentors    int
		ActiveMentees    int
		NewUsersToday    int
		ActiveSessions   int
		OpenJobs         int
	}

	// Collect statistics
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_approved = false AND is_mentor = true").Scan(&stats.PendingApprovals)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_mentor = true AND is_approved = true").Scan(&stats.ActiveMentors)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_mentor = false").Scan(&stats.ActiveMentees)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE").Scan(&stats.NewUsersToday)
	h.db.QueryRow("SELECT COUNT(*) FROM mentorship_sessions WHERE status = 'scheduled'").Scan(&stats.ActiveSessions)
	h.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE status = 'active'").Scan(&stats.OpenJobs)

	// Get pending mentors
	query := `
		SELECT u.id, p.first_name, p.last_name, u.email, u.is_mentor
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.is_approved = false AND u.is_mentor = true
		ORDER BY u.created_at DESC
	`
	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		log.Printf("Error fetching pending profiles: %v", err)
		http.Error(w, "Failed to fetch pending profiles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pendingProfiles []map[string]interface{}
	for rows.Next() {
		profile := make(map[string]interface{})
		var id int
		var firstName, lastName, email string
		var isMentor bool

		err := rows.Scan(&id, &firstName, &lastName, &email, &isMentor)
		if err != nil {
			log.Printf("Error scanning profile: %v", err)
			continue
		}

		profile["id"] = id
		profile["firstName"] = firstName
		profile["lastName"] = lastName
		profile["email"] = email
		profile["isMentor"] = isMentor

		pendingProfiles = append(pendingProfiles, profile)
	}

	data := map[string]interface{}{
		"Stats":           stats,
		"PendingProfiles": pendingProfiles,
		"Website":         "NEXUS Mentorship Platform",
	}

	err = h.templates.ExecuteTemplate(w, "admin_dashboard.html", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	common.RespondJSON(w, http.StatusOK, users)
}

func (h *AdminHandler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.profileRepo.SearchProfiles(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to fetch profiles", http.StatusInternalServerError)
		return
	}
	common.RespondJSON(w, http.StatusOK, profiles)
}

func (h *AdminHandler) ApproveMentor(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Approved bool `json:"approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update both approval and email verification status
	result, err := h.db.ExecContext(r.Context(), `
        UPDATE users
        SET is_approved = $1, email_verified = $1
        WHERE id = $2 AND is_mentor = true`,
		req.Approved, userID)

	if err != nil {
		log.Printf("Failed to update mentor status: %v", err)
		http.Error(w, "Failed to update mentor status", http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Mentor not found", http.StatusNotFound)
		return
	}

	common.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Mentor status updated successfully",
	})
}

// Job management methods
func (h *AdminHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
        SELECT id, title, company, location, description, salary_range, 
               remote_policy, status, is_featured, created_at
        FROM jobs
        WHERE status != $1
        ORDER BY created_at DESC`,
		models.JobStatus.Closed)
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		if err := rows.Scan(
			&job.ID, &job.Title, &job.Company, &job.Location, &job.Description,
			&job.SalaryRange, &job.RemotePolicy, &job.Status, &job.IsFeatured,
			&job.CreatedAt); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	common.RespondJSON(w, http.StatusOK, jobs)
}

func (h *AdminHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job.Status = models.JobStatus.Active

	err := h.db.QueryRow(`
        INSERT INTO jobs (title, company, location, description, salary_range, remote_policy, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		job.Title, job.Company, job.Location, job.Description,
		job.SalaryRange, job.RemotePolicy, job.Status).Scan(&job.ID)

	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	common.RespondJSON(w, http.StatusCreated, job)
}

func (h *AdminHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.Atoi(chi.URLParam(r, "jobId"))
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("UPDATE jobs SET status = $1 WHERE id = $2",
		models.JobStatus.Closed, jobID)
	if err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) FeatureJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.Atoi(chi.URLParam(r, "jobId"))
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Featured bool `json:"featured"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("UPDATE jobs SET is_featured = $1 WHERE id = $2",
		req.Featured, jobID)
	if err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
