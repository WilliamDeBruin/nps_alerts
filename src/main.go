package main

import (
	"fmt"

	"github.com/WilliamDeBruin/nps_alerts/src/config"
	"github.com/WilliamDeBruin/nps_alerts/src/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to load environment from env file: %s", err))
	}

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
