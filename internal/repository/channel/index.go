package channelrepo

import (
	"context"

	"github.com/thnthien/impa/repository"

	"forwarding-bot/internal/model/entity"
)

type IRepo interface {
	repository.IInsert[entity.Channel, int64]
	repository.IFindByID[entity.Channel, int64]
	repository.IUpdate[entity.Channel, int64]

	GetOrCreate(ctx context.Context, id int64) (*entity.Channel, error)
}
