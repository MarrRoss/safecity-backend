package DiplomGoNew

import "github.com/google/uuid"

type SystemIntegration struct {
	ID         int    `db:"id"`
	ExternalID string `db:"external_id"`
}

type UserIntegration struct {
	UUID   uuid.UUID         `db:"user_id"`
	System SystemIntegration `db:"system"`
}

//		Select("ns.receiver_id").
//		From("notifications_settings ns").
//		Join("users_integrations ui ON ui.user_id = ns.receiver_id AND ui.system_id = ?", 1).
//		Where(squirrel.Eq{
//			"ns.id":       stringIds,
//			"ns.ended_at": nil,
//		}).

// SELECT
// user_integration.user_id as user_id,
// user_integration.system_id as 'system.id',
// user_integration.external_id as 'system.external_id'
// FROM user_integration
// INNER JOIN notifications_settings ns ON ns.receiver_id = user_integration.user_id
// WHERE user_integration.external_id IS NOT NULL AND ns.id IN (?, ?, ?) AND ns.ended_at IS NULL
