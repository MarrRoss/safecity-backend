package request

import (
	"github.com/google/uuid"
)

type GetUsersRequest struct {
	Login *string `query:"login"`
} // @name GetUsersRequest

type GetByUserIDRequest struct {
	ID uuid.UUID `params:"id" format:"uuid"`
} // @name GetByUserIDRequest

//type GetUserByPhoneRequest struct {
//	Phone string `json:"phone" path:"phone"`
//}

type RestoreUserRequest struct {
	ID string `json:"id" path:"id"`
}

type AddUserRequest struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Login      string    `json:"login"`
	ExternalID uuid.UUID `json:"external_id"`
} // @name AddUserRequest

type ChangeUserRequest struct {
	ID        uuid.UUID `params:"id" format:"uuid"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
} // @name ChangeUserRequest

type ChangeUserPasswordRequest struct {
	ID       uuid.UUID `json:"id" format:"uuid"`
	Password string    `json:"password"`
} // @name ChangeUserPasswordRequest

type ChangeUserEmailRequest struct {
	ID    uuid.UUID `json:"id" format:"uuid"`
	Email string    `json:"email"`
} // @name ChangeUserEmailRequest

type ChangeUserPhoneRequest struct {
	ID    uuid.UUID `json:"id" format:"uuid"`
	Phone string    `json:"phone"`
} // @name ChangeUserPhoneRequest

//type ChangeUserTrackingRequest struct {
//	ID uuid.UUID `params:"id" format:"uuid"`
//}

//type DeleteUserRequest struct {
//	ID uuid.UUID `params:"id" format:"uuid"`
//}

type DeleteUserEmailRequest struct {
	ID uuid.UUID `params:"user_id" format:"uuid"`
} // @name DeleteUserEmailRequest
