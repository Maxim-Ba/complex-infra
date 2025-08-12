package main

import (
	"fmt"
	"go-websocket/internal/config"
	"go-websocket/internal/services"
	"go-websocket/pkg/kafka"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("start go-websocket")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	cfg := config.New()
	c := cfg.GetConfig()
	
	kproducer := kafka.NewProducer(cfg)
	wss := services.NewWebSoc(kproducer, c)

	messageHandler := services.NewMsgHandler(wss)

	kconsumer, err := kafka.NewConsumer(cfg, messageHandler )
	if err != nil {
		panic(err)
	}
	go kconsumer.StartRead()


	http.HandleFunc("/ws/", wss.HandleConnections)
	log.Println("http server started on ", c.ServerAddr)
	err = http.ListenAndServe(c.ServerAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	<-exit
}
