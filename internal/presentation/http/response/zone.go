package response

import (
	"awesomeProjectDDD/internal/application/app_model"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"github.com/google/uuid"
	"time"
)

type GetZoneResponse struct {
	ID         uuid.UUID         `json:"zone_id"`
	Name       string            `json:"zone_name"`
	Safety     bool              `json:"zone_safety"`
	Boundaries []BoundariesModel `json:"zone_boundaries"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
} // @name GetZoneResponse

type BoundariesModel struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewGetZoneResponse(zone *app_model.ApplicationZone) *GetZoneResponse {
	boundaries := make([]BoundariesModel, len(zone.Boundaries))
	for i, boundary := range zone.Boundaries {
		boundaries[i] = BoundariesModel{
			Latitude:  boundary.Latitude,
			Longitude: boundary.Longitude,
		}
	}
	return &GetZoneResponse{
		ID:         zone.ID,
		Name:       zone.Name,
		Safety:     zone.Safety,
		Boundaries: boundaries,
		CreatedAt:  zone.CreatedAt,
		UpdatedAt:  zone.UpdatedAt,
	}
}

func NewGetZonesResponse(zones []*app_model.ApplicationZone) []*GetZoneResponse {
	respZones := make([]*GetZoneResponse, len(zones))
	for i, zone := range zones {
		respZones[i] = NewGetZoneResponse(zone)
	}
	return respZones
}

type AddZoneResponse struct {
	ID uuid.UUID `json:"id"`
} // @name AddZoneResponse

func NewAddZoneResponse(id value_object.ID) *AddZoneResponse {
	return &AddZoneResponse{ID: id.ToRaw()}
}
