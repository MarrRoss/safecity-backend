package db

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, id value_object.ID) (*entity.User, error)
	GetUserByExternalID(ctx context.Context, id value_object.ID) (*entity.User, error)
	GetUsers(ctx context.Context, login *value_object.Login) ([]*entity.User, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	//PhoneExists(ctx context.Context, phone value_object.Phone) (bool, error)
	EmailExists(ctx context.Context, email value_object.Email) (bool, error)
}
