// YOU CAN EDIT YOUR CUSTOM CONFIG HERE

package config

import (
	"forwarding-bot/pkg/mysql"
	telegrambot "forwarding-bot/pkg/telegram-bot"
)

// Config ...
//easyjson:json
type Config struct {
	Base         `mapstructure:",squash"`
	SentryConfig SentryConfig `json:"sentry" mapstructure:"sentry"`

	MysqlConfig       mysql.Config       `json:"mysql" mapstructure:"mysql"`
	TelegramBotConfig telegrambot.Config `json:"telegram_bot" mapstructure:"telegram_bot"`

	ChannelConfigFile string `json:"channel_config_file" mapstructure:"channel_config_file"`
	MaxPoolSize       int    `json:"max_pool_size" mapstructure:"max_pool_size"`
	RateLimit         int    `json:"rate_limit" mapstructure:"rate_limit"`
}

// SentryConfig ...
type SentryConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled"`
	DNS     string `json:"dns" mapstructure:"dns"`
	Trace   bool   `json:"trace" mapstructure:"trace"`
}
