package models

import (
	"RealTimeChat/config"
	"time"

	"github.com/lib/pq"
)

type Message struct {
	ChatId    uint64        `json:"chat_id,omitempty" gorm:"type:int8;not null"`
	SenderId  uint64        `json:"sender_id,omitempty" gorm:"type:int8;not null"`
	Message   string        `json:"message,omitempty" gorm:"type:varchar(4096);not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"type:timestamptz;not null"`
	ReadedBy  pq.Int64Array `json:"readed_by" gorm:"type:int8[];default:'{}'"`

	IsReaded  bool  `json:"is_readed" gorm:"->;-:migration"`
	TotalRows int64 `json:"-" gorm:"->;-:migration"`

	Chat *ChatRoom `json:"-" gorm:"foreignKey:chat_id"`
	User *User     `json:"-" gorm:"foreignKey:sender_id"`
}

func (m *Message) Create() error {
	return config.DB.Create(&m).Error
}

func (*Message) TableName() string {
	return "message"
}
