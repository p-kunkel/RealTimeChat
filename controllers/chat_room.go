package controllers

import (
	"RealTimeChat/config"
	dict "RealTimeChat/dictionaries"
	"RealTimeChat/models"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateChat(c *gin.Context) {
	var (
		chatRoom  models.ChatRoom
		err       error
		session   models.Session
		chatAdmin = models.ChatMember{RoleId: dict.Dicts.ChatRole["admin"].Id}
	)

	if err = session.GetFromContext(c); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = chatRoom.Create(tx); err != nil {
			return err
		}

		chatAdmin.UserId = session.UserId
		chatAdmin.ChatId = chatRoom.Id
		if err = chatAdmin.Create(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.JSON(200, chatRoom)
}

func AddMembersToChat(c *gin.Context) {
	var (
		chatRoom models.ChatRoom
		err      error
		session  models.Session
		reqBody  map[string][]uint64
	)

	if chatRoom.Id, err = strconv.ParseUint(c.Param("chat_id"), 10, 64); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = session.GetFromContext(c); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	if err = chatRoom.AddMembers(reqBody["user_id"], config.DB); err != nil {
		HandleErrResponse(c, MakeErrResponse(err))
		return
	}

	c.Status(200)
}
