package commandservice

import (
	"context"
	channelrepo "forwarding-bot/internal/repository/channel"
	"forwarding-bot/pkg/gpooling"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	mediagroupsendhelper "forwarding-bot/internal/pkg/helper/media-group"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/l"
)

type IService interface {
	ProcessCommand(ctx context.Context, message *tgbotapi.Message) error
}

type serviceImpl struct {
	ll               l.Logger               `container:"name"`
	teleBot          *tgbotapi.BotAPI       `container:"name"`
	gpooling         gpooling.IPool         `container:"name"`
	mediaMessageRepo mediamessagerepo.IRepo `container:"name"`
	channelRepo      channelrepo.IRepo      `container:"name"`

	sendHelper *mediagroupsendhelper.MediaGroupSendHelper
}

func New() *serviceImpl {
	service := &serviceImpl{
		sendHelper: mediagroupsendhelper.New(),
	}
	container.Fill(service)

	return service
}
