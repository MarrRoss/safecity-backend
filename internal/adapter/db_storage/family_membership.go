package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/aggregate"
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

type FamilyMembershipRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewFamilyMembershipRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability) (*FamilyMembershipRepositoryImpl, error) {
	return &FamilyMembershipRepositoryImpl{pg: pg, observer: observer}, nil
}

func (r *FamilyMembershipRepositoryImpl) AddFamilyMembership(
	ctx context.Context, familyMembership *aggregate.FamilyMembership) error {
	stmt := r.pg.Builder.Insert("families_memberships").Columns(
		"id",
		"user_id",
		"role_id",
		"family_id",
		"created_at",
		"ended_at",
		"updated_at",
	).Values(
		familyMembership.ID.String(),
		familyMembership.User.ID.String(),
		familyMembership.Role.ID,
		familyMembership.Family.ID.String(),
		familyMembership.CreatedAt,
		familyMembership.EndedAt,
		familyMembership.UpdatedAt,
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

func (r *FamilyMembershipRepositoryImpl) GetFamilyMembershipByID(
	ctx context.Context,
	id value_object.ID) (*response.FamilyMembershipDB, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"user_id",
		"role_id",
		"family_id",
		"created_at",
		"updated_at").
		From("families_memberships").
		Where(squirrel.Eq{"id": id.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}
	var familyMembership response.FamilyMembershipDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()
	err = r.pg.Scan.ScanOne(&familyMembership, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("membership not found")
		return nil, adapter.ErrMembershipNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return &familyMembership, nil
}

func (r *FamilyMembershipRepositoryImpl) GetFamiliesByUserID(ctx context.Context,
	userID value_object.ID) ([]*entity.Family, error) {
	stmt := r.pg.Builder.Select(
		"f.id",
		"f.name",
		"f.created_at",
		"f.updated_at").
		From("families_memberships fm").
		Join("families f ON fm.family_id = f.id").
		Where(squirrel.Eq{
			"fm.user_id":  userID.String(),
			"fm.ended_at": nil,
			"f.ended_at":  nil,
		})

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
		return nil, nil
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

func (r *FamilyMembershipRepositoryImpl) MembershipExists(
	ctx context.Context, userID, familyID value_object.ID) (bool, error) {
	stmt := r.pg.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("families_memberships").
		Where(squirrel.Eq{
			"user_id":   userID.String(),
			"family_id": familyID.String(),
		}).
		Suffix(")").
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build an SQL string from the query")
		return false, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return false, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	var exists bool
	err = r.pg.Scan.ScanOne(&exists, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return false, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	return exists, nil
}

func (r *FamilyMembershipRepositoryImpl) CheckUsersBelongToFamilyByZone(
	ctx context.Context,
	zoneID, senderID, receiverID value_object.ID,
) (bool, error) {
	stmt := r.pg.Builder.
		Select("COUNT(DISTINCT fm.user_id) AS cnt").
		From("families_memberships fm").
		Join("families f ON fm.family_id = f.id").
		Where(squirrel.Eq{
			"fm.user_id":  []string{senderID.String(), receiverID.String()},
			"fm.ended_at": nil,
			"f.ended_at":  nil,
		}).
		Where(squirrel.Expr(
			"fm.family_id = (SELECT family_id FROM zones WHERE id = ? AND ended_at IS NULL)",
			zoneID.String(),
		))

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return false, fmt.Errorf("%w: failed to build SQL query", adapter.ErrStorage)
	}

	var count int
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.observer.Logger.Trace().Err(err).Msg("membership not found")
			return false, adapter.ErrMembershipNotFound
		}
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return false, fmt.Errorf("%w: failed to execute query", adapter.ErrStorage)
	}

	if count != 2 {
		r.observer.Logger.Trace().Msg("receiver and sender are not family members")
		return false, adapter.ErrNotFamilyMembers
	}

	return true, nil
}

func (r *FamilyMembershipRepositoryImpl) CheckUsersShareCommonActiveFamily(
	ctx context.Context,
	senderID, receiverID value_object.ID,
) (bool, error) {
	// выбираем family_id из таблицы families_memberships, где для каждого family_id
	// найдены активные записи (ended_at IS NULL) для обоих пользователей, а сама семья активна.
	subquery := r.pg.Builder.
		Select("fm.family_id").
		From("families_memberships fm").
		Join("families f ON fm.family_id = f.id").
		Where(squirrel.Eq{
			"fm.user_id":  []string{senderID.String(), receiverID.String()},
			"fm.ended_at": nil,
			"f.ended_at":  nil,
		}).
		GroupBy("fm.family_id").
		Having("COUNT(DISTINCT fm.user_id) = 2").
		PlaceholderFormat(squirrel.Dollar)

	// Внешний запрос: считаем количество семей из подзапроса
	stmt := r.pg.Builder.
		Select("COUNT(*) AS cnt").
		FromSelect(subquery, "common_families").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return false, fmt.Errorf("%w: failed to build SQL query", adapter.ErrStorage)
	}

	var count int
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.observer.Logger.Trace().Err(err).Msg("no common family membership found")
			return false, adapter.ErrMembershipNotFound
		}
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return false, fmt.Errorf("%w: failed to execute query", adapter.ErrStorage)
	}

	if count == 0 {
		r.observer.Logger.Trace().Msg("receiver and sender are not family members")
		return false, adapter.ErrNotFamilyMembers
	}

	return true, nil
}

//func (r *FamilyMembershipRepositoryImpl) GetMembershipsByUserID(ctx context.Context,
//	userID value_object.externalID) ([]*response.FamilyMembershipDB, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"user_id",
//		"role_id",
//		"family_id",
//		"created_at").
//		From("families_memberships").
//		Where(squirrel.Eq{"user_id": userID.String()})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var familyMembershipsDb []*response.FamilyMembershipDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	err = r.pg.Scan.ScanAll(&familyMembershipsDb, rows)
//	if err != nil {
//		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
//	}
//	return familyMembershipsDb, nil
//}

func (r *FamilyMembershipRepositoryImpl) GetMembershipsByFamilyID(
	ctx context.Context,
	familyID value_object.ID,
	usersIDs *[]string,
) ([]*response.FamilyMembershipDB, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"user_id",
		"role_id",
		"family_id",
		"created_at",
		"updated_at").
		From("families_memberships").
		Where(squirrel.Eq{"family_id": familyID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	if usersIDs != nil {
		stmt = stmt.Where(squirrel.Eq{"user_id": usersIDs})
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
	}

	var familyMembershipsDb []*response.FamilyMembershipDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&familyMembershipsDb, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan result to variable: %w", err)
	}
	return familyMembershipsDb, nil
}

func (r *FamilyMembershipRepositoryImpl) GetMembershipsParticipantsByUser(ctx context.Context,
	userID value_object.ID) ([]*response.FamilyMembershipParticipants, error) {
	stmt := r.pg.Builder.Select(
		"fm.id AS membership_id",
		"u.id AS user_id",
		"r.name AS role_name",
		"f.id AS family_id").
		From("families_memberships fm").
		Join("users u ON fm.user_id = u.id").
		Join("roles r ON fm.role_id = r.id").
		Join("families f ON fm.family_id = f.id").Where(squirrel.Eq{"fm.ended_at": nil}).
		Where(squirrel.Expr(
			"fm.family_id = ANY (SELECT family_id FROM families_memberships WHERE user_id = ?::uuid)",
			userID.String(),
		))

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build SQL string: %w", err, adapter.ErrStorage)
	}

	var memberships []*response.FamilyMembershipParticipants
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&memberships, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result to variable")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	if len(memberships) == 0 {
		r.observer.Logger.Trace().Str("user_id", userID.String()).Msg("no memberships found for user")
	}

	return memberships, nil
}

func (r *FamilyMembershipRepositoryImpl) GetAvailableNotificationSenders(
	ctx context.Context,
	receiverID, familyID value_object.ID) ([]*entity.User, error) {
	stmt := r.pg.Builder.
		Select(
			"u.id",
			"u.first_name",
			"u.last_name",
			"u.email",
			"u.username",
			"u.tracking",
			"ui.external_id AS telegram_id",
			"u.created_at",
			"u.updated_at",
		).
		From("families_memberships fm").
		Join("users u ON u.id = fm.user_id").
		LeftJoin("users_integrations ui ON ui.user_id = u.id").
		Where(squirrel.Eq{
			"fm.family_id": familyID.String(),
			"fm.ended_at":  nil,
			"u.ended_at":   nil,
			"u.tracking":   true,
			"fm.role_id":   2,
		}).
		Where(squirrel.NotEq{"u.id": receiverID.String()})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var usersDb []*response.UserDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanAll(&usersDb, rows)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	users, err := response.UserDbListToEntityList(usersDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return users, nil
}

// TODO: переделать через домен
func (r *FamilyMembershipRepositoryImpl) DeleteFamilyMembership(
	ctx context.Context,
	membershipID value_object.ID,
) error {
	now := time.Now()
	stmt := r.pg.Builder.Update("families_memberships").
		SetMap(map[string]interface{}{
			"updated_at": now,
			"ended_at":   now,
		}).
		Where(squirrel.Eq{"id": membershipID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	if rows.RowsAffected() == 0 {
		r.observer.Logger.Trace().Err(err).Msg("membership not found")
		return adapter.ErrMembershipNotFound
	}
	return nil
}

// TODO: переделать через домен
func (r *FamilyMembershipRepositoryImpl) UpdateFamilyMembershipRole(
	ctx context.Context,
	membershipID value_object.ID,
	roleID int,
) error {
	now := time.Now()
	stmt := r.pg.Builder.Update("families_memberships").
		SetMap(map[string]interface{}{
			"role_id":    roleID,
			"updated_at": now,
		}).
		Where(squirrel.Eq{"id": membershipID.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	rows, err := r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	if rows.RowsAffected() == 0 {
		r.observer.Logger.Trace().Err(err).Msg("membership not found")
		return adapter.ErrMembershipNotFound
	}
	return nil
}

//func (r *FamilyMembershipRepositoryImpl) GetMembershipByUserAndZone(ctx context.Context,
//	userID, zoneID value_object.externalID) (*response.FamilyMembershipDB, error) {
//	// Подзапрос для получения family_id, связанного с zone_id
//	zoneFamilyQuery := r.pg.Builder.
//		Select("family_id").
//		From("zones").
//		Where(squirrel.Eq{"id": zoneID.String()})
//
//	// Основной запрос для получения membership
//	stmt := r.pg.Builder.
//		Select("fm.id", "fm.user_id", "fm.role_id", "fm.family_id", "fm.created_at").
//		From("families_memberships fm").
//		Where(squirrel.And{
//			squirrel.Eq{"fm.user_id": userID.String()},
//			squirrel.Expr("fm.family_id IN ?", zoneFamilyQuery),
//		})
//
//	// Преобразование в SQL
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to build SQL string: %w", err)
//	}
//
//	// Выполнение запроса
//	var membershipDB response.FamilyMembershipDB
//	rows, err := r.pg.Pool.Query(ctx, sql, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	defer rows.Close()
//
//	// Сканируем результат
//	err = r.pg.Scan.ScanOne(&membershipDB, rows)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, fmt.Errorf("no membership found for user %s in zone %s", userID.String(), zoneID.String())
//		}
//		return nil, fmt.Errorf("failed to scan membership: %w", err)
//	}
//
//	return &membershipDB, nil
//}

//func (r *FamilyMembershipRepositoryImpl) DeleteFamilyMembership(ctx context.Context, id value_object.externalID) error {
//	stmt := r.pg.Builder.Delete("families_memberships").
//		Where(squirrel.Eq{"id": id.String()})
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

//func (r *FamilyMembershipRepositoryImpl) DeleteFamilyMembershipWithCheck(ctx context.Context, membershipID value_object.externalID) error {
//	// Запрос для получения family_id, user_id и author_id
//	query := `
//		SELECT fm.family_id, fm.user_id, f.author_id
//		FROM families_memberships fm
//		JOIN families f ON fm.family_id = f.id
//		WHERE fm.id = $1
//	`
//
//	var familyID, userID, authorID string
//	err := r.pg.Pool.QueryRow(ctx, query, membershipID.String()).Scan(&familyID, &userID, &authorID)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return fmt.Errorf("membership not found")
//		}
//		return fmt.Errorf("failed to fetch family and user details: %w", err)
//	}
//
//	// Если пользователь — создатель семьи, удаляем семью и связанные данные
//	if userID == authorID {
//		tx, err := r.pg.Pool.Begin(ctx)
//		if err != nil {
//			return fmt.Errorf("failed to start transaction: %w", err)
//		}
//		defer func() {
//			if p := recover(); p != nil {
//				err := tx.Rollback(ctx)
//				if err != nil {
//					log.Errorf("failed to rollback transaction: %v", err)
//				}
//			} else if err != nil {
//				err := tx.Rollback(ctx)
//				if err != nil {
//					log.Errorf("failed to rollback transaction: %v", err)
//				}
//			} else {
//				err = tx.Commit(ctx)
//			}
//		}()
//
//		// Удаляем все приглашения для семьи
//		deleteInvitationsQuery := `
//			DELETE FROM invitations_to_families WHERE family_id = $1
//		`
//		_, err = tx.Exec(ctx, deleteInvitationsQuery, familyID)
//		if err != nil {
//			return fmt.Errorf("failed to delete invitations: %w", err)
//		}
//
//		// Удаляем все memberships для семьи
//		deleteMembershipsQuery := `
//			DELETE FROM families_memberships WHERE family_id = $1
//		`
//		_, err = tx.Exec(ctx, deleteMembershipsQuery, familyID)
//		if err != nil {
//			return fmt.Errorf("failed to delete memberships: %w", err)
//		}
//
//		// Удаляем семью
//		deleteFamilyQuery := `
//			DELETE FROM families WHERE id = $1
//		`
//		_, err = tx.Exec(ctx, deleteFamilyQuery, familyID)
//		if err != nil {
//			return fmt.Errorf("failed to delete family: %w", err)
//		}
//
//		return nil
//	}
//
//	// Если пользователь не создатель семьи, удаляем только указанный membership
//	deleteMembershipQuery := `
//		DELETE FROM families_memberships WHERE id = $1
//	`
//	_, err = r.pg.Pool.Exec(ctx, deleteMembershipQuery, membershipID.String())
//	if err != nil {
//		return fmt.Errorf("failed to delete membership: %w", err)
//	}
//
//	return nil
//}
