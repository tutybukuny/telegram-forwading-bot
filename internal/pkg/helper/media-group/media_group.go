package mediagroupsendhelper

import (
	"errors"

	"forwarding-bot/internal/model/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	"forwarding-bot/pkg/l"
)

type MediaGroupSendHelper struct {
	ll       l.Logger         `container:"name"`
	teleBot  *tgbotapi.BotAPI `container:"name"`
	gpooling gpooling.IPool   `container:"name"`
}

func New() *MediaGroupSendHelper {
	h := &MediaGroupSendHelper{}
	container.Fill(h)

	return h
}

func (h *MediaGroupSendHelper) Send(channelID int64, message *entity.MediaMessage) error {
	if len(message.Messages) == 1 && message.Messages[0].Type == entity.MediaTypeAnimation {
		msg := tgbotapi.NewAnimation(channelID, tgbotapi.FileID(message.Messages[0].FileID))
		sentMsg, err := h.teleBot.Send(msg)
		if err != nil {
			h.ll.Error("error when sent message", l.Object("msg", msg), l.Error(err))
			return err
		}
		h.ll.Debug("sent message", l.Object("sent_msg", sentMsg))
	} else {
		files := make([]any, 0, len(message.Messages))
		for _, m := range message.Messages {
			fileID := tgbotapi.FileID(m.FileID)
			switch m.Type {
			case entity.MediaTypeAudio:
				files = append(files, tgbotapi.NewInputMediaAudio(fileID))
			case entity.MediaTypePhoto:
				files = append(files, tgbotapi.NewInputMediaPhoto(fileID))
			case entity.MediaTypeVideo:
				files = append(files, tgbotapi.NewInputMediaVideo(fileID))
			default:
				h.ll.Error("unhandled media type", l.String("media_type", string(m.Type)))
				return errors.New("unhandled media type")
			}
		}
		msg := tgbotapi.NewMediaGroup(channelID, files)
		sentMsg, err := h.teleBot.SendMediaGroup(msg)
		if err != nil {
			h.ll.Error("error when sent message", l.Object("msg", msg), l.Error(err))
			return err
		}
		h.ll.Debug("sent message", l.Object("sent_msg", sentMsg))
	}
	return nil
}
