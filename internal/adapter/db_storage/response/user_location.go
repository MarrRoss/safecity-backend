package response

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type UserLatestLocations struct {
	UserId     string    `db:"id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Email      string    `db:"email"`
	Login      string    `db:"username"`
	Tracking   bool      `db:"tracking"`
	TelegramID *string   `db:"telegram_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Latitude   float64   `db:"latitude"`
	Longitude  float64   `db:"longitude"`
}

type UserLocationDB struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Location  string    `db:"location"`
	CreatedAt time.Time `db:"created_at"`
}

type LocationContextDB struct {
	LocationID string    `db:"location_log_id"`
	ZoneID     *string   `db:"zone_id"`
	NotifyType string    `db:"notification_type"`
	CreatedAt  time.Time `db:"created_at"`
}

func UserLocationDbToEntity(db *UserLocationDB) (*entity.UserLocation, error) {
	id, err := value_object.NewIDFromString(db.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to id: %w", err)
	}
	userID, err := value_object.NewIDFromString(db.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to user id: %w", err)
	}

	var lon, lat float64
	_, err = fmt.Sscanf(db.Location, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return nil, fmt.Errorf("failed to parse location: %w", err)
	}
	location := value_object.NewLocation(value_object.Latitude(lat), value_object.Longitude(lon))

	return &entity.UserLocation{
		ID:        id,
		UserID:    userID,
		Location:  location,
		CreatedAt: db.CreatedAt,
	}, nil
}

func UserLocationDbListToEntityList(dbList []*UserDB) ([]*entity.User, error) {
	var users []*entity.User
	for _, db := range dbList {
		user, err := UserDbToEntity(db)
		if err != nil {
			return nil, fmt.Errorf("failed to convert db data to entity")
		}
		users = append(users, user)
	}
	return users, nil
}
