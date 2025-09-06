package http

import (
	"awesomeProjectDDD/internal/application"
	"fmt"
)

var ErrAuthorizationNotFound = fmt.Errorf("authorization not found: %w", application.ErrGeneral)
var ErrAccessTokenExpired = fmt.Errorf("access token expired: %w", application.ErrGeneral)
