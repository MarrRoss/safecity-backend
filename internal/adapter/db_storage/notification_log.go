package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type NotificationLogRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewNotificationLogRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*NotificationLogRepositoryImpl, error) {
	return &NotificationLogRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *NotificationLogRepositoryImpl) AddNotificationLog(
	ctx context.Context,
	notificationLog *entity.NotificationLog) error {
	//r.observer.Logger.Info().Interface("notificationLog", notificationLog).Msg("AddNotificationLog")
	//r.observer.Logger.Info().Interface("t", notificationLog.CreatedAt.String()).Msg("AddNotificationLog")
	formattedTime := notificationLog.CreatedAt.Format("2006-01-02 15:04:05.999999-07:00")
	//r.observer.Logger.Info().Interface("t", formattedTime).Msg("AddNotificationLog")

	stmt := r.pg.Builder.Insert("notifications_logs").Columns(
		"id",
		"notification_id",
		"created_at",
		"context",
		"system_id",
	).Values(
		notificationLog.ID,
		notificationLog.NotificationSettingID,
		formattedTime,
		notificationLog.Context,
		notificationLog.SystemID,
	)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}
	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	return nil
}

//func (r *NotificationLogRepositoryImpl) GetUserLogs(
//	ctx context.Context,
//	notifySettingID value_object.ID) ([]*entity.NotificationLog, error) {
//	// TODO: переделать
//	stmt := r.pg.Builder.Select(
//		"id",
//		"notify_setting_id",
//		"created_at",
//		"context",
//	).From("А").Where("notify_setting_id = $1", notifySettingID)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	var notifyLogs []*entity.NotificationLog
//	for rows.Next() {
//		var notifyLog entity.NotificationLog
//		err = rows.Scan(&notifyLog.ID, &notifyLog.NotificationSettingID, &notifyLog.CreatedAt, &notifyLog.Context)
//		if err != nil {
//			return nil, fmt.Errorf("failed to scan rows: %w", err)
//		}
//		notifyLogs = append(notifyLogs, &notifyLog)
//	}
//	return notifyLogs, nil
//}

func (r *NotificationLogRepositoryImpl) GetLastLogByNotificationID(
	ctx context.Context,
	notifySettingID value_object.ID) (*entity.NotificationLog, error) {
	stmt := r.pg.Builder.
		Select(
			"id",
			"notification_id",
			"created_at",
			"context",
		).
		From("notifications_logs").
		Where("notification_id = ?", notifySettingID.String()).
		OrderBy("created_at DESC").
		Limit(1)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultLog response.NotificationLogDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&resultLog, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("log not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	log, err := response.NotificationLogDbToEntity(&resultLog)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}

	return log, nil
}

func (r *NotificationLogRepositoryImpl) GetLastLogsByNotificationIDs(
	ctx context.Context,
	ids []value_object.ID) ([]*response.NotificationLogDB, error) {
	stringIds := make([]string, len(ids))
	for i, id := range ids {
		stringIds[i] = id.String()
	}
	stmt := r.pg.Builder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"DISTINCT ON (notification_id) id",
			"notification_id",
			"created_at",
			"context",
		).
		From("notifications_logs").
		Where(squirrel.Eq{
			"notification_id": stringIds,
		}).
		OrderBy("notification_id", "created_at DESC")

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var notifyLogsDb []*response.NotificationLogDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&notifyLogsDb, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("notify logs not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return notifyLogsDb, nil
}
