package adapter

import (
	"awesomeProjectDDD/internal/application"
	"fmt"
)

var ErrStorage = fmt.Errorf("storage error: %w", application.ErrApplication)
var ErrUserNotFound = fmt.Errorf("user not found: %w", ErrStorage)
var ErrMembershipsNotFound = fmt.Errorf("memberships not found: %w", ErrStorage)
var ErrMembershipUsersNotFound = fmt.Errorf("membership users not found: %w", ErrStorage)
var ErrMembershipFamiliesNotFound = fmt.Errorf("membership families not found: %w", ErrStorage)
var ErrFamilyNotFound = fmt.Errorf("family not found: %w", ErrStorage)
var ErrRoleNotFound = fmt.Errorf("role not found: %w", ErrStorage)
var ErrZoneNotFound = fmt.Errorf("zone not found: %w", ErrStorage)
var ErrMembershipNotFound = fmt.Errorf("membership not found: %w", ErrStorage)
var ErrInvitationNotFound = fmt.Errorf("invitation not found: %w", ErrStorage)
var ErrNotFamilyMembers = fmt.Errorf("not family members: %w", ErrStorage)
var ErrFrequencyNotFound = fmt.Errorf("frequency not found: %w", ErrStorage)
var ErrNotificationSettingExists = fmt.Errorf("notification setting already exists: %w", ErrStorage)
var ErrNotificationSettingNotFound = fmt.Errorf("notification setting not found: %w", ErrStorage)
var ErrIntegrationAlreadyExists = fmt.Errorf("integration already exists: %w", ErrStorage)
