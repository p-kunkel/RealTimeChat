package models

import (
	"RealTimeChat/config"
	dict "RealTimeChat/dictionaries"
	"RealTimeChat/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User
}

type Token struct {
	Uuid   string `json:"uuid,omitempty" gorm:"type:varchar;not null;default:null"`
	RtUuid string `json:"rt_uuid,omitempty" gorm:"type:varchar;default:null"`
	UserId uint64 `json:"user_id,omitempty" gorm:"type:int8;default:null"`
	Token  string `json:"-" gorm:"type:varchar;not null;default:null"`
	TypeId int8   `json:"type_id,omitempty" gorm:"type:int2;not null;default:null"`
	Exp    int64  `json:"exp,omitempty" gorm:"type:int8;not null;default:null"`

	claims jwt.MapClaims `json:"-" gorm:"-"`

	User *User `json:"-" gorm:"foreignKey:user_id"`
}

func NewToken(TokenType dict.DTokenType, DB *gorm.DB) (Token, error) {
	if TokenType.Id <= 0 {
		return Token{}, errors.New("invalid token type")
	}

	t := Token{
		Uuid:   uuid.New().String(),
		TypeId: TokenType.Id,
	}

	if TokenType.ExpMinutes == 0 {
		return Token{}, errors.New("the expiry time of the token was not specified")
	}
	t.Exp = time.Now().Add(time.Duration(int64(time.Minute) * TokenType.ExpMinutes)).Unix()

	return t, nil
}

func NewLoginToken(user User, DB *gorm.DB) (LoginToken, error) {
	var (
		at, rt Token
		lt     = LoginToken{User: &user}
		err    error
	)

	if user.Id == 0 {
		return LoginToken{}, errors.New("invalid user_id")
	}

	if at, err = NewToken(dict.Dicts.TokenType["access_token"], DB); err != nil {
		return LoginToken{}, err
	}

	if rt, err = NewToken(dict.Dicts.TokenType["refresh_token"], DB); err != nil {
		return LoginToken{}, err
	}

	at.RtUuid = rt.Uuid
	at.UserId = user.Id
	rt.UserId = user.Id

	if err = at.GenerateJWT([]byte("SECRET_ACCESS_TOKEN")); err != nil {
		return LoginToken{}, err
	}

	if err = rt.GenerateJWT([]byte("SECRET_REFRESH_TOKEN")); err != nil {
		return LoginToken{}, err
	}

	lt.AccessToken = at.Token
	lt.RefreshToken = rt.Token
	return lt, DB.Create(&[]Token{at, rt}).Error
}

func (lt *LoginToken) Refresh() error {
	var (
		err error
		t   Token
	)

	t.Token = lt.RefreshToken
	if err = t.Decode([]byte("SECRET_REFRESH_TOKEN")); err != nil {
		return err
	}

	if t.ValidType(dict.Dicts.TokenType["refresh_token"]); err != nil {
		return err
	}

	if err = config.DB.Transaction(func(DB *gorm.DB) error {
		if err = t.DeleteFromDB(DB); err != nil {
			return err
		}

		if *lt, err = NewLoginToken(User{Id: t.UserId}, DB); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (t *Token) GenerateJWT(secret []byte) error {
	var (
		err error
	)

	if err = t.MarshalClaims(); err != nil {
		return err
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, t.claims)
	if t.Token, err = at.SignedString(secret); err != nil {
		return fmt.Errorf("error form create token: %s", err)
	}
	return nil
}

func (t *Token) MarshalClaims() error {
	var (
		err error
		b   []byte
	)

	if b, err = json.Marshal(t); err != nil {
		return fmt.Errorf("error form marshal token: %s", err)
	}
	if err := json.Unmarshal(b, &t.claims); err != nil {
		return fmt.Errorf("error form marshal token: %s", err)
	}
	return nil
}

func (t *Token) UnmarshalClaims() error {
	var (
		err error
		b   []byte
	)

	if b, err = json.Marshal(t.claims); err != nil {
		return fmt.Errorf("error form marshal token: %s", err)
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return fmt.Errorf("error form marshal token: %s", err)
	}
	return nil
}

func (t *Token) Decode(secret []byte) error {
	if err := t.Valid(secret); err != nil {
		return err
	}
	return t.UnmarshalClaims()
}

func (t *Token) Valid(secretKey []byte) error {
	var (
		err      error
		ok       bool
		jwtToken *jwt.Token
	)
	if jwtToken, err = jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}); err != nil {
		return err
	}

	t.claims, ok = jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return errors.New("token is invalid")
	}

	return nil
}

func (t *Token) ValidType(tType dict.DTokenType) error {
	if t.TypeId != tType.Id {
		return errors.New("invalid token type")
	}
	return nil
}

func (t *Token) AddToDB(DB *gorm.DB) error {
	return DB.Create(&t).Error
}

func (t *Token) DeleteFromDB(DB *gorm.DB) error {
	switch t.TypeId {
	case dict.Dicts.TokenType["access_token"].Id:
		DB = DB.Where("uuid IN (?, ?)", t.Uuid, t.RtUuid)

	case dict.Dicts.TokenType["refresh_token"].Id:
		DB = DB.Where("? IN (uuid, rt_uuid)", t.Uuid)

	default:
		return errors.New("delete token: invalid token type")
	}

	return helpers.RecordMustExist(DB.Delete(&t))
}
