/**
* @Author:zhoutao
* @Date:2020/7/10 下午2:45
 */

package http_proxy_middleware

import (
	"context"
	"github.com/gin-gonic/gin"
)

//匹配接入方式
//使用请求信息和服务列表的匹配关系获得所需要的服务的配置，放到上下文中
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		context.WithValue(c)
		c.Next()
	}
}
