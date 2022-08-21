package messagehistorystore

import (
	gormrepository "github.com/thnthien/impa/repository/gorm"
	"gorm.io/gorm"

	"forwarding-bot/internal/model/entity"
)

type repoImpl struct {
	*gormrepository.BaseRepo
	*gormrepository.InsertBaseRepo[entity.MessageHistory, int64]
}

func New(db *gorm.DB) *repoImpl {
	base := gormrepository.NewBaseRepo[entity.MessageHistory](db)
	insert := gormrepository.NewInsertBaseRepo[entity.MessageHistory, int64](base)
	return &repoImpl{base, insert}
}
