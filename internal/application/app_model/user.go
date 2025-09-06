package app_model

import (
	"awesomeProjectDDD/internal/domain/model/entity"
	"github.com/google/uuid"
	"time"
)

type ApplicationUser struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Login        string
	Tracking     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TgIntegrated bool
}

func NewApplicationUser(user *entity.User) *ApplicationUser {
	return &ApplicationUser{
		ID:           user.ID.ToRaw(),
		FirstName:    user.Name.FirstName.String(),
		LastName:     user.Name.LastName.String(),
		Email:        user.Email.String(),
		Login:        user.Login.String(),
		Tracking:     user.Tracking,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		TgIntegrated: user.TgIntegrated,
	}
}

func NewApplicationUsers(users []*entity.User) []*ApplicationUser {
	appUsers := make([]*ApplicationUser, len(users))
	for i, user := range users {
		appUsers[i] = NewApplicationUser(user)
	}
	return appUsers
}
