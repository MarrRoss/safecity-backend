package value_object

import (
	"awesomeProjectDDD/internal/domain"
	"errors"
	"fmt"
	"strconv"
)

type Longitude float64

func NewLongitude(long float64) (Longitude, error) {
	longFloat, err := strconv.ParseFloat(fmt.Sprintf("%f", long), 64)
	if err != nil {
		return 000, fmt.Errorf("failed to convert to float 64: %w", domain.ErrInvalidLongitude)
	}
	newLong := Longitude(longFloat)
	if err := newLong.ValidateLongitude(); err != nil {
		return 000, fmt.Errorf("%w: %w", err, domain.ErrInvalidLatitude)
	}
	return newLong, nil
}

func (long Longitude) ValidateLongitude() error {
	if long < -180 || long > 180 {
		return errors.New("longitude must be in the range [-180;180]")
	}
	return nil
}
