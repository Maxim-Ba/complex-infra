package main

import (
	"context"
	"fmt"
	"go-game/cmd/wire"
	"go-game/internal/models"
	"go-game/internal/server"
	"go-game/internal/services"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	deps, err := wire.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Error on wire.Initialize() %v", err))
	}

	var wg sync.WaitGroup
go deps.Consumer.StartRead()

    gameService := services.NewGameService(deps.RTCManager)
    
    go func() {
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        
        for range ticker.C {
            if err := gameService.BroadcastGameState(ctx, "game1", models.GameState{
                
            }); err != nil {
                slog.Error("Failed to broadcast game state", "error", err)
            }
        }
    }()


	select {
	case <-exit:
	case <-ctx.Done():
	}

	deps.Consumer.Close()
	deps.Producer.Close()
	if err := server.Inst.Close(); err != nil {
		slog.Error(err.Error())
	}

	wg.Wait()
}
