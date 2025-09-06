package value_object

import (
	"awesomeProjectDDD/internal/domain"
	"errors"
	"fmt"
	"strconv"
)

type Latitude float64

func NewLatitude(lat float64) (Latitude, error) {
	latFloat, err := strconv.ParseFloat(fmt.Sprintf("%f", lat), 64)
	if err != nil {
		return 000, fmt.Errorf("failed to convert to float 64: %w", domain.ErrInvalidLatitude)
	}
	newLat := Latitude(latFloat)
	if err := newLat.ValidateLatitude(); err != nil {
		return 000, fmt.Errorf("%w: %w", err, domain.ErrInvalidLatitude)
	}
	return newLat, nil
}

func (lat Latitude) ValidateLatitude() error {
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be in the range [-90;90]")
	}
	return nil
}
