package channelmessagerepo

import (
	"context"

	"github.com/thnthien/impa/repository"

	"forwarding-bot/internal/model/entity"
)

type IRepo interface {
	repository.IInsert[entity.ChannelMessage, int64]
	repository.IUpdate[entity.ChannelMessage, int64]

	GetOrCreate(ctx context.Context, channelID int64, messageType int) (*entity.ChannelMessage, error)
}
