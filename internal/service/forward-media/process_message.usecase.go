package forwardmediaservice

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	serviceconfig "forwarding-bot/internal/pkg/config"
	"forwarding-bot/pkg/l"
)

func (s *serviceImpl) ProcessMessage(ctx context.Context, message *tgbotapi.Message) error {
	if !s.isMediaMessage(message) {
		s.ll.Info("this is not media message, just ignore it")
		return nil
	}

	channelConfig, ok := s.channelConfig[fmt.Sprintf("%d", message.Chat.ID)]
	if !ok {
		s.ll.Info("ignore message from unconfigured channel", l.Int64("channel_id", message.Chat.ID))
		return nil
	}

	if channelConfig.IsDBChannel {
		return nil
		return s.handleDBChannel(ctx, channelConfig, message)
	}

	return s.handleForwardingChannel(ctx, channelConfig, message)
}

func (s *serviceImpl) handleDBChannel(ctx context.Context, channelConfig *serviceconfig.ChannelConfig, message *tgbotapi.Message) error {
	return s.dbMessageHelper.Save(ctx, message)
}

func (s *serviceImpl) handleForwardingChannel(ctx context.Context, channelConfig *serviceconfig.ChannelConfig, message *tgbotapi.Message) error {
	sender, ok := s.senderMap[channelConfig.ChannelID]
	if !ok {
		s.ll.Error("cannot find sender", l.Object("channel_config", channelConfig), l.Object("sender_map", s.senderMap))
		return errors.New("cannot find sender")
	}

	sender.Send(message)
	return nil
}
