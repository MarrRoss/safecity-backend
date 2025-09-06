package zone

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"github.com/google/uuid"
)

type GetZonesByUserIDQuery struct {
	ID uuid.UUID
}

type GetZonesByUserIDHandler struct {
	zoneStorage db.ZoneRepository
	userStorage db.UserRepository
	observer    *observability.Observability
}

func NewGetZonesByUserIDHandler(
	zoneStorage db.ZoneRepository,
	userStorage db.UserRepository,
	observer *observability.Observability,
) *GetZonesByUserIDHandler {
	return &GetZonesByUserIDHandler{
		zoneStorage: zoneStorage,
		userStorage: userStorage,
		observer:    observer,
	}
}

func (h *GetZonesByUserIDHandler) Handle(ctx context.Context,
	query GetZonesByUserIDQuery) ([]*app_model.ApplicationZone, error) {
	//userID, err := value_object.NewIDFromString(query.externalID.String())
	//if err != nil {
	//	h.observer.Logger.Trace().Err(err).Msg("invalid author id")
	//	return nil, err
	//}
	//user, err := h.userStorage.GetUser(ctx, userID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrUserNotFound) {
	//		h.observer.Logger.Trace().Err(err).Msg("author not found")
	//		return nil, err
	//	}
	//
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting zone author")
	//		return nil, err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone author")
	//	return nil, err
	//}
	//if user.EndedAt != nil {
	//	h.observer.Logger.Error().Err(err).Msg("author is deleted")
	//	return nil, domain.ErrUserIsDeleted
	//}
	//zones, err := h.zoneStorage.GetZonesByAuthorID(ctx, user.externalID)
	//if err != nil {
	//	if errors.Is(err, adapter.ErrStorage) {
	//		h.observer.Logger.Error().Err(err).Msg("database error while getting zones")
	//		return nil, err
	//	}
	//
	//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zones")
	//	return nil, err
	//}

	return app_model.NewApplicationZones(nil), nil
}
