package commandservice

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"

	mediagroupsendhelper "forwarding-bot/internal/pkg/helper/media-group"
	channelrepo "forwarding-bot/internal/repository/channel"
	mediamessagerepo "forwarding-bot/internal/repository/media-message"
	"forwarding-bot/pkg/container"
	"forwarding-bot/pkg/gpooling"
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
	limiter    limiter.Store
}

func New(rate int) *serviceImpl {
	li, err := memorystore.New(&memorystore.Config{
		Tokens:   uint64(rate),
		Interval: time.Minute,
	})

	if err != nil {
		log.Fatalf("cannot create limiter: %v", err)
	}

	service := &serviceImpl{
		sendHelper: mediagroupsendhelper.New(),
		limiter:    li,
	}
	container.Fill(service)

	return service
}
