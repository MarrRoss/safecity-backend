package app_model

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"github.com/google/uuid"
	"time"
)

type ApplicationZone struct {
	ID         uuid.UUID
	Name       string
	Safety     bool
	Boundaries []ApplicationLocation
	FamilyID   uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ApplicationLocation struct {
	Latitude  float64
	Longitude float64
}

func NewApplicationLocation(lat, lon float64) *ApplicationLocation {
	return &ApplicationLocation{
		Latitude:  lat,
		Longitude: lon,
	}
}

func NewApplicationZone(zone *entity.Zone) *ApplicationZone {
	appBoundaries := make([]ApplicationLocation, len(*zone.Boundaries))
	for i, boundary := range *zone.Boundaries {
		appBoundaries[i] = ApplicationLocation{
			Latitude:  float64(boundary.Latitude),
			Longitude: float64(boundary.Longitude),
		}
	}

	return &ApplicationZone{
		ID:         zone.ID.ToRaw(),
		Name:       zone.Name.String(),
		Safety:     zone.Safety.Bool(),
		Boundaries: appBoundaries,
		FamilyID:   zone.FamilyID.ToRaw(),
		CreatedAt:  zone.CreatedAt,
		UpdatedAt:  zone.UpdatedAt,
	}
}

func NewApplicationZones(zones []*entity.Zone) []*ApplicationZone {
	appZones := make([]*ApplicationZone, len(zones))
	for i, zone := range zones {
		appZones[i] = NewApplicationZone(zone)
	}
	return appZones
}
