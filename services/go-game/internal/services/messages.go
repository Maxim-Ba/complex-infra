package services

import (
	"context"
	"encoding/json"
	"go-game/internal/models"
	"go-game/pkg/webrtc"
	"log/slog"
)

type MessageService struct {
    rtcManager *webrtc.RTCManager
}

func NewMessageService(rtcManager *webrtc.RTCManager) *MessageService {
    return &MessageService{
        rtcManager: rtcManager,
    }
}

func (s *MessageService) HandleMessage(ctx context.Context, msg models.MessageDTO) error {
    var signal models.WebRTCSignal
    if err := json.Unmarshal([]byte(msg.Payload), &signal); err != nil {
        return err
    }

    switch signal.Type {
    case "offer":
        var offer models.WebRTCOffer
        if err := json.Unmarshal(signal.Payload, &offer); err != nil {
            return err
        }
        return s.rtcManager.HandleOffer(ctx, offer)
    case "answer":
        var answer models.WebRTCAnswer
        if err := json.Unmarshal(signal.Payload, &answer); err != nil {
            return err
        }
        return s.rtcManager.HandleAnswer(ctx, answer)
    case "candidate":
        var candidate models.ICECandidate
        if err := json.Unmarshal(signal.Payload, &candidate); err != nil {
            return err
        }
        return s.rtcManager.HandleICECandidate(ctx, candidate)
    default:
        slog.Warn("Unknown WebRTC signal type", "type", signal.Type)
        return nil
    }
}
