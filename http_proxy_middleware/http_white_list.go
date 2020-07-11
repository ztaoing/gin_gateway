/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
	"strings"
)

// ip 白名单 访问控制
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)

		WhiteIpList := []string{}
		if serviceDetails.AccessControl.WhiteList != "" {
			WhiteIpList = strings.Split(serviceDetails.AccessControl.WhiteList, ",")
		}

		//权限是否开启
		if serviceDetails.AccessControl.OpenAuth == 1 && len(WhiteIpList) > 0 {
			//当前的ip是否在白名单中
			if !public.InStringSlice(WhiteIpList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("your ip %s are not in whiteIp list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
