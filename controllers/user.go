package controllers

import (
	"RealTimeChat/config"
	"RealTimeChat/helpers"
	"RealTimeChat/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var (
		err        error
		user       models.User
		loginToken models.LoginToken
	)

	if err = c.ShouldBindJSON(&user); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = helpers.LoadStructData(user, false).CheckRequiredField([]string{"create"}); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = user.Password.GenerateHash(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if !helpers.IsValidEmail(user.Email) {
		HandleErrResponse(c, MakeErrResponse(errors.New("invalid email")))
		return
	}

	if err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = user.Create(tx); err != nil {
			return err
		}

		if loginToken, err = models.NewLoginToken(user, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, loginToken)
}

func GetUser(c *gin.Context) {
	var (
		err     error
		user    models.User
		session models.Session
	)

	if err = session.GetFromContext(c); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	user.Id = session.UserId
	if err = user.FindById(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, user)
}
