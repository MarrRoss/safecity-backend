package hydra

import (
	"awesomeProjectDDD/pkg/http"
	"context"
	"fmt"
	client "github.com/ory/hydra-client-go"
)

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Introspect(ctx context.Context, token string) (*client.OAuth2TokenIntrospection, error) {
	formValues := map[string]string{
		"token": token,
	}
	var successModel client.OAuth2TokenIntrospection
	var errModel client.JsonError
	res, err := s.client.R().
		SetContext(ctx).
		SetFormData(formValues).
		SetResult(&successModel).
		SetError(&errModel).
		Post("/admin/oauth2/introspect")

	if err != nil {
		return nil, err
	}
	if (200 <= res.StatusCode()) && (res.StatusCode() >= 300) {
		return nil, fmt.Errorf("error: %s", errModel.GetError())
	}
	return &successModel, nil
}
