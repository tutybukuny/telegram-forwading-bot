package forwardmediaservice

import (
	"context"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"forwarding-bot/config"
	"forwarding-bot/internal/pkg/config"
	dbmessagehelper "forwarding-bot/internal/pkg/helper/db-message"
	"forwarding-bot/internal/pkg/helper/media-group-raw"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	"forwarding-bot/pkg/json"
	"forwarding-bot/pkg/l"
)

type IService interface {
	ProcessMessage(ctx context.Context, message *tgbotapi.Message) error
}

type serviceImpl struct {
	ll               l.Logger               `container:"name"`
	teleBot          *tgbotapi.BotAPI       `container:"name"`
	gpooling         gpooling.IPool         `container:"name"`
	db               *gorm.DB               `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo `container:"name"`

	channelConfig   serviceconfig.ChannelConfigMap
	senderMap       serviceconfig.SenderMap
	dbMessageHelper *dbmessagehelper.DBMessageHelper
}

func New(cfg *config.Config) *serviceImpl {
	service := &serviceImpl{
		senderMap:       make(map[int64]*mediagrouprawsendhelper.MediaGroupRawSendHelper),
		dbMessageHelper: dbmessagehelper.New(),
	}
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
		h := mediagrouprawsendhelper.New(channelConfig.ChannelID, channelConfig.ForwardToIDs)
		s.senderMap[h.ChannelID] = h
	}
}

func (s *serviceImpl) isMediaMessage(message *tgbotapi.Message) bool {
	return message.Video != nil || len(message.Photo) > 0 || message.Audio != nil || message.Animation != nil
}
