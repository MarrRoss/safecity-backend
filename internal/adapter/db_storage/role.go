package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type RoleRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewRoleRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*RoleRepositoryImpl, error) {
	return &RoleRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *RoleRepositoryImpl) AddRole(ctx context.Context, role *entity.Role) error {
	stmt := r.pg.Builder.Insert("roles").Columns(
		"id",
		"name",
		"created_at",
	).Values(
		role.ID,
		role.Name,
		role.CreatedAt,
	)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	// TODO: переделать
	return fmt.Errorf("failed to execute query: %w", err)
}

func (r *RoleRepositoryImpl) GetRole(ctx context.Context, id int) (*entity.Role, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"created_at").
		From("roles").
		Where("id = ?", id)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultRole response.RoleDB
	result, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer result.Close()

	err = r.pg.Scan.ScanOne(&resultRole, result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.observer.Logger.Trace().Err(err).Msg("role not found")
			return nil, adapter.ErrRoleNotFound
		} else if err != nil {
			r.observer.Logger.Error().Err(err).Msg("failed to scan result")
			return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
		}
	}

	role, err := response.RoleDbToEntity(&resultRole)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return role, nil
}

func (r *RoleRepositoryImpl) GetRoles(ctx context.Context) ([]*entity.Role, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"created_at").
		From("roles")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	// TODO: переделать
	var rolesDb []*response.RoleDB
	err = pgxscan.Select(ctx, r.pg.Pool, &rolesDb, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	roles, err := response.RoleDBListToEntityList(rolesDb)
	if err != nil {
		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
	}
	return roles, nil
}
