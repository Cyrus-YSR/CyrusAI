package router

import (
	"GopherAI/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	enterRouter := r.Group("/api/v1")
	{
		RegisterUserRouter(enterRouter.Group("/user"))
	}
	//后续登录的接口需要jwt鉴权
	{
		AIGroup := enterRouter.Group("/AI")
		AIGroup.Use(jwt.Auth())
		AIRouter(AIGroup)
	}

	{
		ImageGroup := enterRouter.Group("/image")
		ImageGroup.Use(jwt.Auth())
		ImageRouter(ImageGroup)
	}

	{
		FileGroup := enterRouter.Group("/file")
		FileGroup.Use(jwt.Auth())
		FileRouter(FileGroup)
	}

	return r
}
