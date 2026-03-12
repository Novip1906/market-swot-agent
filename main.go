package main

import (
	"log"

	"github.com/philippschepnov/market-agent/internal/bot"
	"github.com/philippschepnov/market-agent/internal/config"
	"github.com/philippschepnov/market-agent/internal/research"
)

func main() {
	cfg := config.Load()

	engine, err := research.NewEngine(cfg.OpenAIToken, cfg.TavilyToken, cfg.OpenAIBase, cfg.ModelName)
	if err != nil {
		log.Fatalf("Failed to initialize research engine: %v", err)
	}

	tgBot, err := bot.NewBot(cfg.TelegramToken, engine)
	if err != nil {
		log.Fatalf("Failed to initialize telegram bot: %v", err)
	}

	log.Printf("Bot started as %s", tgBot.SelfName())
	tgBot.Start()
}
