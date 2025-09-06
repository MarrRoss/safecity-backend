package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"strings"
)

type ZoneRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewZoneRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability,
) (*ZoneRepositoryImpl, error) {
	return &ZoneRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *ZoneRepositoryImpl) AddZone(ctx context.Context, zone *entity.Zone) error {
	polygonWKT := ConvertBoundariesToPolygon(zone.Boundaries)
	if polygonWKT == "" {
		r.observer.Logger.Trace().Msg("invalid boundaries: not enough points or invalid polygon")
		return fmt.Errorf("invalid boundaries: %w", domain.ErrInvalidZoneBoundaries)
	}

	stmt := r.pg.Builder.Insert("zones").Columns(
		"id",
		"name",
		"safety",
		"boundaries",
		"created_at",
		"updated_at",
		"ended_at",
		"family_id",
	).Values(
		zone.ID.String(),
		zone.Name.String(),
		zone.Safety,
		squirrel.Expr("ST_GeomFromText(?, 4326)", polygonWKT),
		zone.CreatedAt,
		zone.UpdatedAt,
		zone.EndedAt,
		zone.FamilyID,
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

func (r *ZoneRepositoryImpl) ZoneOverlaps(
	ctx context.Context,
	familyID value_object.ID,
	boundaries *[]value_object.Location) (bool, error) {
	polygonWKT := ConvertBoundariesToPolygon(boundaries)
	if polygonWKT == "" {
		r.observer.Logger.Trace().Msg("invalid boundaries: not enough points or invalid polygon")
		return false, fmt.Errorf("invalid boundaries: %w", domain.ErrInvalidZoneBoundaries)
	}

	// выбираем количество зон для данной семьи, у которых пересечение текущего boundaries с новой геометрией (построенной из polygonWKT) имеет ненулевую площадь
	stmt := r.pg.Builder.
		Select("count(*)").
		From("zones").
		Where(squirrel.Eq{
			"family_id": familyID.String(),
			"ended_at":  nil,
		}).
		Where(squirrel.Expr(
			"ST_Area(ST_Intersection(boundaries, ST_GeomFromText(?, 4326))) > 0",
			polygonWKT,
		))

	sqlQuery, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query for zone overlap check")
		return false, fmt.Errorf("failed to build SQL query: %w", adapter.ErrStorage)
	}

	var count int
	err = r.pg.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&count)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute zone overlap check query")
		return false, fmt.Errorf("failed to execute query: %w", adapter.ErrStorage)
	}

	// Если count > 0, значит существует хотя бы одна зона, пересечение с которой имеет ненулевую площадь
	return count > 0, nil
}

func (r *ZoneRepositoryImpl) CoordinatesInFamilyZones(
	ctx context.Context,
	familyID value_object.ID,
	lon, lat float64,
) ([]*entity.Zone, error) {
	builder := r.pg.Builder.PlaceholderFormat(squirrel.Dollar)
	pointExpr := squirrel.Expr(
		"ST_Contains(boundaries, ST_SetSRID(ST_MakePoint(?, ?), 4326))",
		lon, lat,
	)
	stmt := builder.
		Select(
			"id",
			"name",
			"safety",
			"ST_AsText(boundaries) AS boundaries",
			"created_at",
			"updated_at",
			"family_id",
		).
		From("zones").
		Where(squirrel.Eq{
			"family_id": familyID.String(),
			"ended_at":  nil,
		}).
		Where(pointExpr).
		Limit(1)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query for zones by point")
		return nil, fmt.Errorf("%w: failed to build SQL query: %w", err, adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute zones by point query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	var zoneDBs []*response.ZoneDB
	err = r.pg.Scan.ScanAll(&zoneDBs, rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || len(zoneDBs) == 0 {
			r.observer.Logger.Trace().Msg("no zones contain this point")
			return []*entity.Zone{}, nil
		}
		r.observer.Logger.Error().Err(err).Msg("failed to scan zones result")
		return nil, fmt.Errorf("%w: failed to scan zones result: %w", err, adapter.ErrStorage)
	}

	var zones []*entity.Zone
	for _, zdb := range zoneDBs {
		zone, err := response.ZoneDBToEntity(zdb)
		if err != nil {
			r.observer.Logger.Error().Err(err).Msg("failed to convert zone DB to entity")
			return nil, fmt.Errorf("%w: failed to convert db data to entity: %w", err, adapter.ErrStorage)
		}
		zones = append(zones, zone)
	}
	fmt.Printf("zones: %v\n", zones)

	return zones, nil
}

func (r *ZoneRepositoryImpl) GetZone(ctx context.Context, id value_object.ID) (*entity.Zone, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"safety",
		"ST_AsText(boundaries) AS boundaries",
		"created_at",
		"updated_at",
		"family_id",
	).
		From("zones").
		Where(squirrel.Eq{"id": id.String()}).
		Where(squirrel.Eq{"ended_at": nil})

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

	var resultZone response.ZoneDB
	err = r.pg.Scan.ScanOne(&resultZone, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("zone not found")
		return nil, adapter.ErrZoneNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	zone, err := response.ZoneDBToEntity(&resultZone)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert db data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}

	return zone, nil
}

func (r *ZoneRepositoryImpl) GetZones(ctx context.Context) ([]*entity.Zone, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"name",
		"safety",
		"boundaries",
		"created_at",
		"updated_at",
		"family_id").
		From("zones").
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var zonesDb []*response.ZoneDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&zonesDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	zones, err := response.ZoneDbListToEntityList(zonesDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return zones, nil
}

//func (r *ZoneRepositoryImpl) GetZonesByAuthorID(ctx context.Context, id value_object.externalID) ([]*entity.Zone, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"z_name",
//		"author_id",
//		"safety",
//		"created_at",
//		"updated_at",
//	).Column(squirrel.Expr("ST_AsText(boundaries) AS boundaries")).
//		From("zones").
//		Where(squirrel.Eq{"author_id": id.String()})
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
//		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
//		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
//	}
//
//	zones, err := response.ZoneDbListToEntityList(zonesDb)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
//		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
//	}
//	return zones, nil
//}

//func (r *ZoneRepositoryImpl) GetAvailableZonesForSubscription(ctx context.Context,
//	receiverID, senderID value_object.externalID) ([]*entity.Zone, error) {
//	// Полный SQL-запрос с подзапросами
//	query := `
//		SELECT z.id, z.z_name, z.author_id, z.safety, z.boundaries, z.created_at, z.updated_at
//		FROM zones z
//		WHERE z.author_id IN (
//		    SELECT user_id
//		    FROM families_memberships
//		    WHERE family_id IN (
//		        SELECT family_id
//		        FROM families_memberships
//		        WHERE user_id = $1
//		          AND family_id IN (
//		              SELECT family_id
//		              FROM families_memberships
//		              WHERE user_id = $2
//		          )
//		    )
//		)
//	`
//
//	// Аргументы для SQL-запроса
//	args := []interface{}{receiverID.String(), senderID.String()}
//
//	// Выполняем SQL-запрос
//	var zonesDb []*response.ZoneDB
//	rows, err := r.pg.Pool.Query(ctx, query, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	// Сканируем результат в срез структур ZoneDB
//	for rows.Next() {
//		var zone response.ZoneDB
//		if err := rows.Scan(
//			&zone.externalID, &zone.Name, &zone.AuthorID, &zone.Safety, &zone.Boundaries,
//			&zone.CreatedAt, &zone.UpdatedAt,
//		); err != nil {
//			return nil, fmt.Errorf("failed to scan row: %w", err)
//		}
//		zonesDb = append(zonesDb, &zone)
//	}
//
//	if rows.Err() != nil {
//		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
//	}
//
//	zones, err := response.ZoneDbListToEntityList(zonesDb)
//	if err != nil {
//		return nil, fmt.Errorf("failed to convert storage data to entity list: %w", err)
//	}
//
//	return zones, nil
//}

func (r *ZoneRepositoryImpl) UpdateZone(ctx context.Context, zone *entity.Zone) error {
	stmt := r.pg.Builder.Update("zones").
		SetMap(map[string]interface{}{
			"name":       zone.Name.String(),
			"safety":     zone.Safety,
			"boundaries": squirrel.Expr("ST_GeomFromText(?, 4326)", ConvertBoundariesToPolygon(zone.Boundaries)),
			"updated_at": zone.UpdatedAt,
			"ended_at":   zone.EndedAt,
		}).
		Where(squirrel.Eq{"id": zone.ID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

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

//func (r *ZoneRepositoryImpl) DeleteZone(ctx context.Context, id value_object.externalID) error {
//	deleteFamiliesZonesStmt := r.pg.Builder.Delete("families_zones").
//		Where(squirrel.Eq{"zone_id": id.String()})
//
//	sqlFamiliesZones, argsFamiliesZones, err := deleteFamiliesZonesStmt.ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build SQL string for deleting families_zones: %w", err)
//	}
//
//	_, err = r.pg.Pool.Exec(ctx, sqlFamiliesZones, argsFamiliesZones...)
//	if err != nil {
//		return fmt.Errorf("failed to execute query for deleting families_zones: %w", err)
//	}
//
//	deleteZoneStmt := r.pg.Builder.Delete("zones").
//		Where(squirrel.Eq{"id": id.String()})
//
//	sqlZone, argsZone, err := deleteZoneStmt.ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build SQL string for deleting zone: %w", err)
//	}
//
//	_, err = r.pg.Pool.Exec(ctx, sqlZone, argsZone...)
//	if err != nil {
//		return fmt.Errorf("failed to execute query for deleting zone: %w", err)
//	}
//
//	return nil
//}

func ConvertBoundariesToPolygon(boundaries *[]value_object.Location) string {
	fmt.Println(boundaries, "aboba")
	if len(*boundaries) < 4 {
		return ""
	}

	var points []string
	for _, boundary := range *boundaries {
		points = append(points, fmt.Sprintf("%f %f", boundary.Longitude, boundary.Latitude))
	}

	return fmt.Sprintf("POLYGON((%s))", strings.Join(points, ", "))
}

//func ConvertBoundariesToPolygon(boundaries *[]value_object.Location) string {
//	if len(*boundaries) < 3 {
//		return ""
//	}
//
//	pts := *boundaries
//	// Проверяем, совпадают ли первая и последняя точки
//	first := pts[0]
//	last := pts[len(pts)-1]
//	if first.Latitude != last.Latitude ||
//		first.Longitude != last.Longitude {
//		pts = append(pts, first)
//	}
//
//	// Формируем WKT‑строку
//	wktPoints := make([]string, len(pts))
//	for i, p := range pts {
//		wktPoints[i] = fmt.Sprintf("%f %f", p.Longitude, p.Latitude)
//	}
//
//	return fmt.Sprintf("POLYGON((%s))", strings.Join(wktPoints, ", "))
//}
