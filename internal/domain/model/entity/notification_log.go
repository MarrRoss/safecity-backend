package entity

import (
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)
import _ "time/tzdata"

type NotificationLog struct {
	ID                    value_object.ID
	NotificationSettingID value_object.ID
	CreatedAt             time.Time
	Context               string
	SystemID              int
}

func NewNotificationLog(
	settingID value_object.ID,
	//createdAt time.Time,
	context string,
	systemID int,
) (*NotificationLog, error) {
	id := value_object.NewID()
	if len(context) == 0 {
		return nil, fmt.Errorf("context is empty: %v", domain.ErrInvalidNotificationLogContext)
	}
	//init the loc
	loc, _ := time.LoadLocation("Europe/Moscow")
	//set timezone,
	now := time.Now().In(loc)
	//now := time.Now()
	//sixMonthsAgo := now.AddDate(0, -6, 0)
	//if createdAt.Before(sixMonthsAgo) {
	//	return nil, fmt.Errorf("time must be within the last six months: %v", domain.ErrInvalidSendTime)
	//}
	newNotificationLog := NotificationLog{
		ID:                    id,
		NotificationSettingID: settingID,
		CreatedAt:             now,
		Context:               context,
		SystemID:              systemID}
	return &newNotificationLog, nil
}
