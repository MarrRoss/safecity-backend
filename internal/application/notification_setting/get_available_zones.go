package notification_setting

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/port/db"
	"context"
)

type GetAvailableZonesQuery struct {
	ReceiverID string
	SenderID   string
}

type GetAvailableZonesHandler struct {
	zoneStorage db.ZoneRepository
}

func NewGetAvailableZonesHandler(zoneStorage db.ZoneRepository) *GetAvailableZonesHandler {
	return &GetAvailableZonesHandler{
		zoneStorage: zoneStorage,
	}
}

func (h *GetAvailableZonesHandler) Handle(ctx context.Context, query GetAvailableZonesQuery) ([]*app_model.ApplicationZone, error) {
	//receiverID, err := value_object.NewIDFromString(query.ReceiverID)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid receiver id: %w", err)
	//}
	//senderID, err := value_object.NewIDFromString(query.SenderID)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid sender id: %w", err)
	//}
	//
	//zones, err := h.zoneStorage.GetAvailableZonesForSubscription(ctx, receiverID, senderID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get available zones: %w", err)
	//}

	return app_model.NewApplicationZones(nil), nil
}
