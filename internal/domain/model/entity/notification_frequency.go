package entity

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type NotificationFrequency struct {
	ID        value_object.ID
	Frequency time.Duration
}

func NewFrequency(f string) (NotificationFrequency, error) {
	if f == "" {
		return NotificationFrequency{}, fmt.Errorf("notification frequency is empty: %w", domain.ErrInvalidFrequency)
	}
	freq, err := time.ParseDuration(f)
	if err != nil {
		return NotificationFrequency{}, fmt.Errorf("error parsing frequency to time.Duration: %w", domain.ErrInvalidFrequency)
	}
	id := value_object.NewID()
	newFrequency := NotificationFrequency{
		ID:        id,
		Frequency: freq,
	}
	return newFrequency, nil
}

func (f NotificationFrequency) String() string {
	return f.Frequency.String()
}
