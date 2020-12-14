package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nate-trojian/ccg-game-server/internal"
	"github.com/nate-trojian/ccg-game-server/pkg"
	"github.com/nate-trojian/ccg-game-server/pkg/matchmaking"
	"go.uber.org/zap"
)

func main() {
	_ = internal.InitializeLogger("debug")
	logger := internal.NewLogger("main")

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())

	matchmaker := matchmaking.InitializeMatchmaker()
	go matchmaker.Start(ctx)

	hub := pkg.NewHub(matchmaker.Out())
	go hub.Start(ctx)

	server := pkg.NewServer(hub.RegisterClientsChan())
	go server.Start()

	<-stop
	cancel()
	// TODO - Should wait until all sub components have finished
	if err := server.Shutdown(); err != nil {
		logger.Error("Server shutdown encountered error", zap.Error(err))
	}
}