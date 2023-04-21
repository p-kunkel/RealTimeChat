package models

import (
	"RealTimeChat/config"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatMembers []ChatMember

type ChatMember struct {
	UserId uint64 `json:"user_id" gorm:"type:int8;not null;uniqueIndex:idx_user_chat"`
	ChatId uint64 `json:"chat_id" gorm:"type:int8;not null;uniqueIndex:idx_user_chat"`
	RoleId int8   `json:"role_id" gorm:"type:int2;not null"`

	Chat *ChatRoom `json:"-" gorm:"foreignKey:chat_id"`
	User *User     `json:"-" gorm:"foreignKey:user_id"`
}

func (cm *ChatMembers) Create() error {
	return config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&cm).Error
}

func (cm *ChatMember) Create() error {
	return config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&cm).Error
}

func (cm *ChatMembers) Find(scopes ...func(*gorm.DB) *gorm.DB) error {
	return config.DB.Find(&cm).Error
}

func (cm *ChatMember) Find(scopes ...func(*gorm.DB) *gorm.DB) error {
	return config.DB.Find(&cm).Error
}

func (*ChatMember) TableName() string {
	return "chat_member"
}
