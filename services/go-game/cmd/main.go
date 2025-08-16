package main

import (
	"context"
	"go-game/internal/config"
	"go-game/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	cfg := config.New()

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()

		if err := server.Start(cfg.ServerAddr); err != nil {
			cancel()
		}
	}()

	select {
	case <-exit:
	case <-ctx.Done():
	}

	if err := server.Inst.Close(); err != nil {
		slog.Error(err.Error())
	}

	wg.Wait()
}
