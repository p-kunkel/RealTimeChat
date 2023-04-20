package mappings

import (
	"RealTimeChat/controllers"

	"github.com/gin-gonic/gin"
)

func userMapping(r *gin.Engine) {

	user := r.Group("user")
	user.POST("registration", controllers.CreateUser)
	user.POST("login", controllers.LoginUser)
	user.POST("refresh_token", controllers.RefreshToken)
}
