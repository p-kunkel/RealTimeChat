package models

import (
	dict "RealTimeChat/dictionaries"

	"errors"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ChatRooms []ChatRoom

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

func (cr *ChatRooms) FindByUserId(userId uint64, DB *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) error {
	scopes = append(scopes, func(tx *gorm.DB) *gorm.DB {
		return tx.Select("chat_room.*, (? = ANY (readed_by)) AS is_readed", userId).
			Joins("INNER JOIN chat_member AS cm ON cm.chat_id = chat_room.id").
			Where("user_id = ?", userId)
	})

	return cr.Find(DB, scopes...)
}

func (cr *ChatRooms) Find(DB *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) error {
	return DB.Scopes(scopes...).Find(&cr).Error
}

func (cr *ChatRoom) AddMembers(userIds []uint64, DB *gorm.DB) error {
	var chatMembers ChatMembers

	if cr == nil || cr.Id <= 0 {
		return errors.New("invalid chat id")
	}

	for _, uId := range userIds {
		chatMembers = append(chatMembers, ChatMember{ChatId: cr.Id, UserId: uId, RoleId: dict.Dicts.ChatRole["user"].Id})
	}

	return chatMembers.Create(DB)
}

func (cr *ChatRoom) Create(DB *gorm.DB) error {
	return DB.Create(&cr).Error
}

func (*ChatRoom) TableName() string {
	return "chat_room"
}
