package dictionaries

import "RealTimeChat/config"

var Dicts dicts

type Dictionary struct {
	Id    int8   `json:"-" gorm:"type:smallserial;primarykey"`
	Key   string `json:"-" gorm:"type:varchar; not null"`
	Value string `json:"-" gorm:"type:varchar"`
}

type DTokenType struct {
	Dictionary
	ExpMinutes int64 `json:"-" gorm:"type:int8;not null"`
}

type dicts struct {
	TokenType map[string]DTokenType
}

func (d *dicts) LoadFromDB() error {
	var (
		err error
		tt  []DTokenType
	)

	d.TokenType = map[string]DTokenType{}
	if err = config.DB.Find(&tt).Error; err != nil {
		return err
	}

	for _, v := range tt {
		d.TokenType[v.Key] = v
	}
	return nil
}
