package initialization

import (
	"mentorApp/internal/api/handlers"
	"mentorApp/internal/services"
)

// InitializeHandlers creates and returns all handler instances
func InitializeHandlers(
	userService services.IUserService,
	mentorshipService services.IMentorshipService,
) (*handlers.UserHandler, *handlers.MentorshipHandler, *handlers.ProfileHandler, *handlers.HomeHandler) {

	userHandler := handlers.NewUserHandler(userService)
	mentorshipHandler := handlers.NewMentorshipHandler(mentorshipService)
	profileHandler := handlers.NewProfileHandler(userService)
	homeHandler := handlers.NewHomeHandler(userService, mentorshipService)

	return userHandler, mentorshipHandler, profileHandler, homeHandler
}
