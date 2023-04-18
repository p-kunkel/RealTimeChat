package controllers

import (
	"RealTimeChat/helpers"
	"RealTimeChat/models"
	"errors"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var (
		err  error
		user models.User
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

	if err = user.Create(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, user)
}
