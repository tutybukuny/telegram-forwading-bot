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

func (r *repoImpl) GetRandomMessage(ctx context.Context, channelID int64, messageType int) (*entity.MediaMessage, error) {
	query := `SELECT m.* FROM media_messages AS m 
           		WHERE m.type = ? AND m.id NOT IN (SELECT mm.id FROM media_messages AS mm 
           		                                 INNER JOIN message_histories AS mh ON mm.id = mh.media_message_id 
           		                                                                           AND mh.channel_id = ? 
           		                                                                           AND mh.message_type = ?)
           		ORDER BY RAND() LIMIT 1`
	var mediaMessage entity.MediaMessage
	err := r.GetDB(ctx).Raw(query, messageType, channelID, messageType).Scan(&mediaMessage).Error
	if err != nil {
		return nil, err
	}
	if mediaMessage.ID == 0 {
		return nil, nil
	}
	return &mediaMessage, nil
}
