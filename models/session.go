package models

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Session struct {
	UserId uint64
}

func (s *Session) SetInContext(c *gin.Context) error {
	if s == nil {
		return errors.New("session is null")
	}

	c.Set("session", s)
	return nil
}

func (s *Session) GetFromContext(c *gin.Context) error {
	var (
		ok bool
		v  interface{}
	)

	if v, ok = c.Get("session"); !ok {
		return errors.New("no session in context")
	}

	if s, ok = v.(*Session); !ok {
		return errors.New("invalid session model")
	}

	return nil
}
