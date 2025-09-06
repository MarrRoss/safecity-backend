package app_model

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
)

type ApplicationUserLocation struct {
	User     *ApplicationUser
	Location *ApplicationLocation
}

func NewApplicationUserLocation(usersLocation *response.UserLatestLocations) *ApplicationUserLocation {
	id, err := value_object.NewIDFromString(usersLocation.UserId)
	if err != nil {
		return nil
	}
	telegramExists := false
	if usersLocation.TelegramID != nil {
		telegramExists = true
	}
	return &ApplicationUserLocation{
		User: NewApplicationUser(&entity.User{
			ID:           id,
			Name:         value_object.NewFullName(value_object.FirstName(usersLocation.FirstName), value_object.LastName(usersLocation.LastName)),
			Email:        value_object.Email(usersLocation.Email),
			Login:        value_object.Login(usersLocation.Login),
			Tracking:     usersLocation.Tracking,
			CreatedAt:    usersLocation.CreatedAt,
			UpdatedAt:    usersLocation.UpdatedAt,
			TgIntegrated: telegramExists,
		}),
		Location: NewApplicationLocation(usersLocation.Latitude, usersLocation.Longitude),
	}
}

func NewApplicationUserLocations(usersLocations []*response.UserLatestLocations) []*ApplicationUserLocation {
	appLocations := make([]*ApplicationUserLocation, len(usersLocations))
	for i, location := range usersLocations {
		appLocations[i] = NewApplicationUserLocation(location)
	}
	return appLocations
}
