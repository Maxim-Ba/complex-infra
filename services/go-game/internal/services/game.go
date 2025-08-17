package services

import (
	"context"
	"encoding/json"
	"go-game/internal/models"
	"go-game/pkg/webrtc"
)

type GameService struct {
    rtcManager *webrtc.RTCManager
}

func NewGameService(rtcManager *webrtc.RTCManager) *GameService {
    return &GameService{
        rtcManager: rtcManager,
    }
}

func (s *GameService) BroadcastGameState(ctx context.Context, gameID string, state models.GameState) error {
    data, err := json.Marshal(state)
    if err != nil {
        return err
    }

    
    
    s.rtcManager.BroadcastToGame(gameID, data)
    return nil
}
