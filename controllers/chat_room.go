package controllers

import (
	"RealTimeChat/models"

	"github.com/gin-gonic/gin"
)

func CreateChat(c *gin.Context) {
	var (
		chatRoom models.ChatRoom
		err      error
		session  models.Session
	)

	if err = session.GetFromContext(c); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = chatRoom.Create(); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = chatRoom.AddMembers([]uint64{session.UserId}); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, chatRoom)
}
