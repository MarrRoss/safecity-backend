package aggregate

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type InvitationToFamily struct {
	ID        value_object.ID
	Author    *entity.User
	Role      *entity.Role
	Family    *entity.Family
	CreatedAt time.Time
	EndedAt   *time.Time
}

func NewInvitationToFamily(
	author *entity.User,
	role *entity.Role,
	family *entity.Family,
) (*InvitationToFamily, error) {
	if author == nil {
		return nil, fmt.Errorf("author is empty: %w", domain.ErrUserNotFound)
	}
	if role == nil {
		return nil, fmt.Errorf("role is empty: %w", domain.ErrRoleNotFound)
	}
	if family == nil {
		return nil, fmt.Errorf("family is empty: %w", domain.ErrFamilyNotFound)
	}
	id := value_object.NewID()
	now := time.Now()
	newInvitationToFamily := InvitationToFamily{
		ID:        id,
		Author:    author,
		Role:      role,
		Family:    family,
		CreatedAt: now,
		EndedAt:   nil,
	}
	return &newInvitationToFamily, nil
}

func (inv *InvitationToFamily) StopExistence() error {
	if inv.EndedAt != nil {
		return domain.ErrInvitationToFamilyIsDeleted
	}
	timeNow := time.Now()
	inv.EndedAt = &timeNow
	return nil
}
