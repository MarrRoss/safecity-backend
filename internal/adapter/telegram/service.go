package telegram

import (
	"awesomeProjectDDD/pkg/http"
	"context"
)

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) SendMessage(ctx context.Context, telegramID string, message string) error {
	data := map[string]string{
		"chat_id": telegramID,
		"text":    message,
	}
	_, err := s.client.R().
		SetContext(ctx).
		SetBody(data).
		Post("/sendMessage")

	if err != nil {
		return err
	}
	//if (200 <= res.StatusCode()) && (res.StatusCode() >= 300) {
	//	return nil, fmt.Errorf("error: %s", errModel.GetError())
	//}
	return nil
}
