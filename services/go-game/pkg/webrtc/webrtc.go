package webrtc

import (
	"context"
	"encoding/json"
	"errors"
	"go-game/internal/app"
	"go-game/internal/models"
	"sync"

	"github.com/pion/webrtc/v3"
)

type RTCManager struct {
    producer  app.KProducer
    peers     sync.Map // map[sessionID]*PeerConnection
    config    webrtc.Configuration
}

type PeerConnection struct {
    *webrtc.PeerConnection
    PlayerID string
    GameID   string
}
var responseTopic string

func NewRTCManager(producer app.KProducer, cfg app.AppConfig) *RTCManager {
	responseTopic = cfg.GetConfig().RTCResponseTopic
    var iceServersConfig []webrtc.ICEServer
		iceServers:= cfg.GetConfig().WebRTCIceServers
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
    }
}

func (m *RTCManager) HandleOffer(ctx context.Context, offer models.WebRTCOffer) error {
    peerConnection, err := webrtc.NewPeerConnection(m.config)
    if err != nil {
        return err
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
            return
        }

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
            // TODO handle game messages
        })
    })

    // set sdp 
    if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
        Type: webrtc.SDPTypeOffer,
        SDP:  offer.SDP,
    }); err != nil {
        return err
    }

  
    answer, err := peerConnection.CreateAnswer(nil)
    if err != nil {
        return err
    }

    
    if err = peerConnection.SetLocalDescription(answer); err != nil {
        return err
    }

    
    response := models.WebRTCAnswer{
        SDP:       answer.SDP,
        PlayerID:  offer.PlayerID,
        GameID:    offer.GameID,
        SessionID: offer.SessionID,
    }

    payload, _ := json.Marshal(response)
    signal := models.WebRTCSignal{
        Type:    "answer",
        Payload: payload,
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
