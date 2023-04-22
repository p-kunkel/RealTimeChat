package mappings

import (
	contr "RealTimeChat/controllers"
	middl "RealTimeChat/middlewares"

	"github.com/gin-gonic/gin"
)

func chatMapping(r *gin.Engine) {
	r.GET("chats", middl.Authenticate(), contr.GetUserChats)

	chat := r.Group("chat", middl.Authenticate())
	{
		chat.POST("", contr.CreateChat)
		chat.POST(":chat_id/user", contr.AddMembersToChat)
		chat.POST(":chat_id/message")

		chat.GET(":chat_id/messages")
		chat.GET(":chat_id/listen")
	}

}
