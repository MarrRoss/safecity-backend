package db_storage

import (
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
)

type FamilyZoneRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewFamilyZoneRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*FamilyZoneRepositoryImpl, error) {
	return &FamilyZoneRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

//func (r *FamilyZoneRepositoryImpl) AddFamilyZone(ctx context.Context, familyZone *entity.FamilyZone) error {
//	stmt := r.pg.Builder.Insert("families_zones").Columns(
//		"id",
//		"family_id",
//		"zone_id",
//		"created_at",
//	).Values(
//		familyZone.externalID.String(),
//		familyZone.FamilyID.String(),
//		familyZone.ZoneID.String(),
//		familyZone.CreatedAt,
//	)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
//		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
//	}
//
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
//	}
//
//	return nil
//}

//func (r *FamilyZoneRepositoryImpl) GetFamilyZones(ctx context.Context, familyID value_object.externalID) ([]*entity.Zone, error) {
//	stmt := r.pg.Builder.Select(
//		"z.id",
//		"z.z_name",
//		"z.author_id",
//		"z.safety",
//		"ST_AsText(z.boundaries) AS boundaries",
//		"z.created_at",
//		"z.updated_at",
//	).
//		From("zones z").
//		Join("families_zones fz ON fz.zone_id = z.id").
//		Where(squirrel.Eq{"fz.family_id": familyID.String()})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
//	}
//
//	var zonesDb []*response.ZoneDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
//	}
//	defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&zonesDb, rows)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
//		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
//	}
//
//	zones, err := response.ZoneDbListToEntityList(zonesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
//		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
//	}
//	return zones, nil
//}

//func (r *FamilyZoneRepositoryImpl) DeleteFamilyZone(
//	ctx context.Context,
//	familyID value_object.externalID,
//	zoneID value_object.externalID,
//) error {
//	stmt := r.pg.Builder.Delete("families_zones").
//		Where("family_id = ? AND zone_id = ?", familyID.String(), zoneID.String())
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
//		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
//	}
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
//	}
//	return nil
//}

//func (r *FamilyZoneRepositoryImpl) DeleteFamilyZones(ctx context.Context,
//	familyID value_object.externalID) error {
//	stmt := r.pg.Builder.Delete("families_zones").
//		Where("family_id = ?", familyID.String())
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	return fmt.Errorf("failed to execute query: %w", err)
//}
