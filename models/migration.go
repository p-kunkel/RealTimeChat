package models

import (
	"RealTimeChat/config"
	"errors"
)

func DBAutoMigrate() error {
	if config.DB == nil {
		return errors.New("not database connection")
	}
	return config.DB.AutoMigrate(User{})
}
