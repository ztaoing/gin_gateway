package http_proxy_router

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/http_proxy_middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//默认配置的中间件
	router.Use(http_proxy_middleware.HTTPAccessModeMiddleware())

	return router
}
