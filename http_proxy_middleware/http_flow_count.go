/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
)

// 限流控制
func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)

		//全站的统计

		totalCount, err := public.FlowCountHandler.GetCounter(public.FloatTotal)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCount.Increace()

		//服务的统计
		ServiceCount, err := public.FlowCountHandler.GetCounter(public.FloatCountService + serviceDetails.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 4002, err)
			c.Abort()
			return
		}
		ServiceCount.Increace()

		//租户的统计
		APPCount, err := public.FlowCountHandler.GetCounter(public.FloatCountAPP)
		if err != nil {
			middleware.ResponseError(c, 4003, err)
			c.Abort()
			return
		}
		APPCount.Increace()

		c.Next()
	}
}
