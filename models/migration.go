package models

import (
	"RealTimeChat/config"
	dict "RealTimeChat/dictionaries"
	"errors"
)

func DBAutoMigrate() error {
	if config.DB == nil {
		return errors.New("not database connection")
	}
	return config.DB.AutoMigrate(
		User{},
		dict.DTokenType{},
		Token{},
		ChatRoom{},
		ChatMember{},
		Message{})
}
