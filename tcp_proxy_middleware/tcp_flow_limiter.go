/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package tcp_proxy_middleware

import (
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/public"
)

// 限流控制
func TCPFlowLimiterMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		service := c.Get("service")
		if service == nil {
			//没有拿到service的设置
			//不允许继续向下传递，即不再执行后边的中间件
			//输出
			c.conn.Write([]byte("get none service"))
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
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
		}

		ClientQPS := serviceDetails.AccessControl.ClientIPFlowLimit
		if ClientQPS > 0 {
			//clientip的限流
			clientLimiter, err := public.FlowLimitHandler.GetLimiter(public.FloatCountService+serviceDetails.Info.ServiceName+"_"+c.ClientIP(), float64(ClientQPS))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
