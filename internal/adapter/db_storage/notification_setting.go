package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"time"
)

type NotificationSettingRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewNotificationSettingRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*NotificationSettingRepositoryImpl, error) {
	return &NotificationSettingRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *NotificationSettingRepositoryImpl) AddNotificationSetting(ctx context.Context,
	setting *aggregate.NotificationSetting) error {
	var freqValue interface{} = nil
	if setting.Frequency != nil {
		freqValue = setting.Frequency.ID.String()
	}
	stmt := r.pg.Builder.Insert("notifications_settings").
		Columns(
			"id",
			"receiver_id",
			"sender_id",
			"event_type",
			"frequency_id",
			"created_at",
			"updated_at",
			"ended_at").
		Values(
			setting.ID.String(),
			setting.Receiver.ID.String(),
			setting.Sender.ID.String(),
			setting.EventType,
			freqValue,
			setting.CreatedAt,
			setting.UpdatedAt,
			setting.EndedAt,
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

func (r *NotificationSettingRepositoryImpl) AddZoneNotificationSetting(
	ctx context.Context,
	id, notificationID, zoneID value_object.ID) error {
	now := time.Now()
	stmt := r.pg.Builder.Insert("zones_notifications_details").
		Columns(
			"id",
			"notification_id",
			"zone_id",
			"created_at",
			"ended_at").
		Values(
			id.String(),
			notificationID.String(),
			zoneID.String(),
			now,
			nil,
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

func (r *NotificationSettingRepositoryImpl) AddBatteryNotificationSetting(
	ctx context.Context,
	id, notificationID value_object.ID,
	battery int) error {
	now := time.Now()
	stmt := r.pg.Builder.Insert("battery_notifications_details").
		Columns(
			"id",
			"notification_id",
			"battery_threshold",
			"created_at",
			"ended_at").
		Values(
			id.String(),
			notificationID.String(),
			battery,
			now,
			nil,
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

func (r *NotificationSettingRepositoryImpl) NotificationSettingExists(
	ctx context.Context,
	setting *aggregate.NotificationSetting,
) (bool, error) {
	qb := r.pg.Builder.
		Select("ns.id").
		From("notifications_settings ns").
		Where(squirrel.Eq{
			"ns.receiver_id": setting.Receiver.ID.String(),
			"ns.sender_id":   setting.Sender.ID.String(),
			"ns.event_type":  setting.EventType,
			"ns.ended_at":    nil,
		})

	if setting.EventType == "zone" && setting.Zone != nil {
		qb = qb.
			Join("zones_notifications_details zd ON ns.id = zd.notification_id").
			Where(squirrel.Eq{
				"zd.zone_id":  setting.Zone.ID.String(),
				"zd.ended_at": nil,
			})
	} else if setting.EventType == "battery" && setting.BatteryThreshold != nil {
		qb = qb.
			Join("battery_notifications_details bd ON ns.id = bd.notification_id").
			Where(squirrel.Eq{
				"bd.ended_at": nil,
			})
	} else {
		return false, domain.ErrInvalidZoneOrBattery
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return false, fmt.Errorf("%w: failed to build SQL query", adapter.ErrStorage)
	}

	var id string
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Если подписки не найдено, возвращаем false
			return false, nil
		}
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return false, fmt.Errorf("%w: failed to execute query", adapter.ErrStorage)
	}

	return true, adapter.ErrNotificationSettingExists
}

func (r *NotificationSettingRepositoryImpl) GetNotificationSetting(
	ctx context.Context,
	id value_object.ID,
) (*response.NotificationSettingDB, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"receiver_id",
		"sender_id",
		"event_type",
		"frequency_id",
		"created_at",
		"updated_at").
		From("notifications_settings").
		Where("id = ?", id.String()).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultSetting response.NotificationSettingDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()
	err = r.pg.Scan.ScanOne(&resultSetting, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("notification setting not found")
		return nil, adapter.ErrNotificationSettingNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return &resultSetting, nil
}

func (r *NotificationSettingRepositoryImpl) GetNotificationSettingsByZoneForUser(
	ctx context.Context,
	receiverID, zoneID value_object.ID,
) ([]*response.NotificationSettingMinBatteryDB, error) {
	stmt := r.pg.Builder.
		Select(
			"ns.id",
			"ns.receiver_id",
			"ns.sender_id",
			"ns.event_type",
			"ns.frequency_id",
			"ns.created_at",
			"ns.updated_at",
			"bd.battery_threshold",
		).
		From("notifications_settings ns").
		Join("zones_notifications_details zd ON ns.id = zd.notification_id AND zd.ended_at IS NULL").
		LeftJoin("battery_notifications_details bd ON bd.notification_id = ns.id AND bd.ended_at IS NULL").
		Where(squirrel.Eq{
			"ns.receiver_id": receiverID.String(),
			"zd.zone_id":     zoneID.String(),
			"ns.ended_at":    nil,
			"ns.event_type":  "zone",
		})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query for getting notification settings by zone")
		return nil, fmt.Errorf("%w: failed to build SQL query", adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query for getting notification settings by zone")
		return nil, fmt.Errorf("%w: failed to execute query", adapter.ErrStorage)
	}
	defer rows.Close()

	var settings []*response.NotificationSettingMinBatteryDB
	err = r.pg.Scan.ScanAll(&settings, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("notification settings not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan notification settings rows")
		return nil, fmt.Errorf("%w: failed to scan result to variable", adapter.ErrStorage)
	}

	return settings, nil
}

func (r *NotificationSettingRepositoryImpl) GetDangerZoneNotificationSettings(
	ctx context.Context,
	senderID value_object.ID,
	zoneID value_object.ID,
) ([]value_object.ID, error) {
	stmt := r.pg.Builder.
		Select("ns.id").
		From("notifications_settings ns").
		Join("zones_notifications_details znd ON ns.id = znd.notification_id").
		Where(squirrel.Eq{
			"ns.sender_id":  senderID.String(),
			"ns.event_type": "zone",
			"ns.ended_at":   nil,
			"znd.zone_id":   zoneID.String(),
			"znd.ended_at":  nil,
		})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	var strIDs []string
	err = r.pg.Scan.ScanAll(&strIDs, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("notification settings not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan notification settings rows")
		return nil, fmt.Errorf("%w: failed to scan result to variable", adapter.ErrStorage)
	}

	var settingIDs []value_object.ID
	for _, s := range strIDs {
		id, err := value_object.NewIDFromString(s)
		if err != nil {
			r.observer.Logger.Error().Err(err).Msg("failed to parse id from db string")
			return nil, fmt.Errorf("%w: failed to parse id from db string: %v", domain.ErrInvalidID, err)
		}
		settingIDs = append(settingIDs, id)
	}

	return settingIDs, nil

	//for rows.Next() {
	//	var id string
	//	if err := rows.Scan(&id); err != nil {
	//		r.observer.Logger.Error().Err(err).Msg("failed to scan notification id")
	//		return nil, fmt.Errorf("failed to scan notification id: %w", adapter.ErrStorage)
	//	}
	//	settingID, err := value_object.NewIDFromString(id)
	//	if err != nil {
	//		r.observer.Logger.Error().Err(err).Msg("failed to parse id from db string")
	//		return nil, fmt.Errorf("failed to parse id from db string: %w", domain.ErrInvalidID)
	//	}
	//	settingIDs = append(settingIDs, settingID)
	//}
	//return settingIDs, nil
}

func (r *NotificationSettingRepositoryImpl) GetNotificationSettingsByIDs(
	ctx context.Context,
	ids []value_object.ID) ([]*response.NotificationSettingWithFrequencyDB, error) {
	stringIds := make([]string, len(ids))
	for i, id := range ids {
		stringIds[i] = id.String()
	}
	builder := r.pg.Builder.PlaceholderFormat(squirrel.Dollar)
	stmt := builder.Select(
		"ns.id",
		"ns.receiver_id",
		"ns.sender_id",
		"ns.event_type",
		"nf.frequency",
		"ns.created_at",
		"ns.updated_at").
		From("notifications_settings ns").
		Join("notifications_frequencies nf ON ns.frequency_id = nf.id").
		Where(squirrel.Eq{"ns.id": stringIds}).
		Where(squirrel.Eq{"ns.ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var notifySettingsDb []*response.NotificationSettingWithFrequencyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&notifySettingsDb, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("notify settings not found")
		return nil, adapter.ErrNotificationSettingNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return notifySettingsDb, nil
}

//func (r *NotificationSettingRepositoryImpl) GetNotificationSettingIDByReceiverSender(ctx context.Context,
//	receiverID, senderID value_object.externalID) (value_object.externalID, error) {
//	stmt := r.pg.Builder.Select("id").
//		From("notifications_settings").
//		Where("receiver_id = ? AND sender_id = ?", receiverID.String(), senderID.String()).
//		Limit(1)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return value_object.externalID{},
//			fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var id value_object.externalID
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&id)
//	if err != nil {
//		return value_object.externalID{}, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return id, nil
//}

func (r *NotificationSettingRepositoryImpl) GetNotificationsSettingsIDsByReceiver(ctx context.Context,
	receiverID value_object.ID) ([]value_object.ID, error) {
	stmt := r.pg.Builder.Select("id").
		From("notifications_settings").
		Where("receiver_id = ?", receiverID.String()).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	var ids []value_object.ID
	result, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer result.Close()
	err = r.pg.Scan.ScanAll(&ids, result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}
	return ids, nil
}

func (r *NotificationSettingRepositoryImpl) GetIntegratedReceiversBySettingIDs(
	ctx context.Context,
	settingIDs []value_object.ID,
) ([]*response.UserIntegration, error) {
	stringIds := make([]string, len(settingIDs))
	for i, id := range settingIDs {
		stringIds[i] = id.String()
	}

	builder := r.pg.Builder.PlaceholderFormat(squirrel.Dollar)
	stmt := builder.
		Select(
			"ns.receiver_id   AS user_id",
			"ns.id            AS notify_setting_id",
			`ui.system_id     AS "system.id"`,
			`ui.external_id   AS "system.external_id"`,
		).
		From("notifications_settings ns").
		Join("users_integrations ui ON ui.user_id = ns.receiver_id").
		Where(squirrel.Eq{
			"ns.id":       stringIds,
			"ns.ended_at": nil,
		}).
		Where(squirrel.Expr("ui.external_id IS NOT NULL"))

	// AND ui.system_id = ?

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL for integrated receivers")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var integratedUsersDb []*response.UserIntegration
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&integratedUsersDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return integratedUsersDb, nil
}

//func (r *NotificationSettingRepositoryImpl) GetNotificationSettingSenderID(ctx context.Context,
//	id value_object.externalID) (value_object.externalID, error) {
//	stmt := r.pg.Builder.Select(
//		"sender_id").
//		From("notifications_settings").
//		Where("id = ?", id.String())
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return value_object.externalID{},
//			fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var senderID value_object.externalID
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&senderID)
//	if err != nil {
//		return value_object.externalID{}, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return senderID, nil
//}

func (r *NotificationSettingRepositoryImpl) GetNotificationSettingsByZoneReceiverAndFamily(
	ctx context.Context,
	userID,
	familyID,
	zoneID value_object.ID,
) ([]*response.NotificationSettingDB, error) {
	stmt := r.pg.Builder.Select(
		"ns.id as id",
		"ns.frequency_id as frequency_id",
		"ns.receiver_id as receiver_id",
		"ns.sender_id as sender_id",
		"ns.created_at as created_at",
		"ns.updated_at as updated_at",
		"ns.zone_id as zone_id",
		"ns.family_id as family_id").
		From("notifications_settings ns").
		Join("notifications_frequencies f ON ns.frequency_id = f.id").
		Join("users sender ON ns.sender_id = sender.id").
		Where(squirrel.Eq{
			"ns.receiver_id": userID.String(),
			"ns.zone_id":     zoneID.String(),
			"ns.family_id":   familyID.String(),
		}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	var notificationSettingsDb []*response.NotificationSettingDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&notificationSettingsDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}
	return notificationSettingsDb, nil
}

func (r *NotificationSettingRepositoryImpl) GetBatteryNotificationSettingsByChild(
	ctx context.Context,
	childID value_object.ID,
) ([]*response.BatterySettingWithFrequencyDB, error) {
	stmt := r.pg.Builder.
		Select(`
            ns.id AS notification_id,
            bd.battery_threshold,
            nf.frequency`).
		From("notifications_settings ns").
		Join("battery_notifications_details bd ON bd.notification_id = ns.id").
		Join("notifications_frequencies nf ON nf.id = ns.frequency_id").
		Where(squirrel.Eq{
			"ns.sender_id":  childID.String(),
			"ns.event_type": "battery",
			"ns.ended_at":   nil,
			"bd.ended_at":   nil,
		})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var settings []*response.BatterySettingWithFrequencyDB
	if err := r.pg.Scan.ScanAll(&settings, rows); err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}
	return settings, nil
}

func (r *NotificationSettingRepositoryImpl) GetBatteryNotificationSettingsByReceiver(
	ctx context.Context,
	receiverID value_object.ID,
) ([]*response.NotificationSettingMinBatteryDB, error) {
	stmt := r.pg.Builder.
		Select(
			"ns.id",
			"ns.receiver_id",
			"ns.sender_id",
			"ns.event_type",
			"ns.frequency_id",
			"ns.created_at",
			"ns.updated_at",
			"bd.battery_threshold AS battery_threshold",
		).
		From("notifications_settings ns").
		Join("battery_notifications_details bd ON bd.notification_id = ns.id AND bd.ended_at IS NULL").
		Where(squirrel.Eq{
			"ns.receiver_id": receiverID.String(),
			"ns.event_type":  "battery",
			"ns.ended_at":    nil,
		})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var settings []*response.NotificationSettingMinBatteryDB
	if err := r.pg.Scan.ScanAll(&settings, rows); err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}
	return settings, nil
}

func (r *NotificationSettingRepositoryImpl) FindLocationSenderIDsByReceiver(
	ctx context.Context,
	receiverID value_object.ID,
) ([]value_object.ID, error) {
	subQ := r.pg.Builder.
		Select("fm2.family_id").
		From("families_memberships fm2").
		Join("families f ON fm2.family_id = f.id").
		Where(squirrel.Eq{"fm2.user_id": receiverID.String()}).
		Where(squirrel.Eq{"f.ended_at": nil})
	subQText, subQArgs, err := subQ.PlaceholderFormat(squirrel.Question).ToSql()

	mainQ := squirrel.
		Select("DISTINCT fm.user_id").
		From("families_memberships fm").
		Where(squirrel.NotEq{"fm.user_id": receiverID.String()}).
		Where("fm.family_id IN ("+subQText+")", subQArgs...).
		PlaceholderFormat(squirrel.Dollar)

	sqlStr, args, err := mainQ.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Query(ctx, sqlStr, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	var strIDs []string
	err = r.pg.Scan.ScanAll(&strIDs, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("location senders not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan location senders rows")
		return nil, fmt.Errorf("%w: failed to scan result to variable", adapter.ErrStorage)
	}

	var sendersIDs []value_object.ID
	for _, s := range strIDs {
		id, err := value_object.NewIDFromString(s)
		if err != nil {
			r.observer.Logger.Error().Err(err).Msg("failed to parse id from db string")
			return nil, fmt.Errorf("%w: failed to parse id from db string: %v", domain.ErrInvalidID, err)
		}
		sendersIDs = append(sendersIDs, id)
	}

	return sendersIDs, nil

	//var senderIDs []uuid.UUID
	//for rows.Next() {
	//	var id uuid.UUID
	//	if err := rows.Scan(&id); err != nil {
	//		r.observer.Logger.Error().
	//			Err(err).
	//			Msg("failed to scan sender_id")
	//		return nil, fmt.Errorf(
	//			"%w: failed to scan sender_id: %v",
	//			adapter.ErrStorage, err,
	//		)
	//	}
	//	senderIDs = append(senderIDs, id)
	//}
	//if err := rows.Err(); err != nil {
	//	r.observer.Logger.Error().
	//		Err(err).
	//		Msg("rows iteration error in FindSenderIDsByReceiver")
	//	return nil, fmt.Errorf(
	//		"%w: rows iteration error: %v",
	//		adapter.ErrStorage, err,
	//	)
	//}
	//
	//return senderIDs, nil
	//
	//stmt := r.pg.Builder.
	//	Select("ns.id").
	//	From("notifications_settings ns").
	//	Join("zones_notifications_details znd ON ns.id = znd.notification_id").
	//	Where(squirrel.Eq{
	//		"ns.sender_id":  senderID.String(),
	//		"ns.event_type": "zone",
	//		"ns.ended_at":   nil,
	//		"znd.zone_id":   zoneID.String(),
	//		"znd.ended_at":  nil,
	//	})
}

func (r *NotificationSettingRepositoryImpl) UpdateNotificationSetting(
	ctx context.Context,
	setting *aggregate.NotificationSetting,
) error {
	stmt := r.pg.Builder.Update("notifications_settings").
		SetMap(map[string]interface{}{
			"frequency_id": setting.Frequency.ID,
			"updated_at":   setting.UpdatedAt,
			"ended_at":     setting.EndedAt,
		}).
		Where(squirrel.Eq{"id": setting.ID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("failed to build SQL string from the query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

//func (r *NotificationSettingRepositoryImpl) DeleteNotificationSetting(
//	ctx context.Context,
//	id value_object.externalID,
//) error {
//	stmt := r.pg.Builder.Delete("notifications_settings").Where("id = ?", id.String())
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
