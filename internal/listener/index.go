package listener

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"forwarding-bot/config"
	"forwarding-bot/internal/pkg/middleware"
	commandservice "forwarding-bot/internal/service/command"
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
	commandService commandservice.IService `container:"name"`
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
	ctx := context.Background()
	var db *gorm.DB
	container.NamedResolve(&db, "db")
	bot, err := tl.teleBot.GetMe()
	if err != nil {
		tl.ll.Fatal("cannot get bot info", l.Error(err))
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message
		tl.ll.Debug("received message", l.Object("message", message))

		if message.IsCommand() {
			if strings.Contains(message.Text, "@") {
				at := strings.Split(message.CommandWithAt(), "@")[1]
				if at != "" && at != bot.UserName {
					tl.ll.Debug("command for other bots, ignore")
					continue
				}
			}
			err = middleware.NewGormTransaction(db, ctx, func(ctx context.Context) error {
				return tl.commandService.ProcessCommand(ctx, message)
			})
		} else {
			err = tl.processService.ProcessMessage(ctx, message)
		}

		if err != nil {
			tl.ll.Error("error when handle message", l.Object("message", message), l.Error(err))
		}
	}
}
