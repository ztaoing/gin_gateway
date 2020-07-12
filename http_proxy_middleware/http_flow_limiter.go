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
)

// 限流控制
func HTTPFlowLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 5001, errors.New("service not find"))
			c.Abort()
			return
		}
		serviceDetails := service.(*dao.ServiceDetail)
		ServiceQPS := serviceDetails.AccessControl.ServiceFlowLimit
		//大于0的时候才会执行
		if ServiceQPS > 0 {

			//服务的限流
			serviceLimiter, err := public.FlowLimitHandler.GetLimiter(public.FloatCountService+serviceDetails.Info.ServiceName, float64(ServiceQPS))
			if err != nil {
				middleware.ResponseError(c, 5002, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				middleware.ResponseError(c, 5003, errors.New(fmt.Sprintf("exceed service limiter %d", ServiceQPS)))
				c.Abort()
				return
			}
		}

		ClientQPS := serviceDetails.AccessControl.ClientIPFlowLimit
		if ClientQPS > 0 {
			//clientip的限流
			clientLimiter, err := public.FlowLimitHandler.GetLimiter(public.FloatCountService+serviceDetails.Info.ServiceName+"_"+c.ClientIP(), float64(ClientQPS))
			if err != nil {
				middleware.ResponseError(c, 5004, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5005, errors.New(fmt.Sprintf("exceed clientIP limiter %d", ClientQPS)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
