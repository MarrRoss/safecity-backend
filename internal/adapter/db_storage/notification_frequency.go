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

type NotificationFrequencyRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewNotificationFrequencyRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*NotificationFrequencyRepositoryImpl, error) {
	return &NotificationFrequencyRepositoryImpl{
		pg:       pg,
		observer: observer}, nil
}

func (r *NotificationFrequencyRepositoryImpl) GetFrequency(
	ctx context.Context,
	id value_object.ID,
) (*entity.NotificationFrequency, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"frequency").
		From("notifications_frequencies").
		Where(squirrel.Eq{"id": id.String()})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}
	var notifyFrequencyDb response.NotificationFrequencyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()
	err = r.pg.Scan.ScanOne(&notifyFrequencyDb, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("frequency not found")
		return nil, adapter.ErrFrequencyNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	notifyFrequency, err := response.NotificationFrequencyDbToEntity(&notifyFrequencyDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return notifyFrequency, nil
}

func (r *NotificationFrequencyRepositoryImpl) GetAllNotificationFrequencies(
	ctx context.Context) ([]*entity.NotificationFrequency, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"frequency").
		From("notifications_frequencies")

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	var notifyFrequenciesDb []*response.NotificationFrequencyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&notifyFrequenciesDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}

	notifyFrequencies, err := response.NotificationFrequencyDbListToEntityList(notifyFrequenciesDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
	}
	return notifyFrequencies, nil
}

func (r *NotificationFrequencyRepositoryImpl) GetNotificationFrequenciesByIDs(
	ctx context.Context,
	ids []string,
) ([]*entity.NotificationFrequency, error) {
	stmt := r.pg.Builder.
		Select("id", "frequency").
		From("notifications_frequencies").
		Where(squirrel.Eq{"id": ids})

	sqlQuery, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	var notifyFrequenciesDb []*response.NotificationFrequencyDB
	rows, err := r.pg.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&notifyFrequenciesDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}

	notifyFrequencies, err := response.NotificationFrequencyDbListToEntityList(notifyFrequenciesDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
	}

	return notifyFrequencies, nil
}
