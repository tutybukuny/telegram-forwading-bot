package dbmessagehelper

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/internal/model/entity"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	"forwarding-bot/pkg/l"
)

type DBMessageHelper struct {
	ll               l.Logger               `container:"name"`
	teleBot          *tgbotapi.BotAPI       `container:"name"`
	gpooling         gpooling.IPool         `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo `container:"name"`

	ChannelID     int64
	ForwardToIDs  []int64
	mediaGroupMap map[string]chan *tgbotapi.Message
}

func New() *DBMessageHelper {
	h := &DBMessageHelper{
		mediaGroupMap: make(map[string]chan *tgbotapi.Message),
	}
	container.Fill(h)

	return h
}

func (h *DBMessageHelper) Save(ctx context.Context, message *tgbotapi.Message) error {
	if message.MediaGroupID != "" {
		messages, ok := h.mediaGroupMap[message.MediaGroupID]
		if !ok {
			messages = make(chan *tgbotapi.Message)
			h.mediaGroupMap[message.MediaGroupID] = messages
			h.gpooling.Submit(func() {
				h.saveGroupMessages(ctx, message.MediaGroupID, messages)
			})
		}
		select {
		case messages <- message:
		}
		return nil
	} else {
		files := h.buildFile(nil, message)
		return h.saveMessages(ctx, files)
	}
}

func (h *DBMessageHelper) buildFile(files []entity.Message, message *tgbotapi.Message) []entity.Message {
	if files == nil {
		files = make([]entity.Message, 0, 1)
	}

	if message.Photo != nil {
		files = append(files, entity.Message{
			Type:   entity.MediaTypePhoto,
			FileID: message.Photo[len(message.Photo)-1].FileID,
		})
	}
	if message.Video != nil {
		files = append(files, entity.Message{
			Type:   entity.MediaTypeVideo,
			FileID: message.Video.FileID,
		})
	}
	if message.Audio != nil {
		files = append(files, entity.Message{
			Type:   entity.MediaTypeAudio,
			FileID: message.Audio.FileID,
		})
	}
	if message.Animation != nil {
		files = append(files, entity.Message{
			Type:   entity.MediaTypeAnimation,
			FileID: message.Animation.FileID,
		})
	}

	return files
}

func (h *DBMessageHelper) saveMessages(ctx context.Context, messages []entity.Message) error {
	if len(messages) == 0 {
		return nil
	}

	msg := &entity.MediaMessage{Messages: messages}

	err := h.mediaMessageRepo.Insert(ctx, msg)
	if err != nil {
		h.ll.Error("cannot create media message", l.Object("msg", msg), l.Error(err))
		return err
	}
	return nil
}

func (h *DBMessageHelper) saveGroupMessages(ctx context.Context, mediaGroupID string, messages chan *tgbotapi.Message) error {
	defer delete(h.mediaGroupMap, mediaGroupID)
	defer close(messages)

	var files []entity.Message
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

loop:
	for {
		select {
		case message := <-messages:
			files = h.buildFile(files, message)
		case <-timer.C:
			break loop
		}
	}

	return h.saveMessages(ctx, files)
}
