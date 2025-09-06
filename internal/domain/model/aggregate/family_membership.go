package aggregate

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type FamilyMembership struct {
	ID        value_object.ID
	User      *entity.User
	Role      *entity.Role
	Family    *entity.Family
	CreatedAt time.Time
	UpdatedAt time.Time
	EndedAt   *time.Time
}

func NewFamilyMembership(
	user *entity.User,
	role *entity.Role,
	family *entity.Family,
) (*FamilyMembership, error) {
	if user == nil {
		return nil, fmt.Errorf("user is empty: %w", domain.ErrUserNotFound)
	}
	if role == nil {
		return nil, fmt.Errorf("role is empty: %w", domain.ErrRoleNotFound)
	}
	if family == nil {
		return nil, fmt.Errorf("family is empty: %w", domain.ErrFamilyNotFound)
	}
	id := value_object.NewID()
	now := time.Now()
	newFamilyMembership := FamilyMembership{
		ID:        id,
		User:      user,
		Role:      role,
		Family:    family,
		CreatedAt: now,
		UpdatedAt: now,
		EndedAt:   nil}
	return &newFamilyMembership, nil
}
