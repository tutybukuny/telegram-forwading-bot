package forwardmediaservice

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/pkg/l"
)

func (s *serviceImpl) ProcessMessage(ctx context.Context, message *tgbotapi.Message) {
	if !s.isMediaMessage(message) {
		s.ll.Info("this is not media message, just ignore it")
		return
	}

	channelConfig, ok := s.channelConfig[fmt.Sprintf("%d", message.Chat.ID)]
	if !ok {
		s.ll.Info("ignore message from unconfigured channel", l.Int64("channel_id", message.Chat.ID))
		return
	}

	sender, ok := s.senderMap[channelConfig.ChannelID]
	if !ok {
		s.ll.Error("cannot find sender", l.Object("channel_config", channelConfig), l.Object("sender_map", s.senderMap))
		return
	}

	sender.Send(message)
}
