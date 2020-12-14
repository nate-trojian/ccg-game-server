package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nate-trojian/ccg-game-server/pkg"
	"go.uber.org/zap"
)

func main() {
	_ = pkg.InitializeLogger("debug")
	logger := pkg.NewLogger("main")

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	server := pkg.NewServer()
	go server.Start()

	<-stop
	if err := server.Shutdown(); err != nil {
		logger.Error("Server shutdown encountered error", zap.Error(err))
	}
}