package main

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"forwarding-bot/config"
	"forwarding-bot/internal/model/entity"
	channelrepo "forwarding-bot/internal/repository/channel"
	channelmessagerepo "forwarding-bot/internal/repository/channel-message"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	commandservice "forwarding-bot/internal/service/command"
	processservice "forwarding-bot/internal/service/forward-media"
	channelstore "forwarding-bot/internal/storage/mysql/channel"
	channelmessagestore "forwarding-bot/internal/storage/mysql/channel-message"
	mediamessagestore "forwarding-bot/internal/storage/mysql/media-message"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	handleossignal "forwarding-bot/pkg/handle-os-signal"
	"forwarding-bot/pkg/l"
	"forwarding-bot/pkg/mysql"
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

	//region init store
	db := mysql.New(cfg.MysqlConfig, ll)
	mysql.AutoMigration(db, []any{
		&entity.MediaMessage{}, &entity.Channel{}, &entity.ChannelMessage{},
	}, ll)

	container.NamedSingleton("db", func() *gorm.DB {
		return db
	})

	container.NamedSingleton("mediaMessageRepo", func() mediamessagerepo.IRepo {
		return mediamessagestore.New(db)
	})

	container.NamedSingleton("channelRepo", func() channelrepo.IRepo {
		return channelstore.New(db)
	})

	container.NamedSingleton("channelMessageRepo", func() channelmessagerepo.IRepo {
		return channelmessagestore.New(db)
	})
	//endregion

	//region init agent
	container.NamedSingleton("teleBot", func() *tgbotapi.BotAPI {
		return telegrambot.New(cfg.TelegramBotConfig)
	})
	//endregion

	//region init service
	container.NamedSingleton("processService", func() processservice.IService {
		return processservice.New(cfg)
	})
	container.NamedSingleton("commandService", func() commandservice.IService {
		return commandservice.New(cfg.RateLimit)
	})
	//endregion
}
