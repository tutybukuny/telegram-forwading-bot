package forwardmediaservice

import (
	"context"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/config"
	"forwarding-bot/internal/service/forward-media/helper"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/json"
	"forwarding-bot/pkg/l"
)

type ChannelConfigMap map[string]*ChannelConfig
type SenderMap map[int64]helper.ISendHelper

type ChannelConfig struct {
	ChannelID    int64   `json:"channel_id,omitempty"`
	ForwardToIDs []int64 `json:"forward_to_ids,omitempty"`
}

type IService interface {
	ProcessMessage(ctx context.Context, message *tgbotapi.Message)
}

type serviceImpl struct {
	ll      l.Logger         `container:"name"`
	teleBot *tgbotapi.BotAPI `container:"name"`

	channelConfig ChannelConfigMap
	senderMap     SenderMap
}

func New(cfg *config.Config) *serviceImpl {
	service := &serviceImpl{senderMap: make(map[int64]helper.ISendHelper)}
	container.Fill(service)

	service.readConfig(cfg)
	service.initHelpers()

	return service
}

func (s *serviceImpl) readConfig(cfg *config.Config) {
	channelConfigFile := cfg.ChannelConfigFile
	if channelConfigFile == "" {
		channelConfigFile = "./config.json"
	}
	file, err := os.Open(cfg.ChannelConfigFile)
	if err != nil {
		s.ll.Fatal("cannot read channel config", l.String("channel_config_file", cfg.ChannelConfigFile), l.Error(err))
	}
	defer file.Close()
	configJson, err := ioutil.ReadAll(file)
	if err != nil {
		s.ll.Fatal("cannot read channel config", l.String("channel_config_file", cfg.ChannelConfigFile), l.Error(err))
	}
	err = json.Unmarshal(configJson, &s.channelConfig)
	if err != nil {
		s.ll.Fatal("cannot parse channel config",
			l.String("channel_config_file", cfg.ChannelConfigFile),
			l.ByteString("config_json", configJson), l.Error(err))
	}
	s.ll.Info("loaded channel config", l.Object("channel_config", s.channelConfig))
}

func (s *serviceImpl) initHelpers() {
	for _, channelConfig := range s.channelConfig {
		h := helper.New(channelConfig.ChannelID, channelConfig.ForwardToIDs)
		s.senderMap[h.ChannelID] = h
	}
}

func (s *serviceImpl) isMediaMessage(message *tgbotapi.Message) bool {
	return message.Video != nil || len(message.Photo) > 0 || message.Audio != nil
}
