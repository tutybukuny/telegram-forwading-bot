package entity

import "time"

type ChannelMessage struct {
	ID                 int64 `gorm:"primaryKey"`
	ChannelID          int64 `gorm:"uniqueIndex:channel_message"`
	MessageType        int   `gorm:"uniqueIndex:channel_message"`
	LastMediaMessageID int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (ChannelMessage) TableName() string {
	return "channel_messages"
}

func (ChannelMessage) IDField() string {
	return "id"
}
