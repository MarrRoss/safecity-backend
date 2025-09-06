package response

import (
	"time"
)

type FamilyMembershipDB struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	RoleID    string    `db:"role_id"`
	FamilyID  string    `db:"family_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FamilyMembershipParticipants struct {
	MembershipID string `db:"membership_id"`
	UserID       string `db:"user_id"`
	RoleName     string `db:"role_name"`
	FamilyID     string `db:"family_id"`
}
