package controllers

import (
	"RealTimeChat/config"
	"RealTimeChat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	var (
		user models.User
		err  error
		lt   models.LoginToken
	)

	if err = c.ShouldBindJSON(&user.LoginData); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = user.Login(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if lt, err = models.NewLoginToken(user, config.DB); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, lt)
}

func RefreshToken(c *gin.Context) {
	var (
		err error
		lt  models.LoginToken
	)

	if err = c.ShouldBindJSON(&lt); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = lt.Refresh(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err, http.StatusUnauthorized))
		return
	}

	c.JSON(200, lt)
}
