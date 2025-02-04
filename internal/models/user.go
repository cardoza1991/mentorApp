package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int            `json:"id"`
	Username          string         `json:"username"`
	Email             string         `json:"email"`
	PasswordHash      string         `json:"-"`
	Role              string         `json:"role"`
	IsMentor          bool           `json:"is_mentor"`
	IsAdmin           bool           `json:"is_admin"`
	IsApproved        bool           `json:"is_approved"`
	EmailVerified     bool           `json:"email_verified"`
	VerificationToken sql.NullString `json:"verification_token"`
	ResetToken        sql.NullString `json:"reset_token"`
	ResetTokenExpiry  *time.Time     `json:"reset_token_expiry,omitempty"`
	LastLoginAt       *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	u.PasswordHash = string(hash)
	log.Printf("Successfully hashed password for user")
	return nil
}

// ValidatePassword checks if the provided password is correct
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("Password validation failed: %v", err)
		return false
	}
	return true
}

// ValidateEmail ensures email is from underground-ops.dev
func (u *User) ValidateEmail() error {
	if !strings.HasSuffix(u.Email, "@underground-ops.dev") {
		return fmt.Errorf("email must be from underground-ops.dev domain")
	}
	return nil
}

// BeforeInsert is called before inserting a new user
func (u *User) BeforeInsert() {
	if u.Role == "" {
		if u.IsMentor {
			u.Role = "mentor"
		} else {
			u.Role = "mentee"
		}
	}
}

// TableUnique sets unique constraints
func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username"},
		[]string{"Email"},
	}
}

// TableIndex sets indexed fields
func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Email"},
		[]string{"Role"},
	}
}

func init() {
	orm.RegisterModel(new(User))
}

// IsAdminUser checks if the user is an admin
func (u *User) IsAdminUser() bool {
	return u.IsAdmin
}

// MakeAdmin promotes a user to admin status
func (u *User) MakeAdmin() {
	u.IsAdmin = true
	u.IsApproved = true // Admins are automatically approved
}

// RemoveAdmin removes admin status from a user
func (u *User) RemoveAdmin() {
	u.IsAdmin = false
}
