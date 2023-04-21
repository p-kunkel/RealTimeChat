package models

import (
	"RealTimeChat/config"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ChatRoom struct {
	Id              uint64         `json:"id" gorm:"type:int8;autoIncrement;primaryKey"`
	LastMessage     string         `json:"last_message,omitempty" gorm:"type:varchar(1000);default:null"`
	MessageSenderId uint64         `json:"message_sender_id,omitempty" gorm:"type:int8;default:null"`
	ReadedBy        pq.Int64Array  `json:"readed_by" gorm:"type:int8[];default:'{}'"`
	CreatedAt       *time.Time     `json:"-" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time     `json:"updated_at,omitempty" gorm:"type:timestamp"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"type:timestamp"`

	IsReaded  bool           `json:"is_readed" gorm:"->;-:migration"`
	TotalRows int64          `json:"-" gorm:"->;-:migration" swaggerignore:"true"`
	Langs     pq.StringArray `json:"-" gorm:"->;-:migration;type:varchar[]"`

	Sender *User `json:"-" gorm:"foreignKey:message_sender_id"`
}

func (cr *ChatRoom) AddMembers(userIds []uint64) error {
	var chatMembers ChatMembers

	for _, uId := range userIds {
		chatMembers = append(chatMembers, ChatMember{ChatId: cr.Id, UserId: uId})
	}

	return chatMembers.Create()
}

func (cr *ChatRoom) Create() error {
	return config.DB.Create(&cr).Error
}

func (*ChatRoom) TableName() string {
	return "chat_room"
}
