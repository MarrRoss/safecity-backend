package request

import "github.com/google/uuid"

type AddTgConnectRequest struct {
	UserID uuid.UUID `json:"user_id" format:"uuid"`
	TgID   string    `json:"tg_id" format:"string"`
} // @name AddTgConnectRequest
