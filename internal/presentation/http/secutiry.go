package http

import (
	"awesomeProjectDDD/internal/adapter/hydra"
	"awesomeProjectDDD/internal/observability"
	"github.com/gofiber/fiber/v2"
	client "github.com/ory/hydra-client-go"
)

type SecurityMiddleware struct {
	observer     *observability.Observability
	hydraService *hydra.Service
}

func NewSecurityMiddleware(
	observer *observability.Observability,
	hydraService *hydra.Service,
) *SecurityMiddleware {
	return &SecurityMiddleware{
		observer:     observer,
		hydraService: hydraService,
	}
}

func (s SecurityMiddleware) Handle(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		s.observer.Logger.Error().Msg("authorization not found")
		return RespondError(c, ErrAuthorizationNotFound)
	}

	identity, err := s.hydraService.Introspect(c.UserContext(), authHeader)

	if err != nil {
		s.observer.Logger.Error().Err(err).Msg("failed to introspect token")
		return RespondError(c, ErrAuthorizationNotFound)
	}

	if !identity.Active {
		s.observer.Logger.Error().Err(err).Msg("token is not active")
		return RespondError(c, ErrAccessTokenExpired)
	}

	c.Locals("identity", identity)
	return c.Next()
}

func GetIdentity(c *fiber.Ctx) (*client.OAuth2TokenIntrospection, error) {
	identity, ok := c.Locals("identity").(*client.OAuth2TokenIntrospection)
	if !ok {
		return nil, ErrAuthorizationNotFound
	}
	return identity, nil
}
