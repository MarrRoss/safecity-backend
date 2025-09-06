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

type UserRepositoryImpl struct {
	observer *observability.Observability
	pg       *postgres.Postgres
}

func NewUserRepositoryImpl(
	observer *observability.Observability,
	pg *postgres.Postgres,
) (*UserRepositoryImpl, error) {
	return &UserRepositoryImpl{
		observer: observer,
		pg:       pg,
	}, nil
}

func (r *UserRepositoryImpl) AddUser(ctx context.Context, user *entity.User) error {
	stmt := r.pg.Builder.Insert("users").Columns(
		"id",
		"first_name",
		"last_name",
		"email",
		"username",
		"external_id",
		"tracking",
		"created_at",
		"updated_at",
		"ended_at",
	).Values(
		user.ID.String(),
		user.Name.FirstName.String(),
		user.Name.LastName.String(),
		user.Email.String(),
		user.Login.String(),
		user.ExternalID.String(),
		user.Tracking,
		user.CreatedAt,
		user.UpdatedAt,
		user.EndedAt,
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

func (r *UserRepositoryImpl) GetUser(ctx context.Context, id value_object.ID) (*entity.User, error) {
	stmt := r.pg.Builder.Select(
		"users.id as id",
		"first_name",
		"last_name",
		"email",
		"username",
		"tracking",
		"users.created_at",
		"users.updated_at",
		"users_integrations.external_id as telegram_id",
	).
		From("users").
		LeftJoin("users_integrations ON users_integrations.user_id = users.id").
		Where(squirrel.Eq{"users.id": id.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultUser response.UserDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&resultUser, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("user not found")
		return nil, adapter.ErrUserNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	user, err := response.UserDbToEntity(&resultUser)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetUserByExternalID(ctx context.Context, id value_object.ID) (*entity.User, error) {
	stmt := r.pg.Builder.Select(
		"users.id as id",
		"first_name",
		"last_name",
		"email",
		"username",
		"tracking",
		"users.created_at",
		"users.updated_at",
		"users_integrations.external_id as telegram_id",
	).
		From("users").
		LeftJoin("users_integrations ON users_integrations.user_id = users.id").
		Where(squirrel.Eq{"users.external_id": id.String()}).
		Where(squirrel.Eq{"ended_at": nil})

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return nil, fmt.Errorf("%w: failed to build an SQL string from the query: %w", err, adapter.ErrStorage)
	}

	var resultUser response.UserDB
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return nil, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	defer rows.Close()

	err = r.pg.Scan.ScanOne(&resultUser, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("user not found")
		return nil, adapter.ErrUserNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	user, err := response.UserDbToEntity(&resultUser)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetUsers(ctx context.Context, login *value_object.Login) ([]*entity.User, error) {
	stmt := r.pg.Builder.Select(
		"users.id as id",
		"first_name",
		"last_name",
		"email",
		"username",
		"tracking",
		"users.created_at",
		"users.updated_at",
		"users_integrations.external_id as telegram_id",
	).
		From("users").
		LeftJoin("users_integrations ON users_integrations.user_id = users.id").
		Where(squirrel.Eq{"ended_at": nil})

	if login != nil {
		l := fmt.Sprintf("%%%s%%", login.String())
		stmt = stmt.Where(squirrel.ILike{"username": l})
	}

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

func (r *UserRepositoryImpl) GetUsersByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	stmt := r.pg.Builder.Select(
		"users.id as id",
		"first_name",
		"last_name",
		"email",
		"username",
		"tracking",
		"users.created_at",
		"users.updated_at",
		"users_integrations.external_id as telegram_id",
	).
		From("users").
		LeftJoin("users_integrations ON users_integrations.user_id = users.id").
		Where(squirrel.Eq{"users.id": ids}).
		Where(squirrel.Eq{"ended_at": nil})

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
	if errors.Is(err, pgx.ErrNoRows) {
		r.observer.Logger.Trace().Err(err).Msg("membership users not found")
		return nil, adapter.ErrMembershipUsersNotFound
	} else if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to scan result")
		return nil, fmt.Errorf("%w: failed to scan result to variable: %w", err, adapter.ErrStorage)
	}

	users, err := response.UserDbListToEntityList(usersDb)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to convert storage data to entity list")
		return nil, fmt.Errorf("%w: failed to convert storage data to entity: %w", err, adapter.ErrStorage)
	}
	return users, nil
}

func (r *UserRepositoryImpl) UpdateUser(
	ctx context.Context,
	user *entity.User,
) error {
	stmt := r.pg.Builder.Update("users").
		SetMap(map[string]interface{}{
			"first_name": user.Name.FirstName.String(),
			"last_name":  user.Name.LastName.String(),
			"updated_at": user.UpdatedAt,
			"ended_at":   user.EndedAt,
			"tracking":   user.Tracking,
		}).
		Where(squirrel.Eq{"id": user.ID.String()}).
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

//func (r *UserRepositoryImpl) PhoneExists(ctx context.Context, phone value_object.Phone) (bool, error) {
//	stmt := r.pg.Builder.Select("1").
//		Prefix("SELECT EXISTS (").
//		From("users").
//		Where(squirrel.Expr("phone = ?", phone.String())).
//		Suffix(")")
//
//	sql, args, err := stmt.ToSql()
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
//		return false, fmt.Errorf("%w: failed to build SQL query: %w", err, adapter.ErrStorage)
//	}
//
//	var exists bool
//	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
//	if err != nil {
//		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
//		return false, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
//	}
//
//	return exists, nil
//}

func (r *UserRepositoryImpl) EmailExists(ctx context.Context, email value_object.Email) (bool, error) {
	stmt := r.pg.Builder.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where(squirrel.Expr("email = ?", email.String())).
		Suffix(")")

	sql, args, err := stmt.ToSql()
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to build SQL query")
		return false, fmt.Errorf("%w: failed to build SQL query: %w", err, adapter.ErrStorage)
	}

	var exists bool
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		r.observer.Logger.Error().Err(err).Msg("failed to execute query")
		return false, fmt.Errorf("%w: failed to execute query: %w", err, adapter.ErrStorage)
	}
	return exists, nil
}
