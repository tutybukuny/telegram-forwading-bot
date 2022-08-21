package mediamessagestore

import (
	"context"
	"errors"

	"github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"forwarding-bot/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.MediaMessage, int64]
	*gormrepository.FindByIDBaseRepo[entity.MediaMessage, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.MediaMessage](db)
	insert := gormrepository.NewInsertBaseRepo[entity.MediaMessage, int64](base)
	find := gormrepository.NewFindByIDBaseRepo[entity.MediaMessage, int64](base)
	return &repoImpl{base, insert, find}
}

func (r *repoImpl) GetNextMessage(ctx context.Context, lastMsgID int64, messageType int) (*entity.MediaMessage, error) {
	var obj entity.MediaMessage
	err := r.GetDB(ctx).First(&obj, "id > ? AND type = ?", lastMsgID, messageType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}
