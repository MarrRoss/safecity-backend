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

type FamilyRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewFamilyRepositoryImpl(pg *postgres.Postgres,
	observer *observability.Observability) (*FamilyRepositoryImpl, error) {
	return &FamilyRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *FamilyRepositoryImpl) AddFamily(ctx context.Context, family *entity.Family) error {
	stmt := r.pg.Builder.Insert("families").Columns(
		"id",
		"name",
		"created_at",
		"updated_at",
		"ended_at",
	).Values(
		family.ID.String(),
		family.Name.String(),
		family.CreatedAt,
		family.UpdatedAt,
		family.EndedAt,
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

func (r *FamilyRepositoryImpl) GetFamily(ctx context.Context, id value_object.ID) (*entity.Family, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"created_at",
		"updated_at",
	).
		From("families").
		Where("id = ?", id.String()).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultFamily response.FamilyDB
	result, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer result.Close()

	err = r.pg.Scan.ScanOne(&resultFamily, result)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("family not found")
		return nil, adapter.ErrFamilyNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	family, err := response.FamilyDbToEntity(&resultFamily)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return family, nil
}

func (r *FamilyRepositoryImpl) GetFamilyZones(ctx context.Context, id value_object.ID) ([]*entity.Zone, error) {
	stmt := r.pg.Builder.
		Select(
			"zones.id",
			"zones.name",
			"zones.safety",
			"ST_AsText(zones.boundaries) AS boundaries",
			"zones.created_at",
			"zones.updated_at",
			"zones.family_id",
		).
		From("zones").
		Join("families ON zones.family_id = families.id").
		Where("zones.family_id = ?", id.String()).
		Where(squirrel.Eq{"zones.ended_at": nil, "families.ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL string for getting family zones with join")
		return nil, fmt.Errorf("%w: failed to build SQL string for getting family zones with join: %w", err, adapter.ErrStorage)
	}

	var zonesDB []*response.ZoneDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query for getting family zones with join")
		return nil, fmt.Errorf("%w: failed to execute query for getting family zones with join: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&zonesDB, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("zones not found for family")
		return nil, adapter.ErrZoneNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan zones result")
		return nil, fmt.Errorf("%w: failed to scan zones result: %w", err, adapter.ErrStorage)
	}

	zones, err := response.ZoneDbListToEntityList(zonesDB)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert zones DB list to entity list")
		return nil, fmt.Errorf("%w: failed to convert zones DB list to entity list: %w", err, adapter.ErrStorage)
	}

	return zones, nil
}

//func (r *FamilyRepositoryImpl) GetFamiliesByAuthorID(ctx context.Context,
//	id value_object.externalID) ([]*entity.Family, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"name",
//		"created_at",
//		"updated_at",
//		"ended_at").
//		From("families").
//		Where("author_id = ?", id.String())
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var familiesDb []*response.FamilyDB
//	err = pgxscan.Select(ctx, r.pg.Pool, &familiesDb, sql, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	families, err := response.FamilyDbListToEntityList(familiesDb)
//	if err != nil {
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//	return families, nil
//}

func (r *FamilyRepositoryImpl) GetFamiliesByIDs(ctx context.Context, ids []string) ([]*entity.Family, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"created_at",
		"updated_at",
	).
		From("families").
		Where(squirrel.Eq{"id": ids}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var familiesDb []*response.FamilyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&familiesDb, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("membership families not found")
		return nil, adapter.ErrMembershipFamiliesNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	families, err := response.FamilyDbListToEntityList(familiesDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return families, nil
}

func (r *FamilyRepositoryImpl) UpdateFamily(ctx context.Context, family *entity.Family) error {
	stmt := r.pg.Builder.Update("families").
		SetMap(map[string]interface{}{
			"name":       family.Name.String(),
			"updated_at": family.UpdatedAt,
			"ended_at":   family.EndedAt,
		}).
		Where(squirrel.Eq{"id": family.ID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

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

//func (r *FamilyRepositoryImpl) DeleteFamily(ctx context.Context, id value_object.externalID) error {
//	stmt := r.pg.Builder.Delete("families").
//		Where("id = ?", id.String())
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		return fmt.Errorf("failed to execute query: %w", err)
//	}
//	return nil
//}
