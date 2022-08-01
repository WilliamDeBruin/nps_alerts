package main

import (
	"fmt"

	"github.com/WilliamDeBruin/nps_alerts/src/config"
	"github.com/WilliamDeBruin/nps_alerts/src/server"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to load config: %s", err))
	}

	srv, err := server.NewServer(&cfg, logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize server: %s", err))
	}

	srv.Serve()
}
