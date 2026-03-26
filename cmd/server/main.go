package main

import (
	"flag"
	"log"
	"os"

	"go.uber.org/zap"

	"sandbox-game/internal/app"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("SANDBOX_GAME_CONFIG")
	}
	if configPath == "" {
		configPath = "configs/local.example.yaml"
	}

	application, err := app.Bootstrap(configPath)
	if err != nil {
		log.Fatalf("bootstrap sandbox game application failed: %v", err)
	}

	addr := application.Config.Server.Address()
	application.Logger.Info("sandbox game server starting",
		zap.String("addr", addr),
		zap.String("config", configPath),
	)

	if err := application.Router.Run(addr); err != nil {
		application.Logger.Fatal("sandbox game server stopped", zap.Error(err))
	}
}
