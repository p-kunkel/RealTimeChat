package models

import (
	"RealTimeChat/config"
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Id          uint64 `json:"id" gorm:"type:int8;autoIncrement;primaryKey"`
	VisibleName string `json:"visible_name" gorm:"type:varchar;not null" chat:"refers_to:create;required:true"`
	LoginData
}

type LoginData struct {
	Email    string    `json:"email" gorm:"type:varchar;not null;unique" chat:"refers_to:create;required:true"`
	Password *Password `json:"password,omitempty" gorm:"column:password" chat:"refers_to:create;required:true"`
}

func (u *User) Create() error {
	return config.DB.Create(&u).Error
}

func (u *User) Login() error {
	var (
		err                 error
		pass                string
		errInvalidLoginData = errors.New("invalid email or password")
	)

	if u.Password == nil || u.Password.InputPassword == "" {
		return errInvalidLoginData
	}

	pass = u.Password.InputPassword

	if err = u.Get(func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", u.Email)
	}); err != nil {
		return err
	}

	if err = u.Password.ComparePassword([]byte(pass)); err != nil {
		return errInvalidLoginData
	}
	return nil
}

func (u *User) Get(scopes ...func(*gorm.DB) *gorm.DB) error {
	return config.DB.Scopes(scopes...).Find(&u).Error
}

func (*User) TableName() string {
	return "user"
}
