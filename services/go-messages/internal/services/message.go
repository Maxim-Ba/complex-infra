package services

import (
	"fmt"
	"go-messages/internal/models"
	"log/slog"
)

type MessageService struct {
}

func New() (*MessageService, error) {

	return &MessageService{}, nil
}

func (s *MessageService) HandleMessage(m models.MessageDTO) error {

	slog.Info(fmt.Sprintf("MessageService HandleMessage Id:%s, Producer:%s, Payload:%s", m.Id,m.Producer, m.Payload))
// TODO save DB
// Send response by Producer
	return  nil
}
