package webrtc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-game/internal/app"
	"go-game/internal/models"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pion/turn/v2"
	"github.com/pion/webrtc/v3"
)

type RTCManager struct {
	producer app.KProducer
	peers    sync.Map // map[sessionID]*PeerConnection
	config   webrtc.Configuration
	api      *webrtc.API
			turnAuth *TurnAuthenticator

}

type PeerConnection struct {
	*webrtc.PeerConnection
	PlayerID string
	GameID   string
	DataChan *webrtc.DataChannel // ссылка на канал

}
type TurnAuthenticator struct {
	username string
	password string
}

func (a *TurnAuthenticator) Authenticate(username, realm string, srcAddr net.Addr) ([]byte, bool) {
	key := turn.GenerateAuthKey(username, realm, a.password)
	return key, username == a.username
}
var responseTopic string

func NewRTCManager(producer app.KProducer, cfg app.AppConfig) *RTCManager {
	settingEngine := webrtc.SettingEngine{}

turnAuth := &TurnAuthenticator{
		username: "user",
		password: "password",
	}

	responseTopic = cfg.GetConfig().RTCResponseTopic

if cfg.GetConfig().ExternalIP != "" {
		settingEngine.SetNAT1To1IPs([]string{cfg.GetConfig().ExternalIP}, webrtc.ICECandidateTypeHost)
		
		// Установка диапазона портов для TURN релеев
		settingEngine.SetEphemeralUDPPortRange(50000, 50100)
	} else {
		// Фильтр интерфейсов для разработки
		settingEngine.SetInterfaceFilter(func(iface string) bool {
			return strings.HasPrefix(iface, "en") || 
				   strings.HasPrefix(iface, "eth") || 
				   strings.HasPrefix(iface, "wlan")||
					  strings.HasPrefix(iface, "br-")
		})
	}

	// Включение более агрессивной сборки ICE кандидатов
	settingEngine.SetLite(false)
	settingEngine.SetNetworkTypes([]webrtc.NetworkType{
		webrtc.NetworkTypeUDP4,
		webrtc.NetworkTypeTCP4,
	})

	api := webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))

	iceServersConfig := []webrtc.ICEServer{
        {
            URLs: []string{"stun:coturn:3478"}, // Используем имя контейнера
        },
        {
            URLs:       []string{"turn:coturn:3478?transport=udp"},
            Username:   "user",
            Credential: "password",
        },
        {
            URLs:       []string{"turn:coturn:3478?transport=tcp"},
            Username:   "user",
            Credential: "password",
        },
    }

	return &RTCManager{
		producer: producer,
		config: webrtc.Configuration{
			ICEServers:   iceServersConfig,
			SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		},
		api:      api,
		turnAuth: turnAuth,
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
// Установка таймаутов для ICE
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// Канал для отслеживания завершения ICE gathering
	iceGatheringComplete := make(chan struct{}, 1)
	//set handler of ICE candidate
//set handler of ICE candidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			slog.Info("ICE gathering complete")
			select {
			case iceGatheringComplete <- struct{}{}:
			default:
			}
			return
		}

		// Фильтруем локальные кандидаты, которые не будут работать через NAT
		candidateJSON := c.ToJSON()
		if strings.Contains(candidateJSON.Candidate, "192.168.") || 
		   strings.Contains(candidateJSON.Candidate, "172.") ||
		   strings.Contains(candidateJSON.Candidate, "10.") ||
		   strings.Contains(candidateJSON.Candidate, "127.0.0.1") ||
		   strings.Contains(candidateJSON.Candidate, "localhost") {
			slog.Debug("Skipping private IP candidate", "candidate", candidateJSON.Candidate)
			return
		}

		slog.Info("Generated ICE candidate",
			"candidate", candidateJSON.Candidate,
			"address", c.Address,
		)

		candidate := models.ICECandidate{
			Candidate: candidateJSON.Candidate,
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

	// Обработка состояния соединения
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		slog.Info("Connection state changed", 
			"state", s.String(),
			"playerID", offer.PlayerID,
			"gameID", offer.GameID)
		
		if s == webrtc.PeerConnectionStateFailed {
			// Попытка восстановления соединения
			slog.Warn("Connection failed, attempting to restart ICE")
			go m.restartICE(offer.SessionID)
		} else if s == webrtc.PeerConnectionStateClosed {
			m.peers.Delete(offer.SessionID)
		}
	})

	// data handler
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
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
		})

		d.OnClose(func() {
			slog.Info("Data channel closed")
			m.peers.Delete(offer.SessionID)
		})
		
		d.OnError(func(err error) {
			slog.Error("Data channel error", 
				"error", err,
				"playerID", peer.PlayerID)
		})
	})

	// set sdp
	if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.SDP,
	}); err != nil {
		return fmt.Errorf("RTCManager HandleOffer SetRemoteDescription: %w", err)
	}

	// Создаем answer с настройками для лучшей совместимости
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return fmt.Errorf("RTCManager HandleOffer CreateAnswer: %w", err)
	}

	// Ждем завершения ICE gathering перед отправкой answer
	select {
	case <-iceGatheringComplete:
	case <-ctx.Done():
		slog.Warn("ICE gathering timeout", "sessionID", offer.SessionID)
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
func (m *RTCManager) restartICE(sessionID string) {
	value, ok := m.peers.Load(sessionID)
	if !ok {
		return
	}

	peer := value.(*PeerConnection)
	
	// Создаем новый offer для инициации restart ICE
	offer, err := peer.CreateOffer(nil)
	if err != nil {
		slog.Error("Failed to create restart offer", "error", err)
		return
	}

	if err := peer.SetLocalDescription(offer); err != nil {
		slog.Error("Failed to set local description for restart", "error", err)
		return
	}

	// Отправляем новый offer через signaling
	// (реализация зависит от вашей signaling инфраструктуры)
}
