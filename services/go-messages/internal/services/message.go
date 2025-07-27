package services

import (
	"context"
	"fmt"
	"go-messages/internal/app"
	"go-messages/internal/models"
	"log/slog"
)

type MessageService struct {
	Repo app.MongoRepository
}

func New(repo app.MongoRepository) (*MessageService, error) {

	return &MessageService{
		Repo: repo,
	}, nil
}

func (s *MessageService) HandleMessage(ctx context.Context, m models.MessageDTO) error {

	slog.Info(fmt.Sprintf("MessageService HandleMessage Id:%s, Producer:%s, Payload:%s", m.Id, m.Producer, m.Payload))
	if err := s.Repo.SaveMessage(ctx, m); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	//TODO Send response by KProducer
	return nil
}
