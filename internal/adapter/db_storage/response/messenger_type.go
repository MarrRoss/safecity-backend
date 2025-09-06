package response

//type MessengerTypeDB struct {
//	Id   string `db:"id"`
//	Type string `db:"mes_type"`
//}

//func MessengerTypeDbToEntity(db *MessengerTypeDB) (*entity.MessengerType, error) {
//	id, err := value_object.NewIDFromString(db.Id)
//	if err != nil {
//		return nil, err
//	}
//	return &entity.MessengerType{
//		ID:      id,
//		MesType: db.Type,
//	}, nil
//}

//func MessengerTypeDbListToEntityList(
//	dbList []*MessengerTypeDB) ([]*entity.MessengerType, error) {
//	var messengerTypes []*entity.MessengerType
//	for _, db := range dbList {
//		mesType, err := MessengerTypeDbToEntity(db)
//		if err != nil {
//			return nil, fmt.Errorf("failed to convert storage data to entity")
//		}
//		messengerTypes = append(messengerTypes, mesType)
//	}
//	return messengerTypes, nil
//}
