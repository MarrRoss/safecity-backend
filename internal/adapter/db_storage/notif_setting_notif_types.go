package db_storage

//type NotifySettingNotifyTypeRepositoryImpl struct {
//	pg       *postgres.Postgres
//	observer *observability.Observability
//}
//
//func NewNotifySettingNotifyTypeRepositoryImpl(
//	pg *postgres.Postgres,
//	observer *observability.Observability,
//) (*NotifySettingNotifyTypeRepositoryImpl, error) {
//	return &NotifySettingNotifyTypeRepositoryImpl{
//		pg:       pg,
//		observer: observer,
//	}, nil
//}
//
////func (r *NotifySettingNotifyTypeRepositoryImpl) AddNotifySettingNotifyType(ctx context.Context,
////	notSettingNotType *entity.NotifySettingNotifyType) error {
////	stmt := r.pg.Builder.Insert("notify_settings_notify_types").Columns(
////		"id",
////		"id_notify_setting",
////		"id_notify_type",
////		"created_at",
////		"ended_at",
////	).Values(
////		notSettingNotType.externalID,
////		notSettingNotType.IDNotifySetting,
////		notSettingNotType.IDNotifyType,
////		notSettingNotType.CreatedAt,
////		notSettingNotType.EndedAt,
////	)
////
////	sql, args, err := stmt.ToSql()
////	if err != nil {
////		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
////	}
////
////	_, err = r.pg.Pool.Exec(ctx, sql, args...)
////	// переделать
////	return fmt.Errorf("failed to execute query: %w", err)
////}
//
//func (r *NotifySettingNotifyTypeRepositoryImpl) AddNotifySettingNotTypes(
//	ctx context.Context,
//	notSettingNotTypes []*entity.NotifySettingNotifyType,
//) error {
//	stmt := r.pg.Builder.Insert("notify_settings_notify_types").Columns(
//		"id",
//		"id_notify_setting",
//		"id_notify_type",
//		"created_at",
//		"ended_at",
//	)
//
//	for _, notSettingNotType := range notSettingNotTypes {
//		if notSettingNotType != nil {
//			stmt = stmt.Values(
//				notSettingNotType.ID.String(),
//				notSettingNotType.IDNotifySetting.String(),
//				notSettingNotType.IDNotifyType.String(),
//				notSettingNotType.CreatedAt,
//				notSettingNotType.EndedAt,
//			)
//		}
//	}
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
//		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	return nil
//}
//
////func (r *NotifySettingNotifyTypeRepositoryImpl) GetNotTypesIDsByNotSetting(ctx context.Context,
////	id value_object.externalID) ([]int, error) {
////	stmt := r.pg.Builder.Select(
////		"id_notif_type").
////		From("notif_settings_notif_types").
////		Where(squirrel.Eq{"id_notif_setting": id.String()})
////
////	sql, args, err := stmt.ToSql()
////	if err != nil {
////		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
////	}
////
////	var typesIDs []int
////	err = pgxscan.Select(ctx, r.pg.Pool, &typesIDs, sql, args...)
////	if err != nil {
////		return nil, fmt.Errorf("failed to execute query: %w", err)
////	}
////
////	return typesIDs, nil
////}
//
//func (r *NotifySettingNotifyTypeRepositoryImpl) GetNotifySettingsNotifyTypesBySettings(
//	ctx context.Context,
//	settingsIDs []string,
//) ([]*response.NotificationSettingNotificationTypeDB, error) {
//	stmt := r.pg.Builder.Select(
//		"nsnt.id_notify_setting",
//		"nt.id AS id_notify_type",
//		"nt.notify_type").
//		From("notify_settings_notify_types nsnt").
//		InnerJoin("notifications_types nt ON nsnt.id_notify_type = nt.id").
//		Where(squirrel.Eq{"nsnt.id_notify_setting": settingsIDs}).
//		Where(squirrel.Eq{"ended_at": nil})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL string")
//		return nil, fmt.Errorf("failed to build SQL string: %w", err)
//	}
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	var notifySettingsNotifyTypesDb []*response.NotificationSettingNotificationTypeDB
//	err = r.pg.Scan.ScanAll(&notifySettingsNotifyTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan results")
//		return nil, fmt.Errorf("failed to scan results: %w", err)
//	}
//
//	return notifySettingsNotifyTypesDb, nil
//}

//func (r *NotifySettingNotifyTypeRepositoryImpl) DeleteNotifySettingNotifyTypesByNotifySetting(
//	ctx context.Context,
//	id value_object.externalID,
//) error {
//	stmt := r.pg.Builder.Delete("notify_settings_notify_types").
//		Where(squirrel.Eq{"id_notify_setting": id.String()})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
//		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return fmt.Errorf("failed to execute query: %w", err)
//	}
//	return nil
//}
