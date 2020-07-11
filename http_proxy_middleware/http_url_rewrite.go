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
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

//请求重写层实现，就是请求来源的信息进行更改
// url_rewrite url重写
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetails.HTTPRule.UrlRewrite, ",") {
			fmt.Println("item ", item)

			elems := strings.Split(item, " ")
			//第一个是正则 第二个是替换后的规则
			if len(elems) != 2 {
				continue
			}
			regexp, err := regexp.Compile(elems[0])
			if err != nil {
				fmt.Println("regexp.compile failed:", err)
				continue
			}
			fmt.Println("before rewrite ", c.Request.URL.Path)

			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(elems[1]))
			//回写到当前的地址中
			c.Request.URL.Path = string(replacePath)

			fmt.Println("after rewrite ", c.Request.URL.Path)

		}
		c.Next()
	}
}
