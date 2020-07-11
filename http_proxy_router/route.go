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
	//HTTPAccessModeMiddleware 匹配接入方式
	//
	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),     //访问模式
		http_proxy_middleware.HTTPWhiteListMiddleware(),      //白名单
		http_proxy_middleware.HTTPBlackListMiddleware(),      //黑名单
		http_proxy_middleware.HTTPHeaderTransferMiddleware(), //在匹配接入层和反向代理之间执行重写请求
		http_proxy_middleware.HTTPStripUriMiddleware(),       //路径去除
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),     //正则重写
		http_proxy_middleware.HTTPReverseProxyMiddelware(),
	)

	return router
}
