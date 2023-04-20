package middlewares

import (
	"RealTimeChat/controllers"
	"RealTimeChat/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Authenticate() func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			token   models.Token
			err     error
			session = &models.Session{}
		)

		if err = token.GetFromHeader(c); err != nil {
			controllers.HandleErrResponse(c, controllers.MakeErrResponse(err, http.StatusForbidden))
			return
		}

		if err = token.Decode([]byte(os.Getenv("SECRET_ACCESS_TOKEN"))); err != nil {
			controllers.HandleErrResponse(c, controllers.MakeErrResponse(err, http.StatusForbidden))
			return
		}

		if err = session.New(token); err != nil {
			controllers.HandleErrResponse(c, controllers.MakeErrResponse(err, http.StatusForbidden))
			return
		}

		if err = session.SetInContext(c); err != nil {
			controllers.HandleErrResponse(c, controllers.MakeErrResponse(err, http.StatusForbidden))
			return
		}
	}
}
