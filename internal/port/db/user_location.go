package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type UserLocationRepository interface {
	AddUserLocation(ctx context.Context, userLocation *entity.UserLocation) error
	AddLocationContext(ctx context.Context, locationLogID value_object.ID, zoneID *value_object.ID, notifyType string, battery *int) error
	GetLastLocationContext(ctx context.Context, userID value_object.ID) (*response.LocationContextDB, error)
	//GetLocationsByUser(ctx context.Context, userID value_object.ID) ([]*entity.UserLocation, error)
	GetLastZoneContext(ctx context.Context, userID value_object.ID) (*response.LocationContextDB, error)
	FindLatestLocationsByUserIDs(ctx context.Context, userIDs []value_object.ID) ([]*response.UserLatestLocations, error)
}
