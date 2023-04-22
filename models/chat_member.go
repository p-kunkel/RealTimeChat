package models

import (
	dict "RealTimeChat/dictionaries"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatMembers []ChatMember

type ChatMember struct {
	UserId uint64 `json:"user_id" gorm:"type:int8;not null;uniqueIndex:idx_user_chat"`
	ChatId uint64 `json:"chat_id" gorm:"type:int8;not null;uniqueIndex:idx_user_chat"`
	RoleId int8   `json:"role_id" gorm:"type:int2;not null"`

	Chat *ChatRoom       `json:"-" gorm:"foreignKey:chat_id"`
	User *User           `json:"-" gorm:"foreignKey:user_id"`
	Role *dict.DChatRole `json:"-" gorm:"foreignKey:role_id"`
}

func (cm *ChatMembers) Create(DB *gorm.DB) error {
	if cm == nil || len(*cm) == 0 {
		return errors.New("empty members list")
	}
	return DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&cm).Error
}

func (cm *ChatMember) Create(DB *gorm.DB) error {
	return DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&cm).Error
}

func (cm *ChatMembers) Find(DB *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) error {
	return DB.Find(&cm).Error
}

func (cm *ChatMember) Find(DB *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) error {
	return DB.Find(&cm).Error
}

func (*ChatMember) TableName() string {
	return "chat_member"
}
