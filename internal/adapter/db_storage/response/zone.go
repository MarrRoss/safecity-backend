package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ZoneDB struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Safety     bool      `db:"safety"`
	Boundaries string    `db:"boundaries"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	FamilyID   string    `db:"family_id"`
}

func ZoneDBToEntity(zoneDB *ZoneDB) (*entity.Zone, error) {
	// Проверяем и преобразуем границы из WKT
	boundaries, err := ConvertWKTToBoundaries(zoneDB.Boundaries)
	if err != nil {
		return nil, fmt.Errorf("failed to convert WKT to boundaries: %w; WKT: %s", err, zoneDB.Boundaries)
	}

	zoneID, err := value_object.NewIDFromString(zoneDB.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to zone id: %w; id: %s", err, zoneDB.ID)
	}

	name, err := value_object.NewZoneName(zoneDB.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to zone name: %w; Name: %s", err, zoneDB.Name)
	}

	safety, err := value_object.NewSafety(zoneDB.Safety)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to safety: %w; Safety: %d", err, zoneDB.Safety)
	}

	familyID, err := value_object.NewIDFromString(zoneDB.FamilyID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to family id: %w; id: %s", err, zoneDB.ID)
	}

	return &entity.Zone{
		ID:         zoneID,
		Name:       name,
		Boundaries: &boundaries,
		Safety:     safety,
		CreatedAt:  zoneDB.CreatedAt,
		UpdatedAt:  zoneDB.UpdatedAt,
		FamilyID:   familyID,
	}, nil
}

func ZoneDbListToEntityList(zoneDBList []*ZoneDB) ([]*entity.Zone, error) {
	var zones []*entity.Zone
	for _, zoneDB := range zoneDBList {
		zone, err := ZoneDBToEntity(zoneDB)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to convert storage data to entity", err)
		}
		zones = append(zones, zone)
	}
	return zones, nil
}

func ConvertWKTToBoundaries(wkt string) ([]value_object.Location, error) {
	if !strings.HasPrefix(wkt, "POLYGON((") || !strings.HasSuffix(wkt, "))") {
		return nil, fmt.Errorf("invalid WKT format: %s", wkt)
	}

	// Убираем "POLYGON((" и "))"
	coords := strings.TrimSuffix(strings.TrimPrefix(wkt, "POLYGON(("), "))")

	var boundaries []value_object.Location
	points := strings.Split(coords, ",")
	for _, point := range points {
		coord := strings.Fields(point) // Разделяем по пробелу
		if len(coord) != 2 {
			return nil, fmt.Errorf("invalid coordinate pair: %s", point)
		}

		lon, err := strconv.ParseFloat(coord[0], 64) // Долгота
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %w", err)
		}

		lat, err := strconv.ParseFloat(coord[1], 64) // Широта
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %w", err)
		}

		boundaries = append(boundaries, value_object.Location{
			Longitude: value_object.Longitude(lon),
			Latitude:  value_object.Latitude(lat),
		})
	}

	return boundaries, nil
}
