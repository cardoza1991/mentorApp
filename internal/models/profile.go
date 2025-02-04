package models

import (
	"time"
)

type Profile struct {
	Id                      int
	UserId                  int
	FirstName               string
	LastName                string
	Bio                     string
	Skills                  string
	Experience              string
	LinkedIn                string
	Github                  string
	Twitter                 string
	Rate                    float64
	Available               bool
	IsApproved              bool // Add this field
	Timezone                string
	ProfilePicture          string
	NotificationPreferences string
	PrivacySettings         string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// NotificationSettings represent user notification preferences
type NotificationSettings struct {
	EmailNotifications   bool `json:"email_notifications"`
	SessionReminders     bool `json:"session_reminders"`
	MessageNotifications bool `json:"message_notifications"`
	UpdatesNotifications bool `json:"updates_notifications"`
}

// PublicProfile represents a user’s public-facing profile information
type PublicProfile struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email,omitempty"`
	Bio            string  `json:"bio"`
	Skills         string  `json:"skills"`
	Rate           float64 `json:"rate"`
	ProfilePicture string  `json:"profile_picture,omitempty"`
	AverageRating  float64 `json:"average_rating,omitempty"`
	RatingCount    int     `json:"rating_count,omitempty"`
}

// ProfileSettings for user preferences and privacy
type ProfileSettings struct {
	NotificationPreferences string `json:"notification_preferences"`
	PrivacySettings         string `json:"privacy_settings"`
}

// Specialty represents areas of expertise
type Specialty struct {
	Id        int
	ProfileId int
	Name      string
	Level     string
	CreatedAt time.Time
}

type Availability struct {
	ID        int       `json:"id"`
	ProfileID int       `json:"profile_id"` // Changed from ProfileId
	UserID    int       `json:"user_id"`    // Added UserID field
	DayOfWeek int       `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Experience represents a user's professional experience
type Experience struct {
	ID          int
	UserID      int
	Title       string
	Company     string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Education represents a user's educational background
type Education struct {
	ID          int
	UserID      int
	Degree      string
	Institution string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rating represents mentor ratings
type Rating struct {
	Id         int
	FromUserId int
	ToUserId   int
	Rating     float64
	Comment    string
	CreatedAt  time.Time
}

type ProfileUpdate struct {
	FirstName      *string
	LastName       *string
	Bio            *string
	Skills         *string
	TimeZone       *string
	Rate           *float64
	Availability   *string
	ProfilePicture *string
}
