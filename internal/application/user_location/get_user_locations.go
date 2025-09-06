package user_location

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
)

type GetUserLocationsQuery struct {
	UserID string
}

type GetUserLocationsHandler struct {
	userStorage     db.UserRepository
	settingsStorage db.NotificationSettingRepository
	locationStorage db.UserLocationRepository
	observer        *observability.Observability
}

func NewGetUserLocationsHandler(
	userStore db.UserRepository,
	settingsStorage db.NotificationSettingRepository,
	locationStore db.UserLocationRepository,
	observer *observability.Observability) *GetUserLocationsHandler {
	return &GetUserLocationsHandler{
		userStorage:     userStore,
		settingsStorage: settingsStorage,
		locationStorage: locationStore,
		observer:        observer,
	}
}

func (h *GetUserLocationsHandler) Handle(ctx context.Context,
	query GetUserLocationsQuery) ([]*app_model.ApplicationUserLocation, error) {
	extID, err := value_object.NewIDFromString(query.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return nil, err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return nil, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return nil, err
	}

	childIDs, err := h.settingsStorage.FindLocationSenderIDsByReceiver(ctx, user.ID)
	fmt.Println(childIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting children ids")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting children ids")
		return nil, err
	}

	usersLocations, err := h.locationStorage.FindLatestLocationsByUserIDs(ctx, childIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting users locations")
			return nil, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting users locations")
		return nil, err
	}

	//usersLocations, err := h.locationStorage.GetLastChildrenLocations(ctx, user.ID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting users locations")
	//		return nil, err
	//	}
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting users locations")
	//	return nil, err
	//}
	//locations, err := h.locationStorage.GetLocationsByUser(ctx, user.ID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get locations: %w", err)
	//}
	//if len(locations) == 0 {
	//	return nil, errors.New("locations not found")
	//}

	//appLocations := make([]*app_model.ApplicationLocationTime, len(locations))
	//for i, loc := range locations {
	//	appLocations[i] = app_model.NewApplicationLocationTime(loc)
	//}

	return app_model.NewApplicationUserLocations(usersLocations), nil
}
