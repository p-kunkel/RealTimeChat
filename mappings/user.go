package mappings

import (
	contr "RealTimeChat/controllers"
	middl "RealTimeChat/middlewares"

	"github.com/gin-gonic/gin"
)

func userMapping(r *gin.Engine) {
	user := r.Group("user")
	{
		user.POST("registration", contr.CreateUser)
		user.POST("login", contr.LoginUser)
		user.POST("refresh_token", contr.RefreshToken)
	}
	user.Use(middl.Authenticate())
	{
		user.GET("", contr.GetUser)
	}

}
