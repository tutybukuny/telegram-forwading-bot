package channelstore

import (
	"context"

	gormrepository "github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"forwarding-bot/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.Channel, int64]
	*gormrepository.UpdateBaseRepo[entity.Channel, int64]
	*gormrepository.FindByIDBaseRepo[entity.Channel, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.Channel](db)
	insert := gormrepository.NewInsertBaseRepo[entity.Channel, int64](base)
	update := gormrepository.NewUpdateBaseRepo[entity.Channel, int64](base)
	find := gormrepository.NewFindByIDBaseRepo[entity.Channel, int64](base)
	return &repoImpl{base, insert, update, find}
}

func (r *repoImpl) GetOrCreate(ctx context.Context, id int64) (*entity.Channel, error) {
	channel, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if channel == nil {
		channel = &entity.Channel{
			ID:                 id,
			LastMediaMessageID: 0,
		}
		err = r.Insert(ctx, channel)
		if err != nil {
			return nil, err
		}
	}
	return channel, nil
}
