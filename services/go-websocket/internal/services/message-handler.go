package services

import (
	"context"
	"fmt"
	"go-websocket/internal/app"
	"go-websocket/internal/models"
)

type MessageHandler struct {
	ws app.WebSocketService
}
func NewMsgHandler( ws app.WebSocketService) *MessageHandler {
	return &MessageHandler{ws: ws}
}

func (s *MessageHandler) HandleConfirmationMessage(ctx context.Context, m models.MessageDTO) error {
	err:= s.ws.SendMessage(m)
	if err != nil {
		return fmt.Errorf("MessageHandler HandleConfirmationMessage: %w", err)
	}
	return nil
}

