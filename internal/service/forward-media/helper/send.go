package helper

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	"forwarding-bot/pkg/l"
)

type MediaGroupSendHelper struct {
	ll       l.Logger         `container:"name"`
	teleBot  *tgbotapi.BotAPI `container:"name"`
	gpooling gpooling.IPool   `container:"name"`

	ChannelID     int64
	ForwardToIDs  []int64
	mediaGroupMap map[string]chan *tgbotapi.Message
}

func New(channelID int64, forwardToIDs []int64) *MediaGroupSendHelper {
	h := &MediaGroupSendHelper{
		ChannelID:     channelID,
		ForwardToIDs:  forwardToIDs,
		mediaGroupMap: make(map[string]chan *tgbotapi.Message),
	}
	container.Fill(h)

	return h
}

func (h *MediaGroupSendHelper) Send(message *tgbotapi.Message) {
	if message.MediaGroupID != "" {
		messages, ok := h.mediaGroupMap[message.MediaGroupID]
		if !ok {
			messages = make(chan *tgbotapi.Message)
			h.mediaGroupMap[message.MediaGroupID] = messages
			h.gpooling.Submit(func() {
				h.sendGroupMessage(message.MediaGroupID, messages)
			})
		}
		select {
		case messages <- message:
		}
	} else {
		files := h.buildFile(nil, message)
		h.gpooling.Submit(func() {
			h.sendMessage(files)
		})
	}
}

func (h *MediaGroupSendHelper) buildFile(files []any, message *tgbotapi.Message) []any {
	if files == nil {
		files = make([]any, 0, 1)
	}

	if message.Photo != nil {
		files = append(files, tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(message.Photo[len(message.Photo)-1].FileID)))
	}
	if message.Video != nil {
		files = append(files, tgbotapi.NewInputMediaVideo(tgbotapi.FileID(message.Video.FileID)))
	}
	if message.Audio != nil {
		files = append(files, tgbotapi.NewInputMediaAudio(tgbotapi.FileID(message.Audio.FileID)))
	}

	return files
}

func (h *MediaGroupSendHelper) sendMessage(files []any) {
	if len(files) == 0 {
		return
	}
	for _, forwardToID := range h.ForwardToIDs {
		msg := tgbotapi.NewMediaGroup(forwardToID, files)
		sentMsg, err := h.teleBot.SendMediaGroup(msg)
		if err != nil {
			h.ll.Error("error when sent message", l.Object("msg", msg), l.Error(err))
		}
		h.ll.Debug("sent message", l.Object("sent_msg", sentMsg))
	}
}

func (h *MediaGroupSendHelper) sendGroupMessage(mediaGroupID string, messages chan *tgbotapi.Message) {
	defer delete(h.mediaGroupMap, mediaGroupID)
	defer close(messages)

	var files []any
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

	h.sendMessage(files)
}