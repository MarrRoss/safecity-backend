package value_object

import (
	"awesomeProjectDDD/internal/domain"
	"fmt"
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func NewIDFromString(id string) (ID, error) {
	if id == "" {
		return ID{}, fmt.Errorf("id is empty: %w", domain.ErrInvalidID)
	}
	u, err := uuid.Parse(id)
	if err != nil {
		return ID{}, fmt.Errorf("failed to parse id from string %s: %w", id, domain.ErrInvalidID)
	}
	return ID(u), nil
}

func (id ID) ToRaw() uuid.UUID {
	return uuid.UUID(id)
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) IsZero() bool {
	if id.IsZero() {
		return true
	}
	return false
}
