package commandservice

import (
	"context"
	"errors"
	commandtype "forwarding-bot/internal/constant/command-type"
	"forwarding-bot/pkg/l"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) ProcessCommand(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case commandtype.SendNudes:
		return s.handleSendNudes(ctx, message)
	default:
		s.ll.Error("not handled command", l.String("command", message.Command()))
		return errors.New("not handled command")
	}
}

func (s *serviceImpl) handleSendNudes(ctx context.Context, message *tgbotapi.Message) error {
	channelID := message.Chat.ID
	channel, err := s.channelRepo.GetOrCreate(ctx, channelID)
	if err != nil {
		s.ll.Error("cannot get or create channel", l.Int64("channel_id", channelID), l.Error(err))
		return err
	}
	nextMessageID := channel.LastMediaMessageID + 1
	mediaMsg, err := s.mediaMessageRepo.FindByID(ctx, nextMessageID)
	if err != nil {
		s.ll.Error("cannot get next message", l.Object("channel", channel), l.Error(err))
		return err
	}
	if mediaMsg == nil {
		s.gpooling.Submit(func() {
			msg := tgbotapi.NewMessage(channelID, "Ôi các bạn ơi, xem sex ít thôi, làm gì có nhiều sex thế để mà các bạn xem \U0001F979")
			sentMsg, err := s.teleBot.Send(msg)
			if err != nil {
				s.ll.Error("error when sent message", l.Object("msg", msg), l.Error(err))
			}
			s.ll.Debug("sent message", l.Object("sent_msg", sentMsg))

		})
	} else {
		err = s.sendHelper.Send(channelID, mediaMsg)
		if err != nil {
			return err
		}
		channel.LastMediaMessageID = mediaMsg.ID
		if err = s.channelRepo.Update(ctx, channel); err != nil {
			s.ll.Error("cannot save channel", l.Object("channel", channel), l.Error(err))
		}
	}

	return err
}
