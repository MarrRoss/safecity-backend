package application

import (
	"errors"
	"fmt"
)

var ErrGeneral = errors.New("general error")
var ErrApplication = fmt.Errorf("application error: %w", ErrGeneral)
var ErrPhoneExists = fmt.Errorf("phone already exists: %w", ErrApplication)
var ErrEmailExists = fmt.Errorf("email already exists: %w", ErrApplication)
var ErrMembershipUserNotFound = fmt.Errorf("membership user not found: %w", ErrApplication)
var ErrMembershipFamilyNotFound = fmt.Errorf("membership family not found: %w", ErrApplication)
var ErrMembershipExists = fmt.Errorf("membership already exists: %w", ErrApplication)
var ErrUserIsNotAdult = fmt.Errorf("user is not adult: %w", ErrApplication)
var ErrReceiverIsSender = fmt.Errorf("receiver is sender: %w", ErrApplication)
var ErrDifferentArrayLength = fmt.Errorf("different array length: %w", ErrApplication)
var ErrNotifySettingsNotFound = fmt.Errorf("notify settings not found: %w", ErrApplication)
var ErrNotifyFrequenciesNotFound = fmt.Errorf("notify frequencies not found: %w", ErrApplication)
var ErrNotifySendersNotFound = fmt.Errorf("notify senders not found: %w", ErrApplication)
var ErrNotifyZoneNotFound = fmt.Errorf("notify zone not found: %w", ErrApplication)
var ErrInvalidNotifySettingDetail = fmt.Errorf("invalid notify setting detail: %w", ErrApplication)
var ErrZoneOverlapsZones = fmt.Errorf("zone overlaps other zones: %w", ErrApplication)
var ErrUserNotTracked = fmt.Errorf("user is not tracked: %w", ErrApplication)
