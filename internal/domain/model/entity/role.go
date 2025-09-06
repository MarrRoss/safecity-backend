package entity

import (
	"awesomeProjectDDD/internal/domain/model/value_object"
	"fmt"
	"time"
)

type Role struct {
	ID        int
	Name      value_object.RoleName
	CreatedAt time.Time
}

//func NewRole(name value_object.RoleName) (*Role, error) {
//	var newRole Role
//	if name.String() == "admin" {
//		newRole = Role{
//			ID:        1,
//			Name:      name,
//			CreatedAt: time.Now(),
//		}
//	}
//	if name.String() == "moderator" {
//		newRole = Role{
//			ID:        2,
//			Name:      name,
//			CreatedAt: time.Now(),
//		}
//	}
//	if name.String() == "adult" {
//		newRole = Role{
//			ID:        3,
//			Name:      name,
//			CreatedAt: time.Now(),
//		}
//	}
//	if name.String() == "child" {
//		newRole = Role{
//			ID:        4,
//			Name:      name,
//			CreatedAt: time.Now(),
//		}
//	}
//	return &newRole, nil
//}

func (r *Role) IDToString() string {
	return fmt.Sprintf("%d", r.ID)
}
