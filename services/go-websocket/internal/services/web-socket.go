package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-websocket/internal/app"
	"go-websocket/internal/models"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	KProducer   app.KProducer
	Config      app.AppConfig
	connections map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		// разрешить все origins (для разработки)
		return true
	},
}
var pingPeriod = 6

func NewWebSoc(p app.KProducer, c app.AppConfig) *WebSocketService {
	return &WebSocketService{KProducer: p, Config: c, connections: make(map[string]*websocket.Conn)}
}

func (s *WebSocketService) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// обновление соединения до WebSocket
	slog.Info("WebSocketService HandleConnections")
	producer := strings.TrimPrefix(r.URL.Path, "/ws/")

	if producer == "" {
		http.Error(w, "Producer ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Cann not upgrade websocket connection", http.StatusBadRequest)
		return
	}
	// TODO проверить есть ли producer в БД
	s.connections[producer] = conn
	slog.Info("WebSocketService HandleConnections producer is: " + producer)

	defer func() {
		// TODO проверить есть ли producer в БД - удалить из БД
		delete(s.connections, producer)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	messageTopic := s.Config.GetConfig().MessageTopic
	webRTCTopic := s.Config.GetConfig().RTCSignalTopic

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println("HandleConnections msg:",string(message) )
		var msg models.MessageDTO
		if err := json.Unmarshal(message, &msg); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Wrong message DTO"))
		}
		err = Validator.Struct(msg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Wrong message DTO: " + err.Error()))

		}
		switch msg.Action {
		case "message":
			s.KProducer.Produce(messageTopic, string(message))
		case "webrtc":
			s.KProducer.Produce(webRTCTopic, string(message))
		}

		// Обновление таймаута после успешного чтения сообщения
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
	s.ping(r.Context(), conn)
}

func (s *WebSocketService) ping(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Duration(pingPeriod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("WebSocketService ping conn.WriteMessage " + err.Error())
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *WebSocketService) SendMessage(m models.MessageDTO) error {

	conn, ok := s.connections[m.Producer]
	if ok {

		msg, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("WebSocketService SendMessage json.Marshal: %w", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return fmt.Errorf("WebSocketService SendMessage conn.WriteMessage: %w", err)
		}
	}

	return nil
}
