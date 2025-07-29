package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-messages/internal/app"
	"go-messages/internal/models"
	"log/slog"
)

type MessageService struct {
	Repo     app.MongoRepository
	Producer app.KProducer
}

func New(repo app.MongoRepository, producer app.KProducer) (*MessageService, error) {

	return &MessageService{
		Repo:     repo,
		Producer: producer,
	}, nil
}

func (s *MessageService) HandleMessage(ctx context.Context, m models.MessageDTO) error {

	slog.Info(fmt.Sprintf("MessageService HandleMessage Id:%s, Producer:%s, Payload:%s", m.Id, m.Producer, m.Payload))
	if err := s.Repo.SaveMessage(ctx, m); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	msgBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("MessageService HandleMessage failed to marshal message to JSON: %w", err)
	}
	msgStr := string(msgBytes)

	s.Producer.Produce("message_confirmations", msgStr)
	return nil
}
