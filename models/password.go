package models

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Password struct {
	InputPassword string
	Hash          []byte
}

func (p *Password) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(p.Hash, password)
}

func (p *Password) GenerateHash() error {
	var (
		cost int
		err  error
	)

	if !p.IsValid() {
		return errors.New("password must have minimum 8 characters, at least one uppercase letter, one lowercase letter, one number and one special character")
	}

	if cost, err = strconv.Atoi(os.Getenv("GENERATE_PASSWORD_COST")); err != nil {
		return err
	}

	if p.Hash, err = bcrypt.GenerateFromPassword([]byte(p.InputPassword), cost); err != nil {
		return err
	}

	return nil
}

func (p *Password) IsValid() bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(p.InputPassword) < 8 {
		return false
	}
	if len(p.InputPassword) > 40 {
		return false
	}

	for _, char := range p.InputPassword {
		if !hasUpper {
			hasUpper = unicode.IsUpper(char)
		}
		if !hasLower {
			hasLower = unicode.IsLower(char)
		}
		if !hasNumber {
			hasNumber = unicode.IsNumber(char)
		}
		if !hasSpecial {
			hasSpecial = unicode.IsPunct(char)
		}
		if !hasSpecial {
			hasSpecial = unicode.IsSymbol(char)
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (p *Password) MarshalJSON() ([]byte, error) {
	return json.Marshal("")
}

func (p *Password) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.InputPassword)
}

func (*Password) GormDataType() string {
	return "varchar"
}

func (p *Password) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if p.Hash == nil {
		return clause.Expr{
			SQL: "null",
		}
	}

	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{string(p.Hash)},
	}
}

func (p *Password) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	p.Hash = []byte(val.(string))
	return nil
}
