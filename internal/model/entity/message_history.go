package entity

import "time"

type MessageHistory struct {
	ID             int64 `gorm:"primaryKey"`
	ChannelID      int64 `gorm:"uniqueIndex:message_history"`
	MediaMessageID int64 `gorm:"uniqueIndex:message_history"`
	MessageType    int   `gorm:"index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (MessageHistory) TableName() string {
	return "message_histories"
}

func (MessageHistory) IDField() string {
	return "id"
}
