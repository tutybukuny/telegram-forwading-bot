package main

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/config"
	processservice "forwarding-bot/internal/service/forward-media"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	handleossignal "forwarding-bot/pkg/handle-os-signal"
	"forwarding-bot/pkg/l"
	telegrambot "forwarding-bot/pkg/telegram-bot"
	validator "forwarding-bot/pkg/validator"
)

func bootstrap(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	_, cancel := context.WithCancel(context.Background())
	shutdown.HandleDefer(cancel)

	container.NamedSingleton("gpooling", func() gpooling.IPool {
		return gpooling.New(cfg.MaxPoolSize, ll)
	})

	container.NamedSingleton("validator", func() validator.IValidator {
		return validator.New()
	})

	//region init agent
	container.NamedSingleton("teleBot", func() *tgbotapi.BotAPI {
		return telegrambot.New(cfg.TelegramBotConfig)
	})
	//endregion

	//region init service
	container.NamedSingleton("processService", func() processservice.IService {
		return processservice.New(cfg)
	})
	//endregion
}
