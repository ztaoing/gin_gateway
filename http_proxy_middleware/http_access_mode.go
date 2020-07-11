/**
* @Author:zhoutao
* @Date:2020/7/10 下午2:45
 */

package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
)

//匹配接入方式
//使用请求信息和服务列表的匹配关系获得所需要的服务的配置，放到上下文中
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//执行服务管理中的访问模式匹配
		serviceDetails, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			//匹配失败
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		//如果能够成功匹配
		fmt.Printf("matched service %s", public.Obj2Json(serviceDetails))
		//添加到上下文中
		c.Set("service", serviceDetails)
		c.Next()
	}
}
