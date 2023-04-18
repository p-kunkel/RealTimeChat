package models

import "RealTimeChat/config"

type User struct {
	Id          uint64    `json:"id" gorm:"type:bigserial;primarykey"`
	Email       string    `json:"email" gorm:"type:varchar;not null;unique" chat:"refers_to:create;required:true"`
	Password    *Password `json:"password,omitempty" gorm:"column:password" chat:"refers_to:create;required:true"`
	VisibleName string    `json:"visible_name" gorm:"type:varchar;not null" chat:"refers_to:create;required:true"`
}

func (u *User) Create() error {
	return config.DB.Create(&u).Error
}

func (*User) TableName() string {
	return "user"
}
