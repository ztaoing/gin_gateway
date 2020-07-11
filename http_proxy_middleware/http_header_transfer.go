/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/pkg/errors"
	"strings"
)

//请求重写层实现，就是请求来源的信息进行更改
// header转换
//在这里增加了header，就需要在下游读取出来
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)
		for _, items := range strings.Split(serviceDetails.HTTPRule.HeaderTransfor, ",") {
			item := strings.Split(items, " ")
			//add name value
			if len(item) != 3 {
				//参数不符合要求
				continue
			}
			//添加和修改时 为重置
			if item[0] == "add" || item[0] == "edit" {
				c.Request.Header.Set(item[1], item[2])
			}
			if item[0] == "del" {
				c.Request.Header.Del(item[1])
			}

		}
		c.Next()
	}
}
