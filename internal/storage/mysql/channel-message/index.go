package channelmessagestore

import (
	"context"
	"errors"

	gormrepository "github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"forwarding-bot/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.ChannelMessage, int64]
	*gormrepository.UpdateBaseRepo[entity.ChannelMessage, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.ChannelMessage](db)
	insert := gormrepository.NewInsertBaseRepo[entity.ChannelMessage, int64](base)
	update := gormrepository.NewUpdateBaseRepo[entity.ChannelMessage, int64](base)

	return &repoImpl{base, insert, update}
}

func (r *repoImpl) GetOrCreate(ctx context.Context, channelID int64, messageType int) (*entity.ChannelMessage, error) {
	obj := &entity.ChannelMessage{}
	err := r.GetDB(ctx).First(obj, "channel_id = ? AND message_type = ?", channelID, messageType).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		obj.ChannelID = channelID
		obj.MessageType = messageType
		obj.LastMediaMessageID = 0
		err = r.GetDB(ctx).Create(obj).Error
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}
