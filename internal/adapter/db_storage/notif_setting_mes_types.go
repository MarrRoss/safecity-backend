package db_storage

//type NotifySettingMesTypeRepositoryImpl struct {
//	pg       *postgres.Postgres
//	observer *observability.Observability
//}
//
//func NewNotifySettingMesTypeRepositoryImpl(
//	pg *postgres.Postgres,
//	observer *observability.Observability,
//) (*NotifySettingMesTypeRepositoryImpl, error) {
//	return &NotifySettingMesTypeRepositoryImpl{
//		pg:       pg,
//		observer: observer,
//	}, nil
//}
//
//func (r *NotifySettingMesTypeRepositoryImpl) GetNotifySettingsMesTypesBySettings(
//	ctx context.Context,
//	settingsIDs []string,
//) ([]*response.NotificationSettingMessengerTypeDB, error) {
//	stmt := r.pg.Builder.Select(
//		"nsmt.id_notify_setting",
//		"mt.id AS id_mes_type",
//		"mt.mes_type",
//	).
//		From("notify_settings_mes_types nsmt").
//		InnerJoin("messengers_types mt ON nsmt.id_mes_type = mt.id").
//		Where(squirrel.Eq{"nsmt.id_notify_setting": settingsIDs}).
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
//	var notifySettingsMesTypesDb []*response.NotificationSettingMessengerTypeDB
//	err = r.pg.Scan.ScanAll(&notifySettingsMesTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan results")
//		return nil, fmt.Errorf("failed to scan results: %w", err)
//	}
//
//	return notifySettingsMesTypesDb, nil
//}
//
//func (r *NotifySettingMesTypeRepositoryImpl) AddNotifySettingMesTypes(
//	ctx context.Context,
//	notSettingMesTypes []*entity.NotifySettingMesType,
//) error {
//	stmt := r.pg.Builder.Insert("notify_settings_mes_types").Columns(
//		"id",
//		"id_notify_setting",
//		"id_mes_type",
//		"created_at",
//		"ended_at",
//	)
//
//	for _, notSettingMesType := range notSettingMesTypes {
//		if notSettingMesType != nil {
//			stmt = stmt.Values(
//				notSettingMesType.ID.String(),
//				notSettingMesType.IDNotifySetting.String(),
//				notSettingMesType.IDMesType.String(),
//				notSettingMesType.CreatedAt,
//				notSettingMesType.EndedAt,
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
//func (r *NotifySettingMesTypeRepositoryImpl) DeleteNotifySettingMesTypesByNotifySetting(ctx context.Context,
//	id value_object.ID) error {
//	//stmt := r.pg.Builder.Delete("notify_settings_mes_types").Where("id_notify_setting = ?", id.String())
//	//
//	//sql, args, err := stmt.ToSql()
//	//if err != nil {
//	//	r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
//	//	return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	//}
//	//_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	//if err != nil {
//	//	r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//	//	return fmt.Errorf("failed to execute query: %w", err)
//	//}
//	return nil
//}
