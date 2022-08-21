package commandservice

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	commandtype "forwarding-bot/internal/constant/command-type"
	mediamessagetype "forwarding-bot/internal/constant/media-message-type"
	"forwarding-bot/internal/model/entity"
	"forwarding-bot/pkg/l"
)

var advices = []string{
	`Nữ sắc suy cho cùng cũng chỉ là máu mủ tanh hôi 😭,
Xem sex ít thôi, mà làm việc khác 😉`,
	`Nữ sắc suy cho cùng cũng chỉ là da với thịt, máu mủ tanh hôi.
Cái bẫy luân hồi đau khổ vô lượng kiếp, sa chân vào lục dục biết bao giờ mới thoát khỏi? 🥺
Đừng vì thế mà sinh lòng lưu luyến 😇`,
	`Nào các bạn, xem sex chừng mực thì quay tay mới có lực. Chậm lại đi bạn ơi 😤`,
	`Ôi bạn ơi, xem sex ít thôi, xem nhiều là liệt đấy 🥶`,
	`https://tienphong.vn/xem-phim-khieu-dam-thuong-xuyen-de-khien-dan-ong-bat-luc-post1422700.tpo`,
	`Đam mê quá nhỉ 😤, thôi, vừa vừa phai phải thôi, để người khác còn xem với 🥲`,
	`Xem nhiều quá rồi, đi nghe nhạc đi tẹo thì xem sex tiếp
https://www.youtube.com/watch?v=Llw9Q6akRo4`,
	`Xem nhiều quá rồi, đi nghe nhạc đi tẹo thì xem sex tiếp
https://www.youtube.com/watch?v=SW8zKGiTUtk`,
	`Thôi
Thôi thôi thôi 😒
nghe nhạc cho tĩnh cái tâm lại đi người https://www.youtube.com/watch?v=T6rcE2wgnSQ`,
}

func (s *serviceImpl) ProcessCommand(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case commandtype.NSFW:
		return s.handleSendNudes(ctx, message, mediamessagetype.NSFW)
	case commandtype.SFW:
		return s.handleSendNudes(ctx, message, mediamessagetype.SFW)
	default:
		s.ll.Error("not handled command", l.String("command", message.Command()))
		return errors.New("not handled command")
	}
}

func (s *serviceImpl) handleSendNudes(ctx context.Context, message *tgbotapi.Message, messageType int) error {
	_, _, _, ok, err := s.limiter.Take(ctx, fmt.Sprintf("%d", message.Chat.ID))
	if err != nil || ok {
		return s.sendNudes(ctx, message, messageType)
	}

	return s.sendAdvice(ctx, message)
}

func (s *serviceImpl) sendNudes(ctx context.Context, message *tgbotapi.Message, messageType int) error {
	channelID := message.Chat.ID
	channel, err := s.channelRepo.GetOrCreate(ctx, channelID, message.Chat.Title)
	if err != nil {
		s.ll.Error("cannot get or create channel", l.Int64("channel_id", channelID), l.Error(err))
		return err
	}

	mediaMsg, err := s.mediaMessageRepo.GetRandomMessage(ctx, channelID, messageType)
	if err != nil {
		s.ll.Error("cannot get random message", l.Object("channel", channel), l.Error(err))
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
		messageHistory := &entity.MessageHistory{
			ChannelID:      channelID,
			MediaMessageID: mediaMsg.ID,
			MessageType:    messageType,
		}
		if err = s.messageHistoryRepo.Insert(ctx, messageHistory); err != nil {
			s.ll.Error("cannot save message history", l.Object("message_history", messageHistory), l.Error(err))
		}
	}

	return err
}

func (s *serviceImpl) sendAdvice(ctx context.Context, message *tgbotapi.Message) error {
	idx := rand.Intn(len(advices))
	msg := tgbotapi.NewMessage(message.Chat.ID, advices[idx])
	s.gpooling.Submit(func() {
		s.teleBot.Send(msg)
	})
	return nil
}
