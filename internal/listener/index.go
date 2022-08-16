package listener

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/config"
	processservice "forwarding-bot/internal/service/forward-media"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/l"
	telegrambot "forwarding-bot/pkg/telegram-bot"
)

type TelegramListener struct {
	ll         l.Logger `container:"name"`
	teleConfig telegrambot.Config

	teleBot *tgbotapi.BotAPI `container:"name"`

	processService processservice.IService `container:"name"`
}

func New(cfg *config.Config) *TelegramListener {
	listener := &TelegramListener{}
	container.Fill(listener)

	return listener
}

func (tl *TelegramListener) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = tl.teleConfig.TimeOut
	if u.Timeout == 0 {
		u.Timeout = 60
	}

	updates := tl.teleBot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message
		tl.ll.Info("received message", l.Object("message", message))

		tl.processService.ProcessMessage(context.Background(), message)
	}
}
