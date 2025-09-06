package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"context"
)

type IntegrationRepository interface {
	AddIntegration(ctx context.Context, userID value_object.ID, systemID int) error
	GetIntegration(ctx context.Context, userID value_object.ID, systemID int) (*response.IntegrationDB, error)
	UpdateIntegration(ctx context.Context, userID value_object.ID, tgID string) error
}
