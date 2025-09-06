package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
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

type IntegrationRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewIntegrationRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability) (*IntegrationRepositoryImpl, error) {
	return &IntegrationRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *IntegrationRepositoryImpl) AddIntegration(
	ctx context.Context,
	userID value_object.ID,
	systemID int,
) error {
	now := time.Now()
	stmt := r.pg.Builder.Insert("users_integrations").Columns(
		"user_id",
		"system_id",
		"external_id",
		"created_at",
		"updated_at",
	).Values(
		userID.String(),
		systemID,
		nil,
		now,
		now,
	)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}

	return nil
}

func (r *IntegrationRepositoryImpl) GetIntegration(
	ctx context.Context,
	userID value_object.ID,
	systemID int) (*response.IntegrationDB, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"user_id",
		"system_id",
		"external_id",
		"created_at",
		"updated_at").
		From("users_integrations").
		Where(squirrel.Eq{"user_id": userID.String()}).
		Where(squirrel.Eq{"system_id": systemID})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}
	var integration response.IntegrationDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()
	err = r.pg.Scan.ScanOne(&integration, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("integration not found")
		return nil, nil
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return &integration, nil
}

func (r *IntegrationRepositoryImpl) UpdateIntegration(
	ctx context.Context,
	userID value_object.ID,
	tgID string) error {
	stmt := r.pg.Builder.Update("users_integrations").
		SetMap(map[string]interface{}{
			"external_id": tgID,
			"updated_at":  time.Now(),
		}).
		Where(squirrel.Eq{"user_id": userID.String()}).
		Where(squirrel.Eq{"system_id": 1})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", adapter.ErrStorage, err)
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", adapter.ErrStorage, err)
	}

	return nil
}
