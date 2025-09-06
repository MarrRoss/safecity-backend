package db_storage

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
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

type InvitationToFamilyRepositoryImpl struct {
	pg       *postgres.Postgres
	observer *observability.Observability
}

func NewInvitationToFamilyRepositoryImpl(
	pg *postgres.Postgres,
	observer *observability.Observability) (*InvitationToFamilyRepositoryImpl, error) {
	return &InvitationToFamilyRepositoryImpl{
		pg:       pg,
		observer: observer,
	}, nil
}

func (r *InvitationToFamilyRepositoryImpl) AddInvitationToFamily(
	ctx context.Context,
	invitation *aggregate.InvitationToFamily,
) error {
	stmt := r.pg.Builder.Insert("invitations_to_families").Columns(
		"id",
		"author_id",
		"role_id",
		"family_id",
		"created_at",
		"ended_at",
	).Values(
		invitation.ID.String(),
		invitation.Author.ID.String(),
		invitation.Role.ID,
		invitation.Family.ID.String(),
		invitation.CreatedAt,
		invitation.EndedAt,
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

//func (r *InvitationToFamilyRepositoryImpl) AddInvitationToFamily(ctx context.Context, invitation *aggregate.InvitationToFamily) error {
//	// Проверка существования активного приглашения
//	checkQuery := `
//		SELECT id
//		FROM invitations_to_families
//		WHERE author_id = $1
//		  AND receiver_id = $2
//		  AND role_id = $3
//		  AND family_id = $4
//		  AND active = TRUE
//		  AND accepted = FALSE
//	`
//	var existingInvitationID string
//	err := r.pg.Pool.QueryRow(ctx, checkQuery,
//		invitation.Author.ID.String(),
//		invitation.Receiver.ID.String(),
//		invitation.Role.ID,
//		invitation.Family.ID.String(),
//	).Scan(&existingInvitationID)
//	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
//		return fmt.Errorf("failed to check existing invitation: %w", err)
//	}
//
//	// Если найдено активное приглашение, обновляем его
//	if existingInvitationID != "" {
//		updateQuery := `
//			UPDATE invitations_to_families
//			SET active = FALSE, accepted = FALSE, updated_at = NOW()
//			WHERE id = $1
//		`
//		_, err = r.pg.Pool.Exec(ctx, updateQuery, existingInvitationID)
//		if err != nil {
//			return fmt.Errorf("failed to update existing invitation: %w", err)
//		}
//	}
//
//	// Добавление нового приглашения
//	stmt := r.pg.Builder.Insert("invitations_to_families").Columns(
//		"id",
//		"author_id",
//		"receiver_id",
//		"role_id",
//		"family_id",
//		"created_at",
//		"active",
//		"accepted",
//		"updated_at",
//	).Values(
//		invitation.ID.String(),
//		invitation.Author.ID.String(),
//		invitation.Receiver.ID.String(),
//		invitation.Role.ID,
//		invitation.Family.ID.String(),
//		invitation.CreatedAt,
//		invitation.Active,
//		invitation.Accepted,
//		invitation.UpdatedAt,
//	)
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//	_, err = r.pg.Pool.Exec(ctx, sql, args...)
//	if err != nil {
//		return fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	return nil
//}

func (r *InvitationToFamilyRepositoryImpl) GetInvitationToFamily(
	ctx context.Context,
	id value_object.ID,
) (*response.ShortInvitationToFamilyDB, error) {
	stmt := r.pg.Builder.Select(
		"id",
		"author_id",
		"role_id",
		"family_id",
		"created_at").
		From("invitations_to_families").
		Where(squirrel.Eq{"id": id.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultInvitation response.ShortInvitationToFamilyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&resultInvitation, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("invitation not found")
		return nil, adapter.ErrInvitationNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	return &resultInvitation, nil
}

func (r *InvitationToFamilyRepositoryImpl) GetInvitationToFamilyByID(
	ctx context.Context,
	invitationID value_object.ID,
) (*aggregate.InvitationToFamily, error) {
	stmt := r.pg.Builder.
		Select(
			"inv.id",
			"inv.created_at   AS invitation_created_at",
			"inv.ended_at     AS invitation_ended_at",
			"auth.id          AS author_id",
			"auth.first_name  AS author_first_name",
			"auth.last_name   AS author_last_name",
			"auth.email",
			"auth.username",
			"auth.tracking",
			"ui.external_id   AS telegram_id",
			"auth.created_at  AS author_created_at",
			"auth.updated_at  AS author_updated_at",
			"inv.family_id",
			"f.name           AS family_name",
			"f.created_at     AS family_created_at",
			"f.updated_at     AS family_updated_at",
			"f.ended_at       AS family_ended_at",
			"inv.role_id",
			"r.name           AS role_name",
		).
		From("invitations_to_families inv").
		Join("users auth ON auth.id = inv.author_id").
		LeftJoin("users_integrations ui ON ui.user_id = auth.id").
		Join("families f ON f.id = inv.family_id").
		Join("roles r ON r.id = inv.role_id").
		Where("inv.id = ?", invitationID.String())

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultInvitation response.InvitationToFamilyDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&resultInvitation, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("invitation not found")
		return nil, adapter.ErrInvitationNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}
	invitation, err := response.InvitationDbToEntity(&resultInvitation)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return invitation, nil
}

//func (r *InvitationToFamilyRepositoryImpl) GetInvitationsToFamilyByReceiverID(ctx context.Context,
//	receiverID value_object.ID) ([]*response.InvitationToFamilyDB, error) {
//	stmt := r.pg.Builder.Select(
//		"id",
//		"author_id",
//		"receiver_id",
//		"role_id",
//		"family_id",
//		"created_at",
//		"active",
//		"accepted",
//		"updated_at").
//		From("invitations_to_families").
//		Where(squirrel.Eq{"receiver_id": receiverID.String()})
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to build an SQL string from the query: %w", err)
//	}
//
//	var familyInvitationsDb []*response.InvitationToFamilyDB
//	err = pgxscan.Select(ctx, r.pg.Pool, &familyInvitationsDb, sql, args...)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//	return familyInvitationsDb, nil
//}

func (r *InvitationToFamilyRepositoryImpl) GetPendingInvitationsByUser(
	ctx context.Context, userID value_object.ID,
) ([]*aggregate.InvitationToFamily, error) {
	stmt := r.pg.Builder.
		Select(
			"inv.id",
			"inv.created_at AS invitation_created_at",
			"inv.ended_at   AS invitation_ended_at",
			"auth.id AS author_id",
			"auth.first_name AS author_first_name",
			"auth.last_name AS author_last_name",
			"auth.email",
			"auth.username",
			"auth.tracking",
			"ui.external_id AS telegram_id",
			"auth.created_at AS author_created_at",
			"auth.updated_at AS author_updated_at",
			"inv.family_id",
			"f.name AS family_name",
			"f.created_at   AS family_created_at",
			"f.updated_at   AS family_updated_at",
			"f.ended_at     AS family_ended_at",
			"inv.role_id",
			"r.name AS role_name",
		).
		From("invitations_to_families inv").
		Join("users auth ON auth.id = inv.author_id").
		LeftJoin("users_integrations ui ON ui.user_id = auth.id").
		Join("families f ON f.id = inv.family_id").
		Join("roles r ON r.id = inv.role_id").
		LeftJoin("invites_activations act ON act.invitation_id = inv.id AND act.user_id = ?", userID.String()).
		Where(squirrel.Eq{
			"inv.ended_at": nil,
		}).
		Where("act.user_id IS NULL").
		Where(squirrel.NotEq{"inv.author_id": userID.String()})

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

	var resultInvitation []*response.InvitationToFamilyDB
	if err := r.pg.Scan.ScanAll(&resultInvitation, rows); err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	invitation, err := response.InvitationDbListToEntityList(resultInvitation)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return invitation, nil
}

func (r *InvitationToFamilyRepositoryImpl) AddInviteActivation(
	ctx context.Context,
	invitationID, userID value_object.ID) error {
	now := time.Now()
	stmt := r.pg.Builder.Insert("invites_activations").Columns(
		"invitation_id",
		"user_id",
		"created_at",
	).Values(
		invitationID.String(),
		userID.String(),
		now,
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

//func (r *InvitationToFamilyRepositoryImpl) AcceptInvitationToFamily(ctx context.Context,
//	invitationID value_object.ID) error {
//	stmt := r.pg.Builder.Update("invitations_to_families").
//		Set("active", false).
//		Set("accepted", true).
//		Where("id = ?", invitationID.String())
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
