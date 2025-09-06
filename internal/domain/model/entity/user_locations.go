package entity

import (
	"awesomeProjectDDD/internal/domain/model/value_object"
	"time"
)

type UserLocation struct {
	ID        value_object.ID
	UserID    value_object.ID
	Location  value_object.Location
	Battery   value_object.BatteryThreshold
	CreatedAt time.Time
}

func NewUserLocation(
	userID value_object.ID,
	location value_object.Location,
	battery value_object.BatteryThreshold,
) *UserLocation {
	id := value_object.NewID()
	return &UserLocation{
		ID:        id,
		UserID:    userID,
		Location:  location,
		Battery:   battery,
		CreatedAt: time.Now().UTC(),
	}
}
