package mappings

import "github.com/gin-gonic/gin"

func RunServer() error {
	r := gin.Default()

	userMapping(r)
	chatMapping(r)

	return r.Run()
}
