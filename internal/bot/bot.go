package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/philippschepnov/market-agent/internal/research"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	engine *research.Engine
}

func NewBot(token string, engine *research.Engine) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:    api,
		engine: engine,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Я Market Research Agent. Пришли мне название компании или продукта, и я проведу SWOT-анализ на основе свежих новостей.")
			b.api.Send(reply)
		}
		return
	}

	if msg.Text == "" {
		return
	}

	company := msg.Text
	waitMsg := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Анализирую %s... Это может занять около 30 секунд.", company))
	b.api.Send(waitMsg)

	b.api.Send(tgbotapi.NewChatAction(msg.Chat.ID, tgbotapi.ChatTyping))

	result, err := b.engine.Analyze(context.Background(), company)
	if err != nil {
		log.Printf("Error analyzing %s: %v", company, err)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Извини, произошла ошибка при анализе. Попробуй позже.")
		b.api.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, result)
	reply.ParseMode = "Markdown"
	b.api.Send(reply)
}
func (b *Bot) SelfName() string {
	return b.api.Self.UserName
}
