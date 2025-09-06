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
	"time"
)

type UserLocationsRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewUserLocationsRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*UserLocationsRepositoryImpl, error) {
	return &UserLocationsRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *UserLocationsRepositoryImpl) AddUserLocation(
	ctx context.Context,
	location *entity.UserLocation,
) error {
	// выражение для geography(Point,4326)
	locExpr := squirrel.Expr(
		"ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography",
		float64(location.Location.Latitude),
		float64(location.Location.Longitude),
	)
	stmt := r.pg.Builder.
		Insert("locations_logs").
		Columns("id", "user_id", "location", "battery", "created_at").
		Values(
			location.ID.String(),
			location.UserID.String(),
			locExpr,
			int(location.Battery),
			location.CreatedAt,
		)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build SQL string: %w", err, adapter.ErrStorage)
	}
	if _, err = r.pg.Pool.Exec(ctx, sql, args...); err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	return nil
}

func (r *UserLocationsRepositoryImpl) AddLocationContext(
	ctx context.Context,
	locationLogID value_object.ID,
	zoneID *value_object.ID,
	notifyType string,
	battery *int,
) error {
	var zoneVal interface{}
	if zoneID == nil {
		zoneVal = nil
	} else {
		zoneVal = zoneID.String()
	}

	stmt := r.pg.Builder.
		Insert("locations_context").
		Columns("location_log_id",
			"zone_id",
			"notification_type",
			"battery",
			"created_at").
		Values(
			locationLogID.String(),
			zoneVal,
			notifyType,
			battery,
			time.Now().UTC(),
		)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build SQL string: %w", err, adapter.ErrStorage)
	}
	if _, err = r.pg.Pool.Exec(ctx, sql, args...); err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	return nil
}

//func (r *UserLocationsRepositoryImpl) GetLastUserLocation(ctx context.Context, userID value_object.ID) (*entity.UserLocation, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"user_id",
//		"location",
//		"created_at").
//		From("locations_logs").
//		Where(squirrel.Eq{"user_id": userID}).
//		OrderBy("created_at DESC").
//		Limit(1)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
//	}
//
//	var resultUserLocation response.UserLocationDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
//	}
//  defer rows.Close()
//
//	err = r.pg.Scan.ScanOne(&resultUser, rows)
//	if errors.Is(err, pgx.ErrNoRows) {
//		r.observer.Logger.Trace().Err(err).Msg("user not found")
//		return nil, adapter.ErrUserNotFound
//	} else if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
//	}
//
//	user, err := response.UserDbToEntity(&resultUser)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
//		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
//	}
//
//	return user, nil
//}

func (r *UserLocationsRepositoryImpl) GetLastLocationContext(
	ctx context.Context,
	userID value_object.ID,
) (*response.LocationContextDB, error) {
	stmt := r.pg.Builder.
		Select(
			"lc.location_log_id",
			"lc.zone_id",
			"lc.notification_type",
			"lc.created_at").
		From("locations_context lc").
		Join("locations_logs ll ON ll.id = lc.location_log_id").
		Where(squirrel.Eq{
			"ll.user_id": userID.String(),
		}).
		OrderBy("lc.created_at DESC").
		Limit(1)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("build SQL for last location context")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var logContext response.LocationContextDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&logContext, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("last user location not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	return &logContext, nil
}

func (r *UserLocationsRepositoryImpl) GetLastZoneContext(
	ctx context.Context,
	userID value_object.ID,
) (*response.LocationContextDB, error) {
	stmt := r.pg.Builder.
		Select(
			"lc.location_log_id",
			"lc.zone_id",
			"lc.notification_type",
			"lc.created_at",
		).
		From("locations_context lc").
		Join("locations_logs ll ON ll.id = lc.location_log_id").
		Where(squirrel.Eq{"ll.user_id": userID.String()}).
		Where("lc.zone_id IS NOT NULL").
		Where(squirrel.Eq{"lc.notification_type": []string{"entry", "inside", "exit"}}).
		OrderBy("lc.created_at DESC").
		Limit(1)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("build SQL for last location context")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var logContext response.LocationContextDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&logContext, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("last user location not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	return &logContext, nil
}

func (r *UserLocationsRepositoryImpl) FindLatestLocationsByUserIDs(
	ctx context.Context,
	userIDs []value_object.ID,
) ([]*response.UserLatestLocations, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	ids := make([]string, len(userIDs))
	for i, id := range userIDs {
		ids[i] = id.String()
	}

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sub := sb.
		Select(
			"DISTINCT ON (user_id) user_id",
			"created_at",
			"location",
		).
		From("locations_logs").
		Where(squirrel.Eq{"user_id": ids}).
		OrderBy("user_id", "created_at DESC")

	stmt := sb.
		Select(
			"u.id",                          // db:"id"
			"u.first_name",                  // db:"first_name"
			"u.last_name",                   // db:"last_name"
			"u.email",                       // db:"email"
			"u.username",                    // db:"username"
			"u.tracking",                    // db:"tracking"
			"ui.external_id AS telegram_id", // db:"telegram_id"
			"ul.created_at",                 // db:"created_at"  (лога)
			"u.updated_at",                  // db:"updated_at"  (пользователя)
			"ST_X(ul.location::geometry) AS latitude",  // db:"latitude"
			"ST_Y(ul.location::geometry) AS longitude", // db:"longitude"
		).
		FromSelect(sub, "ul").
		Join("users u ON u.id = ul.user_id").
		LeftJoin(
			"users_integrations ui ON ui.user_id = u.id AND ui.system_id = ?",
			1,
		).
		Where(squirrel.Eq{"u.tracking": true}).
		OrderBy("u.id")

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("build SQL for last location context")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var userLocations []*response.UserLatestLocations
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&userLocations, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("user locations not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan user locations rows")
		return nil, fmt.Errorf("%w: failed to scan result to variable", adapter.ErrStorage)
	}

	return userLocations, nil
}

//// 4) Выполняем запрос
//rows, err := r.pg.Pool.Query(ctx, sqlStr, args...)
//if err != nil {
//	r.observer.Logger.Error().
//		Err(err).
//		Msg("failed to execute FindLatestLocationsByUserIDs")
//	return nil, fmt.Errorf("%w: %v", adapter.ErrStorage, err)
//}
//defer rows.Close()
//
//// 5) Считываем в слайс моделей
//var result []UserLatestLocations
//for rows.Next() {
//	var rec UserLatestLocations
//	if err := rows.Scan(
//		&rec.UserId,
//		&rec.FirstName,
//		&rec.LastName,
//		&rec.Email,
//		&rec.Login,
//		&rec.Tracking,
//		&rec.TelegramID,
//		&rec.CreatedAt,
//		&rec.UpdatedAt,
//		&rec.Latitude,
//		&rec.Longitude,
//	); err != nil {
//		r.observer.Logger.Error().
//			Err(err).
//			Msg("failed to scan row in FindLatestLocationsByUserIDs")
//		return nil, fmt.Errorf("%w: %v", adapter.ErrStorage, err)
//	}
//	result = append(result, rec)
//}
//if err := rows.Err(); err != nil {
//	r.observer.Logger.Error().
//		Err(err).
//		Msg("rows iteration error in FindLatestLocationsByUserIDs")
//	return nil, fmt.Errorf("%w: %v", adapter.ErrStorage, err)
//}
//
//return result, nil

//
//err = r.pg.Scan.ScanOne(&logContext, rows)
//if errors.Is(err, pgx.ErrNoRows) {
//	r.observer.Logger.Trace().Err(err).Msg("last user location not found")
//	return nil, nil
//} else if err != nil {
//	r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//	return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
//}
//
//return &logContext, nil
