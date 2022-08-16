package helper

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ISendHelper interface {
	Send(message *tgbotapi.Message)
}
