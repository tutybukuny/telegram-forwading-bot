package messagehistoryrepo

import (
	"github.com/thnthien/impa/repository"

	"forwarding-bot/internal/model/entity"
)

type IRepo interface {
	repository.IInsert[entity.MessageHistory, int64]
}
