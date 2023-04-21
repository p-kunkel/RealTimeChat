package controllers

import (
	"RealTimeChat/config"
	"RealTimeChat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	var (
		user       models.User
		err        error
		loginToken models.LoginToken
	)

	if err = c.ShouldBindJSON(&user.LoginData); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = user.Login(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if loginToken, err = models.NewLoginToken(user, config.DB); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, loginToken)
}

func RefreshToken(c *gin.Context) {
	var (
		err        error
		loginToken models.LoginToken
	)

	if err = c.ShouldBindJSON(&loginToken); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = loginToken.Refresh(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err, http.StatusUnauthorized))
		return
	}

	c.JSON(200, loginToken)
}
