package entity

import "time"

type Channel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) IDField() string {
	return "id"
}
