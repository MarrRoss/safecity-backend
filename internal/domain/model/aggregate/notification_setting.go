package aggregate

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type NotificationSetting struct {
	ID               value_object.ID
	Frequency        *entity.NotificationFrequency
	EventType        string
	Receiver         *entity.User
	Sender           *entity.User
	AlertsHistory    []*entity.NotificationLog
	CreatedAt        time.Time
	UpdatedAt        time.Time
	EndedAt          *time.Time
	Zone             *entity.Zone
	BatteryThreshold *value_object.BatteryThreshold
}

func NewNotificationSetting(
	frequency *entity.NotificationFrequency,
	eventType string,
	receiver *entity.User,
	sender *entity.User,
	zone *entity.Zone,
	battery *value_object.BatteryThreshold,
) (*NotificationSetting, error) {
	id := value_object.NewID()
	if eventType == "" {
		return nil, fmt.Errorf("event type is empty: %v", domain.ErrInvalidEventType)
	}
	if receiver == nil {
		return nil, fmt.Errorf("receiver is empty: %v", domain.ErrUserNotFound)
	}
	if sender == nil {
		return nil, fmt.Errorf("sender is empty: %v", domain.ErrUserNotFound)
	}
	if (zone != nil && battery != nil) || (zone == nil && battery == nil) {
		return nil, fmt.Errorf("zone or battery are required: %v", domain.ErrInvalidZoneOrBattery)
	}
	newNotificationSetting := NotificationSetting{
		ID:               id,
		Frequency:        frequency,
		EventType:        eventType,
		Receiver:         receiver,
		Sender:           sender,
		AlertsHistory:    nil,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		EndedAt:          nil,
		Zone:             zone,
		BatteryThreshold: battery}
	return &newNotificationSetting, nil
}

func (n *NotificationSetting) ChangeFrequency(fr *entity.NotificationFrequency) error {
	if n.EndedAt != nil {
		return domain.ErrNotifySettingIsDeleted
	}
	if fr == nil {
		return domain.ErrInvalidFrequency
	}
	n.Frequency = fr
	n.UpdatedAt = time.Now()
	return nil
}

func (n *NotificationSetting) StopExistence() error {
	if n.EndedAt != nil {
		return domain.ErrNotifySettingIsDeleted
	}
	timeNow := time.Now()
	n.EndedAt = &timeNow
	n.UpdatedAt = time.Now()
	return nil
}

//func (n *NotificationSetting) ChangeNotificationTypes(notTypes []*entity.EventType) error {
//	if notTypes == nil {
//		return errors.New("notification types are nil")
//	}
//	n.EventType = notTypes
//	n.UpdatedAt = time.Now()
//	return nil
//}

//func (n *NotificationSetting) ChangeMessengerTypes(mesTypes []*entity.MessengerType) error {
//	if mesTypes == nil {
//		return errors.New("messenger types are nil")
//	}
//	n.MessengerTypes = mesTypes
//	n.UpdatedAt = time.Now()
//	return nil
//}
