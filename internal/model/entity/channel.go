package entity

import "time"

type Channel struct {
	ID                 int64 `gorm:"primaryKey;autoIncrement:false"`
	LastMediaMessageID int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) IDField() string {
	return "id"
}
