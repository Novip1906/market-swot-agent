package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	OpenAIToken   string
	TavilyToken   string
	OpenAIBase    string
	ModelName     string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_APITOKEN"),
		OpenAIToken:   os.Getenv("OPENAI_API_KEY"),
		TavilyToken:   os.Getenv("TAVILY_API_KEY"),
		OpenAIBase:    os.Getenv("OPENAI_API_BASE"),
		ModelName:     os.Getenv("MODEL_NAME"),
	}

	if cfg.ModelName == "" {
		cfg.ModelName = "gpt-4o"
	}

	if cfg.TelegramToken == "" || cfg.OpenAIToken == "" || cfg.TavilyToken == "" {
		log.Println("Warning: some API keys are missing in environment")
	}

	return cfg
}
