package models

import (
	"encoding/json"
	"time"
)

type MessageDTO struct {
	PId       string    `json:"pid"` // id from provider
	Producer  string    `json:"producer"`
	Payload   string    `json:"payload"`
	Group     string    `json:"group"` // чат или комната
	CreatedAt time.Time `json:"сreatedAt"`
}
type WebRTCOffer struct {
    SDP       string `json:"sdp"`
    PlayerID  string `json:"player_id"`
    GameID    string `json:"game_id"`
    SessionID string `json:"session_id"`
}

type WebRTCAnswer struct {
    SDP       string `json:"sdp"`
    PlayerID  string `json:"player_id"`
    GameID    string `json:"game_id"`
    SessionID string `json:"session_id"`
}

type ICECandidate struct {
    Candidate string `json:"candidate"`
    PlayerID  string `json:"player_id"`
    GameID    string `json:"game_id"`
    SessionID string `json:"session_id"`
}

type WebRTCSignal struct {
    Type      string `json:"type"` // "offer", "answer", "candidate"
    Payload   json.RawMessage `json:"payload"`
}
