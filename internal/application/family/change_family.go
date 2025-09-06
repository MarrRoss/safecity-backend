package family

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ChangeFamilyCommand struct {
	ID   uuid.UUID
	Name string
}

type ChangeFamilyHandler struct {
	storage  db.FamilyRepository
	observer *observability.Observability
}

func NewChangeFamilyHandler(
	storage db.FamilyRepository,
	observer *observability.Observability) *ChangeFamilyHandler {
	return &ChangeFamilyHandler{
		storage:  storage,
		observer: observer,
	}
}

func (h *ChangeFamilyHandler) Handle(ctx context.Context, cmd ChangeFamilyCommand) error {
	id, err := value_object.NewIDFromString(cmd.ID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid family id")
		return err
	}
	family, err := h.storage.GetFamily(ctx, id)
	if err != nil {
		if errors.Is(err, adapter.ErrFamilyNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("family not found")
			return err
		}

		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family")
			return err
		}

		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family")
		return err
	}

	name, err := value_object.NewFamilyName(cmd.Name)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create family name")
		return fmt.Errorf("failed to create family name: %w", err)
	}
	err = family.ChangeName(name)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to change family name")
		return fmt.Errorf("failed to change family name: %w", err)
	}
	err = h.storage.UpdateFamily(ctx, family)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while updating family")
			return err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while updating family")
		return err
	}
	return nil
}
