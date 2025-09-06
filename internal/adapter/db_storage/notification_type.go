package db_storage

//type NotificationTypeRepositoryImpl struct {
//	pg       *postgres.Postgres
//	observer *observability.Observability
//}
//
//func NewNotificationTypeRepositoryImpl(
//	pg *postgres.Postgres,
//	observer *observability.Observability,
//) (*NotificationTypeRepositoryImpl, error) {
//	return &NotificationTypeRepositoryImpl{
//		pg:       pg,
//		observer: observer,
//	}, nil
//}
//
//func (r *NotificationTypeRepositoryImpl) GetNotificationType(
//	ctx context.Context,
//	id value_object.ID,
//) (*entity.EventType, error) {
//	stmt := r.pg.Builder.Select(
//		"notify_type").
//		From("notifications_types").
//		Where("id = ?", id).
//		Limit(1)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil,
//			fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var notType *entity.EventType
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&notType)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return notType, nil
//}
//
//func (r *NotificationTypeRepositoryImpl) GetAllNotificationTypes(
//	ctx context.Context) ([]*entity.EventType, error) {
//	stmt := r.pg.Builder.Select(
//		"id", "notify_type").
//		From("notifications_types")
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var notifyTypesDb []*response.NotificationTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&notifyTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	notifyTypes, err := response.NotificationTypeDbListToEntityList(notifyTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return notifyTypes, nil
//}
//
//func (r *NotificationTypeRepositoryImpl) GetNotificationTypesByIDs(
//	ctx context.Context,
//	ids []string,
//) ([]*entity.EventType, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"notify_type").
//		From("notifications_types").
//		Where(squirrel.Eq{"id": ids})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var notifyTypesDb []*response.NotificationTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&notifyTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	notifyTypes, err := response.NotificationTypeDbListToEntityList(notifyTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return notifyTypes, nil
//}
//
//func (r *NotificationTypeRepositoryImpl) GetNotificationTypesBySetting(
//	ctx context.Context,
//	id value_object.ID,
//) ([]*entity.EventType, error) {
//	stmt := r.pg.Builder.Select(
//		"nt.id",
//		"nt.notify_type").
//		From("notify_settings_notify_types AS nsnt").
//		Join("notifications_types AS nt ON nt.id = nsnt.id_notify_type").
//		Where("nsnt.id_notify_setting = ?", id)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var notifyTypesDb []*response.NotificationTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&notifyTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	notifyTypes, err := response.NotificationTypeDbListToEntityList(notifyTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return notifyTypes, nil
//}

//func (r *NotificationTypeRepositoryImpl) GetNotificationTypesIDsByNotificationSetting(ctx context.Context,
//	notifySetID value_object.externalID) ([]value_object.externalID, error) {
//	stmt := r.pg.Builder.Select(
//		"id_notify_type").
//		From("notify_settings_notify_types").
//		Where("id_notify_setting = ?", notifySetID)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil,
//			fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var notTypesIDs []value_object.externalID
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&notTypesIDs)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return notTypesIDs, nil
//}
