package response

//type NotificationTypeDB struct {
//	Id         string `db:"id"`
//	NotifyType string `db:"notify_type"`
//}
//
//func NotificationTypeDbToEntity(notTypeDB *NotificationTypeDB) (*entity.EventType, error) {
//	id, err := value_object.NewIDFromString(notTypeDB.Id)
//	if err != nil {
//		return nil, fmt.Errorf("invalid id: %w", err)
//	}
//	return &entity.EventType{
//		ID:        id,
//		EventType: notTypeDB.NotifyType,
//	}, nil
//}
//
//func NotificationTypeDbListToEntityList(notTypesDb []*NotificationTypeDB) ([]*entity.EventType, error) {
//	var notifyTypes []*entity.EventType
//	for _, db := range notTypesDb {
//		notifyType, err := NotificationTypeDbToEntity(db)
//		if err != nil {
//			return nil, fmt.Errorf("failed to convert storage data to entity: %w", err)
//		}
//		notifyTypes = append(notifyTypes, notifyType)
//	}
//	return notifyTypes, nil
//}
