package webrtc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-game/internal/app"
	"go-game/internal/models"
	"log/slog"
	"strings"
	"sync"

	"github.com/pion/webrtc/v3"
)

type RTCManager struct {
	producer app.KProducer
	peers    sync.Map // map[sessionID]*PeerConnection
	config   webrtc.Configuration
	api      *webrtc.API
}

type PeerConnection struct {
	*webrtc.PeerConnection
	PlayerID string
	GameID   string
	DataChan *webrtc.DataChannel // ссылка на канал
}

var responseTopic string

func NewRTCManager(producer app.KProducer, cfg app.AppConfig) *RTCManager {
	settingEngine := webrtc.SettingEngine{}

	responseTopic = cfg.GetConfig().RTCResponseTopic
	if cfg.GetConfig().ExternalIP != "" {
		settingEngine.SetNAT1To1IPs([]string{cfg.GetConfig().ExternalIP}, webrtc.ICECandidateTypeHost)
	} else {
		settingEngine.SetInterfaceFilter(func(iface string) bool {
			return strings.HasPrefix(iface, "en") || strings.HasPrefix(iface, "eth") || strings.HasPrefix(iface, "wlan")
		})
	}

	var iceServersConfig []webrtc.ICEServer
	iceServers := cfg.GetConfig().WebRTCIceServers

	api := webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))

	for _, server := range iceServers {
		iceServersConfig = append(iceServersConfig, webrtc.ICEServer{
			URLs: []string{server},
		})
	}

	return &RTCManager{
		producer: producer,
		config: webrtc.Configuration{
			ICEServers: iceServersConfig,
		},
		api: api,
	}
}

func (m *RTCManager) HandleOffer(ctx context.Context, offer models.WebRTCOffer) error {
	peerConnection, err := m.api.NewPeerConnection(m.config)
	if err != nil {
		return fmt.Errorf("RTCManager HandleOffer NewPeerConnection: %w", err)
	}

	// save connection
	peer := &PeerConnection{
		PeerConnection: peerConnection,
		PlayerID:       offer.PlayerID,
		GameID:         offer.GameID,
	}
	m.peers.Store(offer.SessionID, peer)

	//set handler of ICE candidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			slog.Info("ICE gathering complete")
			return
		}
		slog.Info("Generated ICE candidate",
			"candidate", c.ToJSON().Candidate,
			"address", c.Address)
		candidate := models.ICECandidate{
			Candidate: c.ToJSON().Candidate,
			PlayerID:  offer.PlayerID,
			GameID:    offer.GameID,
			SessionID: offer.SessionID,
		}

		payload, _ := json.Marshal(candidate)
		signal := models.WebRTCSignal{
			Type:    "candidate",
			Payload: payload,
		}

		data, _ := json.Marshal(signal)
		m.producer.Produce("rtc-response", string(data))
	})

	// data handler
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			peer.DataChan = d // Сохраняем ссылку на канал
			d.OnOpen(func() {
				slog.Info("Data channel opened",
					"playerID", peer.PlayerID,
					"gameID", peer.GameID)
			})
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				// Обработка игровых сообщений
				slog.Info("Received game message",
					"playerID", peer.PlayerID,
					"data", string(msg.Data))

				// TODO handle game messages
			})

			d.OnClose(func() {
				slog.Info("Data channel closed")
				m.peers.Delete(offer.SessionID)
			})
		})
	})

	// set sdp
	if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.SDP,
	}); err != nil {
		return fmt.Errorf("RTCManager HandleOffer SetRemoteDescription: %w", err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return fmt.Errorf("RTCManager HandleOffer CreateAnswer: %w", err)
	}

	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return fmt.Errorf("RTCManager HandleOffer SetLocalDescription: %w", err)
	}

	response := models.WebRTCAnswer{
		SDP:       answer.SDP,
		PlayerID:  offer.PlayerID,
		GameID:    offer.GameID,
		SessionID: offer.SessionID,
	}

	payload, _ := json.Marshal(response)
	signal := models.MessageDTO{
		Action:   "answer",
		Payload:  string(payload),
		Producer: offer.PlayerID,
	}

	data, _ := json.Marshal(signal)
	return m.producer.Produce(responseTopic, string(data))
}

func (m *RTCManager) HandleAnswer(ctx context.Context, answer models.WebRTCAnswer) error {
	value, ok := m.peers.Load(answer.SessionID)
	if !ok {
		return errors.New("peer connection not found")
	}

	peer := value.(*PeerConnection)
	return peer.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer.SDP,
	})
}

func (m *RTCManager) HandleICECandidate(ctx context.Context, candidate models.ICECandidate) error {
	value, ok := m.peers.Load(candidate.SessionID)
	if !ok {
		return errors.New("peer connection not found")
	}

	peer := value.(*PeerConnection)
	return peer.AddICECandidate(webrtc.ICECandidateInit{
		Candidate: candidate.Candidate,
	})
}

func (m *RTCManager) BroadcastToGame(gameID string, data []byte) {
	m.peers.Range(func(key, value interface{}) bool {
		peer := value.(*PeerConnection)
		if peer.GameID == gameID && peer.DataChan != nil {
			if err := peer.DataChan.Send(data); err != nil {
				slog.Error("Failed to send game data",
					"playerID", peer.PlayerID,
					"error", err)
			}
		}
		return true
	})
}

func (m *RTCManager) SendToPlayer(playerID string, data []byte) error {
	var found bool
	m.peers.Range(func(key, value interface{}) bool {
		peer := value.(*PeerConnection)
		if peer.PlayerID == playerID && peer.DataChan != nil {
			if err := peer.DataChan.Send(data); err != nil {
				slog.Error("Failed to send player data",
					"playerID", playerID,
					"error", err)
				return false
			}
			found = true
			return false
		}
		return true
	})

	if !found {
		return errors.New("player not found or data channel not ready")
	}
	return nil
}
