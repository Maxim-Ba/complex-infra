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

	slog.Info(fmt.Sprintf("MessageService HandleMessage Id:%s, Producer:%s, Payload:%s", m.PId, m.Producer, m.Payload))
	if err := s.Repo.SaveMessage(ctx, m); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	msgBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("MessageService HandleMessage failed to marshal message to JSON: %w", err)
	}
	msgStr := string(msgBytes)

	// TODO записать createdAt из БД
	s.Producer.Produce("message_confirmations", msgStr)
	return nil
}
func (s *MessageService) Get(ctx context.Context, r models.RequestMessages) ([]models.MessageDTO, error) {
    slog.Info("MessageService Get", "group", r.GroupiD, "offset", r.Offset, "count", r.Count)
    
    // Валидация параметров запроса
    if r.GroupiD == "" {
        return nil, fmt.Errorf("group ID cannot be empty")
    }
    
    if r.Offset < 0 {
        r.Offset = 0
    }
    
    if r.Count <= 0 {
        r.Count = 10 // Значение по умолчанию
    } else if r.Count > 100 {
        r.Count = 100 // Ограничение максимального количества
    }
    
    messages, err := s.Repo.GetMessagesByGroup(ctx, r)
    if err != nil {
        return nil, fmt.Errorf("failed to get messages: %w", err)
    }
    
    return messages, nil
}
