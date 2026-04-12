package main

import (
	"log"

	"unix-supremacist.github.io/textgame/internal/api"
	"unix-supremacist.github.io/textgame/internal/config"
	"unix-supremacist.github.io/textgame/internal/game"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Warning: failed to load config.json, using defaults: %v", err)
		cfg = &config.Config{
			APIKey:     "dev-key",
			ServerPort: "8080",
		}
	}

	engine := game.NewGameEngine(cfg.DataPath)
	if err := engine.LoadData(); err != nil {
		log.Printf("Warning: failed to load initial data: %v", err)
	}

	engine.RunProductionTicker()

	server := api.NewServer(engine, cfg.APIKey)
	router := server.SetupRouter()

	hostPort := ":" + cfg.ServerPort
	log.Printf("Starting TextGame API on %s", hostPort)
	if err := router.Run(hostPort); err != nil {
		log.Fatal(err)
	}
}
