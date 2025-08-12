package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-transport/internal/config"
	"go-transport/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/quic-go/quic-go/http3"
	webtransport "github.com/quic-go/webtransport-go"
)

func main() {
	fmt.Println("start go-transport")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	cert, err := utils.GenerateSelfSignedCertificate()
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.New()

	c := cfg.GetConfig()

	// Создание WebTransport сервера
	wt := webtransport.Server{
		H3: http3.Server{
			Addr: c.ServerAddr, 
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				NextProtos:   []string{"h3", "webtransport"}, 
				MinVersion:   tls.VersionTLS12,
			},
		},
		CheckOrigin: func(r *http.Request) bool { return true }, // TODO заменить на проверку origin
	}
wt.H3.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Webtransport handler")

    if r.URL.Path == "/webtransport" {
		log.Println("Webtransport handler /webtransport")

			log.Println("WebTransport handler triggered")
        session, err := wt.Upgrade(w, r)
        if err != nil {
            log.Printf("Upgrade failed: %v", err)
            return
        }
        go handleWebTransportSession(session)
    } else {
        w.WriteHeader(404)
    }
})


	log.Println("Starting WebTransport server on ", c.ServerAddr)
	go func() {
		if err := wt.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-exit
}

func handleWebTransportSession(session *webtransport.Session) {
	// Принимаем входящий поток
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Printf("AcceptStream failed: %v", err)
		return
	}
	defer stream.Close()

	// Чтение сообщения
	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Printf("Read failed: %v", err)
		return
	}

	log.Printf("Received message: %s", buf[:n])

	// Отправка ответа
	_, err = stream.Write([]byte("Hello from server via WebTransport!"))
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
