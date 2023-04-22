package dictionaries

import "RealTimeChat/config"

var Dicts dicts

type Dictionary struct {
	Id    int8   `json:"-" gorm:"type:int2;autoIncrement;primaryKey"`
	Key   string `json:"-" gorm:"type:varchar; not null"`
	Value string `json:"-" gorm:"type:varchar"`
}

type DTokenType struct {
	Dictionary
	ExpMinutes int64 `json:"-" gorm:"type:int8;not null"`
}

type DChatRole struct {
	Dictionary
}

type dicts struct {
	TokenType map[string]DTokenType
	ChatRole  map[string]DChatRole
}

func (d *dicts) LoadFromDB() error {
	var (
		err        error
		tokenTypes []DTokenType
		chatRoles  []DChatRole
	)

	{
		d.TokenType = map[string]DTokenType{}
		if err = config.DB.Find(&tokenTypes).Error; err != nil {
			return err
		}

		for _, v := range tokenTypes {
			d.TokenType[v.Key] = v
		}
	}

	{
		d.ChatRole = map[string]DChatRole{}
		if err = config.DB.Find(&chatRoles).Error; err != nil {
			return err
		}

		for _, v := range chatRoles {
			d.ChatRole[v.Key] = v
		}
	}

	return nil
}
