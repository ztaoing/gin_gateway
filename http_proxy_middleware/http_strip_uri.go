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

//请求重写层实现，就是请求来源的信息进行更改
// strip_uri 对路径进行缩减调整
//例如：网页路径地址是： http:127.0.0.1:8080/test_sdf_string/abbb
//实际要访问的下游地址是： http:127.0.0.1:8080/abbb
func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)
		//前缀匹配,并且开启了StripUri
		if serviceDetails.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetails.HTTPRule.NeedStripUri == 1 {
			//替换前的路径
			fmt.Println("c.Request.URL.Path:", c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetails.HTTPRule.Rule, "", 1)
			//替换后的
			fmt.Println("c.Request.URL.Path:", c.Request.URL.Path)
		}

		c.Next()
	}
}
