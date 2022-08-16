package telegrambot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	ApiKey  string `json:"api_key,omitempty" mapstructure:"api_key"`
	TimeOut int    `json:"time_out,omitempty" mapstructure:"time_out"`
	Debug   bool   `json:"debug,omitempty" mapstructure:"debug"`
}

func New(cfg Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(cfg.ApiKey)
	if err != nil {
		log.Fatalf("cannot create bot: %v", err)
	}
	bot.Debug = cfg.Debug
	return bot
}
