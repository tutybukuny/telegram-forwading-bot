package dbsaverservice

import (
	"context"

	dbmessagehelper "forwarding-bot/internal/pkg/helper/db-message"
	channelrepo "forwarding-bot/internal/repository/channel"
	channelmessagerepo "forwarding-bot/internal/repository/channel-message"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	messagehistoryrepo "forwarding-bot/internal/repository/message-history"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
	"forwarding-bot/pkg/l"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IService interface {
	SaveMessage(ctx context.Context, message *tgbotapi.Message) error
}

type serviceImpl struct {
	ll                 l.Logger                 `container:"name"`
	teleBot            *tgbotapi.BotAPI         `container:"name"`
	gpooling           gpooling.IPool           `container:"name"`
	mediaMessageRepo   mediamessagerepo.IRepo   `container:"name"`
	channelRepo        channelrepo.IRepo        `container:"name"`
	channelMessageRepo channelmessagerepo.IRepo `container:"name"`
	messageHistoryRepo messagehistoryrepo.IRepo `container:"name"`

	dbMessageHelper *dbmessagehelper.DBMessageHelper
}

func New(isSaveRaw bool) *serviceImpl {
	service := &serviceImpl{
		dbMessageHelper: dbmessagehelper.New(isSaveRaw),
	}
	container.Fill(service)

	return service
}
