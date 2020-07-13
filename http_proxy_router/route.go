package http_proxy_router

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/controller"
	"github.com/go1234.cn/gin_scaffold/http_proxy_middleware"
	"github.com/go1234.cn/gin_scaffold/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	//router := gin.Default()
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//默认配置的中间件
	//HTTPAccessModeMiddleware 匹配接入方式

	//只有根路径才执行全部的中间件
	root := router.Group("/")
	root.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),     //访问模式
		http_proxy_middleware.HTTPFlowCountMiddleware(),      //限流控制
		http_proxy_middleware.HTTPFlowLimiterMiddleware(),    //流量控制
		http_proxy_middleware.HTTPWhiteListMiddleware(),      //白名单
		http_proxy_middleware.HTTPBlackListMiddleware(),      //黑名单
		http_proxy_middleware.HTTPHeaderTransferMiddleware(), //在匹配接入层和反向代理之间执行重写请求
		http_proxy_middleware.HTTPStripUriMiddleware(),       //路径去除
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),     //正则重写
		http_proxy_middleware.HTTPReverseProxyMiddelware(),
	)
	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	//定义oauth注册方法：
	{
		controller.OAuthRegister(oauth)
	}
	return router
}
