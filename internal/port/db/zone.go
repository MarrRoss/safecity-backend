package db

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type ZoneRepository interface {
	AddZone(ctx context.Context, zone *entity.Zone) error
	ZoneOverlaps(ctx context.Context, familyID value_object.ID, boundaries *[]value_object.Location) (bool, error)
	GetZone(ctx context.Context, id value_object.ID) (*entity.Zone, error)
	CoordinatesInFamilyZones(ctx context.Context, familyID value_object.ID, lon, lat float64) ([]*entity.Zone, error)
	GetZones(ctx context.Context) ([]*entity.Zone, error)
	//GetZonesByAuthorID(ctx context.Context, id value_object.ID) ([]*entity.Zone, error)
	//GetAvailableZonesForSubscription(ctx context.Context, receiverID, senderID value_object.ID) ([]*entity.Zone, error)
	UpdateZone(ctx context.Context, zone *entity.Zone) error
	//DeleteZone(ctx context.Context, id value_object.externalID) error
}
