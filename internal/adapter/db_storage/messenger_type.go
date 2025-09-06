package db_storage

//type MessengerTypeRepositoryImpl struct {
//	pg       *postgres.Postgres
//	observer *observability.Observability
//}
//
//func NewMessengerTypeRepositoryImpl(
//	pg *postgres.Postgres,
//	observer *observability.Observability,
//) (*MessengerTypeRepositoryImpl, error) {
//	return &MessengerTypeRepositoryImpl{
//		pg:       pg,
//		observer: observer,
//	}, nil
//}

//func (r *MessengerTypeRepositoryImpl) GetMessengerType(ctx context.Context,
//	id value_object.ID) (*entity.MessengerType, error) {
//	stmt := r.pg.Builder.Select(
//		"mes_type").
//		From("messengers_types").
//		Where("id = ?", id)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil,
//			fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var mesType *entity.MessengerType
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&mesType)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return mesType, nil
//}
//
//func (r *MessengerTypeRepositoryImpl) GetAllMessengerTypes(
//	ctx context.Context) ([]*entity.MessengerType, error) {
//	stmt := r.pg.Builder.Select(
//		"id", "mes_type").
//		From("messengers_types")
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var mesTypesDb []*response.MessengerTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&mesTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	mesTypes, err := response.MessengerTypeDbListToEntityList(mesTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return mesTypes, nil
//}
//
//func (r *MessengerTypeRepositoryImpl) GetMessengerTypesByIDs(
//	ctx context.Context,
//	ids []string,
//) ([]*entity.MessengerType, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"mes_type").
//		From("messengers_types").
//		Where(squirrel.Eq{"id": ids})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var mesTypesDb []*response.MessengerTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&mesTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	mesTypes, err := response.MessengerTypeDbListToEntityList(mesTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return mesTypes, nil
//}
//
//func (r *MessengerTypeRepositoryImpl) GetMessengerTypesBySetting(
//	ctx context.Context,
//	id value_object.ID,
//) ([]*entity.MessengerType, error) {
//	stmt := r.pg.Builder.Select(
//		"mt.id",
//		"mt.mes_type").
//		From("notify_settings_mes_types AS nsnt").
//		Join("messengers_types AS mt ON mt.id = nsnt.id_mes_type").
//		Where("nsnt.id_notify_setting = ?", id)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var mesTypesDb []*response.MessengerTypeDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&mesTypesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//
//	mesTypes, err := response.MessengerTypeDbListToEntityList(mesTypesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return mesTypes, nil
//}
